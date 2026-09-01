package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type contextKey string

const userIDKey contextKey = "userID"

func RequireAuth(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := GetSession(c.Request.Context(), rdb, cookie)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(c.Request.Context(), userIDKey, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}
