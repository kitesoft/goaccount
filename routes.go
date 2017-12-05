package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"goaccount/server"
)

func (s *Service) RegisterRoutes(router *mux.Router, prefix string) {
	subRouter := router.PathPrefix(prefix).Subrouter()
	server.AddRoutes(s.GetRoutes(), subRouter)
}

func (s *Service) GetRoutes() []server.Route {
	return []server.Route{
		server.Route{
			Name:        "create user",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/create",
			HandlerFunc: s.create_user_handle,
		},
		server.Route{
			Name:        "login",
			Method:      "POST",
			Pattern:     "/login",
			HandlerFunc: s.login_handle,
		},
		server.Route{
			Name:        "update user",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/update",
			HandlerFunc: s.update_user_handle,
		},
		server.Route{
			Name:        "delete user",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/delete",
			HandlerFunc: s.delete_user_handle,
		},
		server.Route{
			Name:        "update password",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/password/update",
			HandlerFunc: s.update_password_handle,
		},
		server.Route{
			Name:        "reset password",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/password/reset",
			HandlerFunc: s.reset_password_handle,
		},

		server.Route{
			Name:        "list user",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/list",
			HandlerFunc: s.list_user_handle,
		},

		server.Route{
			Name:        "login_log",
			Method:      "GET",
			Pattern:     "/{user_id:[0-9]+}/log/login",
			HandlerFunc: s.login_log_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},

		server.Route{
			Name:        "logout",
			Method:      "POST",
			Pattern:     "/{user_id:[0-9]+}/logout",
			HandlerFunc: s.logout_handle,
			Middlewares: []negroni.Handler{NewTokenMiddleware(s._jwt)},
		},
	}
}
