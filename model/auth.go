package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
)

type Claims struct {
	jwt.StandardClaims
}

type User struct {
	ID        gocql.UUID `json:"id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"createdat"`
}
