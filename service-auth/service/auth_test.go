package service

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"math/rand"

	"github.com/bmizerany/assert"
	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/model"
	"github.com/nadern96/Chating-System-go/utils"
)

func TestRegister(t *testing.T) {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra().WithRedis()
	authService := NewAuthService(serviceContext)

	i := rand.Int()
	email := "nader" + strconv.Itoa(i) + "@chat.com"
	err := authService.Register(context.Background(), model.User{
		Username: "nadern96",
		Password: "1234342356",
		Email:    email,
	})
	if err != nil && err.Error() == model.ErrInvalidEmail.Error() {
		t.Errorf(err.Error())
	}

	assert.Equalf(t, errors.Is(err, nil), true, "expected err: nil")
}

func TestRegisterWithInvalidEmail(t *testing.T) {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra().WithRedis()
	authService := NewAuthService(serviceContext)

	err := authService.Register(context.Background(), model.User{
		Username: "nadern96",
		Password: "1234342356",
		Email:    "nadershat.com",
	})

	assert.Equalf(t, errors.Is(err, model.ErrInvalidEmail), true, "expected err: invalid email")
}

func TestRegisterEmailUniqueness(t *testing.T) {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra().WithRedis()
	authService := NewAuthService(serviceContext)

	err := authService.Register(context.Background(), model.User{
		Username: "nadern96",
		Password: "1234342356",
		Email:    "nader@chat.com",
	})

	err = authService.Register(context.Background(), model.User{
		Username: "nadern96",
		Password: "1234342356",
		Email:    "nader@chat.com",
	})

	assert.Equalf(t, errors.Is(err, model.ErrEmailAlreadyExists), true, "expected err: email already exists")
}

func TestLogin(t *testing.T) {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra().WithRedis()
	authService := NewAuthService(serviceContext)

	token, err := authService.Login(context.Background(), model.User{
		Email:    "nader@chat.com",
		Password: "1234342356",
	})

	claims, err := utils.ParseToken(token)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	email := claims.StandardClaims.Subject

	assert.Equalf(t, email, "nader@chat.com", "email does not match")
}
