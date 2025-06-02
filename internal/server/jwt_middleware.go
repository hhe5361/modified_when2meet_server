package server

import (
	"better-when2meet/internal/auth"
	"better-when2meet/internal/db"
	"better-when2meet/internal/domain/room"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Autorization Header Format",
			})
			return
		}

		claims, err := auth.ValidateJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userId", claims["userId"]) //user 특정 해주는 건가.
		c.Next()
	}
}

func JWTRoomAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		roomURL := c.Param("url")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization Header Format",
			})
			return
		}

		claims, err := auth.ValidateJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Get room ID from URL
		roomStorage := room.New(db.InitDB())
		roomData, err := roomStorage.GetRoomByUrl(roomURL)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		// Validate room access
		tokenRoomId := int64(claims["roomId"].(float64))
		if tokenRoomId != roomData.ID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to this room"})
			return
		}

		c.Set("userId", claims["userId"])
		c.Set("roomId", claims["roomId"])
		c.Next()
	}
}
