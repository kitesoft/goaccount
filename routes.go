package main

import (
	"github.com/gorilla/mux"

	"goaccount/server"
)

func (s *AccountService) RegisterRoutes(router *mux.Router, prefix string) {
	subRouter := router.PathPrefix(prefix).Subrouter()
	server.AddRoutes(s.GetRoutes(), subRouter)
}

func (s *AccountService) GetRoutes() []server.Route {
	return []server.Route{
		server.Route{
			Name:    "login",
			Method:  "POST",
			Pattern: "/login",
		},
	}
}
