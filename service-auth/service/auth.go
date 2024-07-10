package service

import (
	"errors"
	"regexp"
	"time"

	"github.com/gocql/gocql"
	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/model"
	"github.com/nadern96/Chating-System-go/utils"
)

type AuthService struct {
	ctx        ctx.ServiceContext
	EmailRegex *regexp.Regexp
}

func NewAuthService(ctx ctx.ServiceContext) *AuthService {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return &AuthService{
		ctx:        ctx,
		EmailRegex: emailRegex,
	}
}

func (s *AuthService) Register(user model.User) error {
	var err error

	if matches := s.EmailRegex.MatchString(user.Email); !matches {
		return errors.New("INVALID_EMAIL")
	}

	if err := s.insertUserIfNotExists(user.Email); err != nil {
		return err
	}

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

func (s *AuthService) insertUserIfNotExists(email string) error {
	var existingEmail string
	query := `SELECT email FROM user WHERE email = ? LIMIT 1`

	if err := s.ctx.GetCassandra().Query(query, email).Scan(&existingEmail); err != nil {
		if err == gocql.ErrNotFound {
			return nil
		}
		return err
	}

	return errors.New("EMAIL_ALREADY_EXISTS")
}
