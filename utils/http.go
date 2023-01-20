package utils

import "github.com/gin-gonic/gin"

// Enum for HTTP methods
type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

// Response - This function writes a success response to the response writer
func Response(c *gin.Context, status int, data any, message string) {
	c.JSON(status, map[string]any{
		"data":    data,
		"message": message,
		"status":  200,
	})
}
