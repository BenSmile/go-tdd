package main

import (
	"github.com/bensmile/go-api-tdd/pkg/common"
	"github.com/bensmile/go-api-tdd/pkg/domain"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (s *server) CreateUser(c *gin.Context) {

	req := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "name, email, password required"})
		return
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	createUser, err := s.store.CreateUser(user)

	if err != nil {
		log.Printf("error creating user: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while creating user"})
		return
	}

	res := struct {
		Id    int64  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Id:    createUser.Id,
		Name:  createUser.Name,
		Email: createUser.Email,
	}

	c.JSON(http.StatusCreated, res)

}

func (s *server) LoginUser(c *gin.Context) {

	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "email, password required"})
		return
	}

	user, err := s.store.FindUserByEmail(req.Email)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"error": "Invalid credentials"})
		return
	}

	err = common.CheckPassword(req.Password, user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
		})
	}

	jwtPayload, err := s.jwt.CreateToken(*user, 24*time.Hour)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "authentication failed",
		})
	}

	c.JSON(http.StatusOK, jwtPayload)

}
