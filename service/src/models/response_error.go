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
var ErrorNotImplemented = errors.New("not implemented")
var ErrorUnauthorized = errors.New("unauthorized")
var ErrorForbidden = errors.New("forbidden")
var ErrorUniqueConstraintViolation = errors.New("unique constraint violation")
var ErrorBadFormat = errors.New("bad format")
var ErrorInvalidEmail = errors.New("invalid email address")
var ErrorInvalidType = errors.New("invalid type")
var ErrorBannedAccount = errors.New("banned account")
var ErrorSignupDisabled = errors.New("signup disabled")
var ErrorMaintenance = errors.New("under maintenance")
