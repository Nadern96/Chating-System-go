package utils

import (
	"context"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/nadern96/Chating-System-go/model"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseToken(tokenString string) (claims *model.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}

func IsAuthorized(ctx context.Context) (bool, string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userId := md.Get("USER_ID")
		if len(userId) > 0 && userId[0] != "" {
			return true, userId[0]
		}
	}

	return false, ""
}
