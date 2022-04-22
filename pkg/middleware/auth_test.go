package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareMissingToken(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)

	cfg, _ := config.NewConfig()
	Auth(auth.NewJWT(cfg), cfg)(c)

	assert.Equal(t, c.Writer.Status(), 401)
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Add(
		"Cookie",
		"Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTEwNjkyNTQsImp0aSI6IjYyMDRlYmUwM2VlZjE1ZTdkMDFhN2RjNSJ9.Pi39ABX70_XPIdci1mc41j5zf-ENgs01l2o3bgVT5eM",
	)

	cfg, _ := config.NewConfig()
	Auth(auth.NewJWT(cfg), cfg)(c)
	assert.Equal(t, c.Writer.Status(), 401)
}
