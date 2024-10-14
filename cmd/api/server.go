package main

import (
	"github.com/bensmile/go-api-tdd/pkg/domain"
	"github.com/gin-gonic/gin"
)

type server struct {
	router *gin.Engine
	store  domain.Store
	jwt    domain.JWT
}

func newServer(store domain.Store, jwt domain.JWT) *server {
	return &server{
		store: store,
		jwt:   jwt,
	}
}

func (s *server) run(addr string) error {
	return s.routes().Run(addr)
}

func (s *server) routes() *gin.Engine {
	if s.router == nil {
		s.setupRoutes()
	}
	return s.router
}
