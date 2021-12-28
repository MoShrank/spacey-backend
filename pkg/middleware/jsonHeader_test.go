package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInvalidHeader(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)

	JSONMiddleware()(c)

	assert.Equal(t, c.Writer.Status(), 406)
}

func TestValidHeader(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	JSONMiddleware()(c)

	assert.Equal(t, c.Writer.Status(), 200)
}
