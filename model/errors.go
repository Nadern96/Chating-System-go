package model

import "errors"

var ErrUnauthorized = errors.New(Unauthorized)
var ErrInvalidEmail = errors.New(InvalidEmail)

var ErrEmailAlreadyExists = errors.New(EmailAlreadyExists)
var ErrEmailNotFound = errors.New(EmailNotFound)
var ErrWrongCred = errors.New(WrongCred)

var ErrChatIdNotExist = errors.New(ChatIdNotExist)
