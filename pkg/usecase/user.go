package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/shou1027/golangJwt/pkg/domain/model"
	"github.com/shou1027/golangJwt/pkg/domain/repository"
	"github.com/shou1027/golangJwt/pkg/myerror"
	"github.com/shou1027/golangJwt/pkg/util"
)

type UseCase interface {
	Signup(c context.Context, username, email, password string) (*model.User, error)
	Login(c context.Context, email, password string) (string, *model.User, error)
}

type useCase struct {
	repository repository.Repository
	timeout    time.Duration
}

func NewUseCase(userRepo repository.Repository) UseCase {
	return &UseCase{
		repository: userRepo,
		timeout:    time.Duration(2) * time.Second,
	}
}

func (uc *useCase) Signup(c context.Context, username, email, password string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	exsitUser, err := uc.repository.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, &myerror.InternalServerError{Err: err}
	}

	if exsitUser.ID != 0 {
		return nil, &myerror.BadRequestError{Err: errors.New("user already exists")}
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, &myerror.InternalServerError{Err: err}
	}

	u := &model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	user, err := uc.repository.CreateUser(ctx, u)
	if err != nil {
		return nil, &myerror.InternalServerError{Err: err}
	}

	return user, nil
}
