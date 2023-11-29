package app

import (
	"fanclub/internal/dal"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	keyUserID = "userid"
)

func (app *application) authenticate(c *gin.Context) {
	// token := extractBearerToken(c.GetHeader("Authorization"))
	token, err := c.Cookie("token")
	fmt.Printf("%v\n", token)
	if err == nil {
		userID, isExpired, err := dal.NewTokenDAL(app.db).GetUserFromToken(token)
		if err != nil {
			log.Println("get user from token error", err)
		} else if isExpired {
			log.Printf("token has expired; user_id = %d\n", userID)
		} else {
			c.Set(keyUserID, userID)
		}
	}
	c.Next()
}

func (app *application) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", app.config.allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
