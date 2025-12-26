package utils

import "github.com/gin-gonic/gin"

// JSONError writes a standardized error response.
func JSONError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error":   true,
		"message": message,
	})
}
