package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	s "github.com/core-go/search"
	"github.com/gorilla/mux"

	. "go-service/internal/model"
	. "go-service/internal/service"
)

const InternalServerError = "Internal Server Error"

type UserHandler struct {
	service     UserService
	Validate    func(context.Context, interface{}) ([]core.ErrorMessage, error)
	LogError    func(context.Context, string, ...map[string]interface{})
	search      func(ctx context.Context, m interface{}, results interface{}, limit int64, offset int64) (int64, error)
	paramIndex  map[string]int
	filterIndex int
}

func NewUserHandler(search func(ctx context.Context, m interface{}, results interface{}, limit int64, offset int64) (int64, error), service UserService, validate func(context.Context, interface{}) ([]core.ErrorMessage, error), logError func(context.Context, string, ...map[string]interface{})) *UserHandler {
	filterType := reflect.TypeOf(UserFilter{})
	paramIndex := s.BuildParamIndex(filterType)
	filterIndex := s.FindFilterIndex(filterType)
	return &UserHandler{service: service, Validate: validate, LogError: logError, search: search, paramIndex: paramIndex, filterIndex: filterIndex}
}

func (h *UserHandler) All(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.All(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, users)
}

func (h *UserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := h.service.Load(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, IsFound(user), user)
}
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	errors, er2 := h.Validate(r.Context(), &user)
	if er2 != nil {
		h.LogError(r.Context(), er2.Error())
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if len(errors) > 0 {
		JSON(w, http.StatusUnprocessableEntity, errors)
		return
	}
	res, er3 := h.service.Create(r.Context(), &user)
	if er3 != nil {
		h.LogError(r.Context(), er3.Error(), MakeMap(user))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusCreated, res)
}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}
	errors, er2 := h.Validate(r.Context(), &user)
	if er2 != nil {
		h.LogError(r.Context(), er2.Error())
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	if len(errors) > 0 {
		JSON(w, http.StatusUnprocessableEntity, errors)
		return
	}
	res, er3 := h.service.Update(r.Context(), &user)
	if er3 != nil {
		http.Error(w, er3.Error(), http.StatusInternalServerError)
		return
	}
	status := GetStatus(res)
	JSON(w, status, res)
}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	res, err := h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	status := GetStatus(res)
	JSON(w, status, res)
}
func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := UserFilter{Filter: &s.Filter{}}
	s.Decode(r, &filter, h.paramIndex, h.filterIndex)

	users, total, err := h.service.Search(r.Context(), &filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.JSON(w, 200, &s.Result{List: &users, Total: total})
}

func JSON(w http.ResponseWriter, code int, res interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(res)
}
func GetStatus(status int64) int {
	if status <= 0 {
		return http.StatusNotFound
	}
	return http.StatusOK
}
func IsFound(res interface{}) int {
	if isNil(res) {
		return http.StatusNotFound
	}
	return http.StatusOK
}
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
func MakeMap(res interface{}, opts ...string) map[string]interface{} {
	key := "request"
	if len(opts) > 0 && len(opts[0]) > 0 {
		key = opts[0]
	}
	m := make(map[string]interface{})
	b, err := json.Marshal(res)
	if err != nil {
		return m
	}
	m[key] = string(b)
	return m
}
