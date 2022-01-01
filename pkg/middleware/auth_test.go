package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareMissingToken(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)

	Auth("")(c)

	assert.Equal(t, c.Writer.Status(), 401)
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.SetCookie("Authorization", "invalid_token", 0, "", "", false, false)

	Auth("")(c)

	assert.Equal(t, c.Writer.Status(), 401)
}
