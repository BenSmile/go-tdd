package main

import (
	"github.com/gin-gonic/gin"
)

func (s *server) setupRoutes() {
	mux := gin.Default()

	v1 := mux.Group("/api/v1")

	user := v1.Group("/users")

	user.POST("", s.CreateUser)
	user.POST("/login", s.LoginUser)

	blog := v1.Group("/blog")

	blog.Use(s.applyAuthentication())
	blog.POST("", s.CreatePost)
	blog.GET("", s.GetUserPosts)

	s.router = mux
}
