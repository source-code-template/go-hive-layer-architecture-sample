package service

import (
	"context"
	. "github.com/core-go/core"

	. "go-service/internal/model"
)

func NewUserService(repository Repository) *UserUseCase {
	return &UserUseCase{repository: repository}
}

type UserUseCase struct {
	repository Repository
}

func (s *UserUseCase) Load(ctx context.Context, id string) (*User, error) {
	var user User
	ok, err := s.repository.LoadAndDecode(ctx, id, &user)
	if !ok {
		return nil, err
	} else {
		return &user, err
	}
}
func (s *UserUseCase) Create(ctx context.Context, user *User) (int64, error) {
	return s.repository.Insert(ctx, user)
}
func (s *UserUseCase) Update(ctx context.Context, user *User) (int64, error) {
	return s.repository.Update(ctx, user)
}
func (s *UserUseCase) Delete(ctx context.Context, id string) (int64, error) {
	return s.repository.Delete(ctx, id)
}
