package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	AuthorizationType   = "Bearer"
	ContextUser         = "context_user"
)

func (s *server) applyAuthentication() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("Vary", AuthorizationHeader)

		authHeader := c.GetHeader(AuthorizationHeader)

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "missing auth header",
			})
			return
		}

		headerParts := strings.Split(authHeader, " ")

		// "Bearer mon_token"

		// ["Bearer", "mon_token"]

		if len(headerParts) != 2 || headerParts[0] != AuthorizationType {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token type",
			})
			return
		}

		token := headerParts[1]
		payload, err := s.jwt.VerifyToken(token)

		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "token expired",
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			return
		}

		user, err := s.store.FindUserById(payload.UserId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			return
		}

		c.Set(ContextUser, *user)
		c.Next()
	}
}
