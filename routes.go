package main

import (
	"github.com/gorilla/mux"

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
	}
}
