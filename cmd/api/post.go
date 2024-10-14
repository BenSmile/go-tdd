package main

import (
	"github.com/bensmile/go-api-tdd/pkg/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *server) CreatePost(c *gin.Context) {

	user := s.getAuthUser(c)

	req := struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{}

	err := c.BindJSON(&req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid string"})
		return
	}

	if req.Title == "" || req.Body == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "title and body required"})
		return
	}

	post, err := s.store.CreatePost(&domain.Post{
		Title:  req.Title,
		Body:   req.Body,
		UserId: user.Id,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.AbortWithStatusJSON(http.StatusCreated, post)

}

func (s *server) GetUserPosts(c *gin.Context) {

	user := s.getAuthUser(c)

	posts, err := s.store.FindPostsByUser(user.Id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.AbortWithStatusJSON(http.StatusOK, posts)

}

func (s *server) getAuthUser(c *gin.Context) domain.User {
	user, ok := c.Get(ContextUser)
	if !ok {
		panic("missing user value in the context")
	}

	return user.(domain.User)
}
