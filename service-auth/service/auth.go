package service

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/model"
	"github.com/nadern96/Chating-System-go/utils"
)

type AuthService struct {
	ctx ctx.ServiceContext
}

func NewAuthService(ctx ctx.ServiceContext) *AuthService {
	return &AuthService{
		ctx: ctx,
	}
}

func (s *AuthService) Register(user model.User) error {
	var err error

	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		return err
	}

	user.ID = gocql.TimeUUID()
	user.CreatedAt = time.Now()

	query := `INSERT INTO user (id, username, email, password, createdat) VALUES (?, ?, ?, ?, ?)`

	err = s.ctx.GetCassandra().Query(query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt).Exec()
	if err != nil {
		return err
	}

	return nil
}
