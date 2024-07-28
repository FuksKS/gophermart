package model

import "errors"

var (
	ErrNotANumber                   = errors.New("order id is not a number")
	ErrOrderAlreadyUploaded         = errors.New("order id has already been uploaded")
	ErrAlreadyUploadedByThisUser    = errors.New("order id has already been uploaded by this user")
	ErrAlreadyUploadedByAnotherUser = errors.New("order id has already been uploaded by another user")
	ErrNotEnoughMoney               = errors.New("not enough money")
)
