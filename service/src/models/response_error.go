package models

import "errors"

type ResponseModel struct {
	Message string `json:"message" example:"invalid"`
	Details string `json:"details,omitempty" example:"email is invalid"`
}

func NewError(err error, details string) ResponseModel {
	e := ResponseModel{
		Message: err.Error(),
	}
	if details != "" {
		e.Details = details
	}

	return e
}

func NewSuccess(message string, details string) ResponseModel {
	e := ResponseModel{
		Message: message,
	}
	if details != "" {
		e.Details = details
	}

	return e
}

var ErrorBadRequest = errors.New("bad request")
var ErrorInternalServerError = errors.New("internal server error")
var ErrorNotFound = errors.New("not found")
var ErrorUnauthorized = errors.New("unauthorized")
var ErrorForbidden = errors.New("forbidden")
var ErrorSignupDisabled = errors.New("signup disabled")
var ErrorTooManyRequests = errors.New("too many requests")
