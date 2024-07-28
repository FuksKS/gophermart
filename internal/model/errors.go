package model

import "errors"

var (
	NotANumberError              = errors.New("order id is not a number")
	OrderAlreadyUploaded         = errors.New("order id has already been uploaded")
	AlreadyUploadedByThisUser    = errors.New("order id has already been uploaded by this user")
	AlreadyUploadedByAnotherUser = errors.New("order id has already been uploaded by another user")
	NotEnoughMoneyError          = errors.New("not enough money")
	NoOrdersForAccrualError      = errors.New("no orders for accrual")
	OrderNotRegisteredError      = errors.New("order is not registered")
	TooManyRequestsErr           = errors.New("too many requests")
)
