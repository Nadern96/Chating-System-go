package service

import (
	"context"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/model"
	"github.com/nadern96/Chating-System-go/utils"
)

type AuthService struct {
	ctx        ctx.ServiceContext
	EmailRegex *regexp.Regexp
	jwtKey     []byte
}

func NewAuthService(ctx ctx.ServiceContext) *AuthService {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return &AuthService{
		ctx:        ctx,
		EmailRegex: emailRegex,
		jwtKey:     []byte(os.Getenv("JWT_SECRET_KEY")),
	}
}

func (s *AuthService) Register(ctx context.Context, user model.User) error {
	var err error

	if matches := s.EmailRegex.MatchString(user.Email); !matches {
		return model.ErrInvalidEmail
	}

	existingUser := model.User{}
	if existingUser, err = s.checkUserExistance(ctx, user.Email); err != nil {
		return err
	}

	if existingUser.Email != "" {
		return model.ErrEmailAlreadyExists
	}

	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		return err
	}

	user.ID = gocql.TimeUUID()
	user.CreatedAt = time.Now()

	query := `INSERT INTO user (id, username, email, password, createdat) VALUES (?, ?, ?, ?, ?)`

	err = s.ctx.GetCassandra().Query(query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt).WithContext(ctx).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) checkUserExistance(ctx context.Context, email string) (model.User, error) {
	var user model.User
	query := `SELECT id, username, password, email, createdAt FROM user WHERE email = ? LIMIT 1`

	if err := s.ctx.GetCassandra().Query(query, email).WithContext(ctx).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt); err != nil {
		if err == gocql.ErrNotFound {
			return model.User{}, nil
		}
		return model.User{}, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, user model.User) (string, error) {
	var err error

	if matches := s.EmailRegex.MatchString(user.Email); !matches {
		return "", model.ErrInvalidEmail
	}

	existingUser := model.User{}
	if existingUser, err = s.checkUserExistance(ctx, user.Email); err != nil {
		if err == gocql.ErrNotFound {
			return "", model.ErrEmailNotFound
		}

		return "", err
	}

	passwordMatch := utils.CompareHashPassword(user.Password, existingUser.Password)
	if !passwordMatch {
		return "", model.ErrWrongCred
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &model.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	redisRes := s.ctx.Redis().Set(user.Email, existingUser.ID.String(), time.Minute*10)
	if redisRes.Err() != nil {
		s.ctx.Logger().Println("err : ", redisRes.Err())
		return "", redisRes.Err()
	}

	return tokenString, err
}
