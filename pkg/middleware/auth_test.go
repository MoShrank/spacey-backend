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
	Auth(auth.NewJWT(cfg))(c)

	assert.Equal(t, c.Writer.Status(), 401)
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.SetCookie("Authorization", "invalid_token", 0, "", "", false, false)

	cfg, _ := config.NewConfig()
	Auth(auth.NewJWT(cfg))(c)
	assert.Equal(t, c.Writer.Status(), 401)
}
