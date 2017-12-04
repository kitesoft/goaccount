package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"goaccount/config"
	"goaccount/jwt"
	"goaccount/mixin"
)

type JWT interface {
	PublicJWT() *jwt.JWT
}

type ServerJWT struct {
	_pubJWT *jwt.JWT
}

func PubJWT(serviceName string) *ServerJWT {
	return &ServerJWT{
		_pubJWT: jwt.NewJwt(serviceName, srvConfig.Jwt.Secret, time.Duration(25200)*time.Second),
	}
}

func (this *ServerJWT) PublicJWT() *jwt.JWT {
	return this._pubJWT
}

type TokenMiddleware struct {
	*mixin.ResponseMixin
	_jwt JWT
}

// Middleware is a struct that has a ServeHTTP method
func NewTokenMiddleware(jwt JWT) *TokenMiddleware {
	return &TokenMiddleware{
		ResponseMixin: mixin.NewResponseMixin(),
		_jwt:          jwt,
	}
}

// The middleware handler 不校验权限
func (this *TokenMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	token := req.Header.Get("Account-Token")
	if token == "" {
		logrus.Error("[TokenMiddleware.ServeHTTP] Account-Token is nil")
		this.ResponseUnauthorized(w)
		return
	}

	id, subject, issuedAt, ok := this._jwt.PublicJWT().Decode(token)
	if !ok || id != mux.Vars(req)["user_id"] {
		logrus.Errorf("[TokenMiddleware.ServeHTTP] token check error => ok: %v, id: %s, subject: %s, issuedAt: %d", ok, id, subject, issuedAt)
		this.ResponseUnauthorized(w)
		return
	}

	next(w, req)
}

type TokenUserMiddleware struct {
	*mixin.ResponseMixin
	_jwt JWT
}

// Middleware is a struct that has a ServeHTTP method
func NewUserTokenMiddleware(jwt JWT) *TokenUserMiddleware {
	return &TokenUserMiddleware{
		ResponseMixin: mixin.NewResponseMixin(),
		_jwt:          jwt,
	}
}

// The middleware handler
func (this *TokenUserMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	token := req.Header.Get("Account-Token")
	if token == "" {
		logrus.Error("[TokenMiddleware.ServeHTTP] Account-Token is nil")
		this.ResponseUnauthorized(w)
		return
	}

	id, subject, issuedAt, ok := this._jwt.PublicJWT().Decode(token)
	if !ok || id != mux.Vars(req)["user_id"] || issuedAt&config.UserPerm != config.UserPerm {
		logrus.Errorf("[TokenMiddleware.ServeHTTP] token check error => ok: %v, id: %s, subject: %s, issuedAt: %d", ok, id, subject, issuedAt)
		this.ResponseUnauthorized(w)
		return
	}

	next(w, req)
}

type AdminTokenMiddleware struct {
	*mixin.ResponseMixin
	_jwt JWT
}

// Middleware is a struct that has a ServeHTTP method
func NewAdminTokenMiddleware(jwt JWT) *AdminTokenMiddleware {
	return &AdminTokenMiddleware{
		ResponseMixin: mixin.NewResponseMixin(),
		_jwt:          jwt,
	}
}

// The middleware handler
func (this *AdminTokenMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	token := req.Header.Get("X-Pear-Token")
	if token == "" {
		logrus.Error("[AdminTokenMiddleware.ServeHTTP] X-Pear-Token is nil")
		this.ResponseUnauthorized(w)
		return
	}

	id, subject, issuedAt, ok := this._jwt.PublicJWT().Decode(token)
	if !ok || id != mux.Vars(req)["user_id"] || issuedAt&config.AdminPerm != config.AdminPerm {
		logrus.Errorf("[AdminTokenMiddleware.ServeHTTP] token check error => ok: %v, id: %s, subject: %s", ok, id, subject)
		this.ResponseUnauthorized(w)
		return
	}

	next(w, req)
}
