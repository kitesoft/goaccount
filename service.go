package main

import (
	"net/http"

	"goaccount/mixin"
)

type RequestValidator interface {
	Validate(*http.Request, interface{}) error
}

type Service struct {
	*mixin.ResponseMixin
	validator RequestValidator
	_jwt      JWT
}

func NewService(validator RequestValidator,
	jwt JWT) *Service {
	return &Service{
		ResponseMixin: mixin.NewResponseMixin(),
		validator:     validator,
		_jwt:          jwt,
	}
}
