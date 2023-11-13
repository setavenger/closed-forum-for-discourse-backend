package server

import (
	"backend/src/common"
	"backend/src/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

const ContextUserKey = "user"

func (d *Daemon) BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, p, ok := c.Request.BasicAuth()
		user, err := db.RetrieveUser(d.DB, u)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !ok || common.HashString(p) != user.PasswordHash {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set(ContextUserKey, user)
		c.Next()
	}
}
