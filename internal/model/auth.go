package model

import "errors"

type ContextKey string

const UserIDKey ContextKey = "user_id"

var (
	ErrLoginAlreadyExist = errors.New("login already exist")
	ErrWrongLogin        = errors.New("login does not exist")
	ErrWrongPas          = errors.New("wrong password")
)

type LogoPass struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
