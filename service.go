package main

import (
	"net/http"

	"goaccount/mixin"
)

type RequestValidator interface {
	Validate(*http.Request, interface{}) error
}

type AccountService struct {
	*mixin.ResponseMixin
	validator RequestValidator
	_jwt      JWT
}

func NewAccountService(validator RequestValidator,
	jwt JWT) *AccountService {
	return &AccountService{
		ResponseMixin: mixin.NewResponseMixin(),
		validator:     validator,
		_jwt:          jwt,
	}
}
