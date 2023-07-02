package service

import (
	"context"
	. "github.com/core-go/core"

	. "go-service/internal/model"
)

func NewUserService(repository Repository) *UserUsecase {
	return &UserUsecase{repository: repository}
}

type UserUsecase struct {
	repository Repository
}

func (s *UserUsecase) Load(ctx context.Context, id string) (*User, error) {
	var user User
	ok, err := s.repository.LoadAndDecode(ctx, id, &user)
	if !ok {
		return nil, err
	} else {
		return &user, err
	}
}
func (s *UserUsecase) Create(ctx context.Context, user *User) (int64, error) {
	return s.repository.Insert(ctx, user)
}
func (s *UserUsecase) Update(ctx context.Context, user *User) (int64, error) {
	return s.repository.Update(ctx, user)
}
func (s *UserUsecase) Delete(ctx context.Context, id string) (int64, error) {
	return s.repository.Delete(ctx, id)
}
