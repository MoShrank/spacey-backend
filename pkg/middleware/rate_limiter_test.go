package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	limiter "github.com/ulule/limiter/v3"
)

func TestRateLimit(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/", nil)

	rate, err := limiter.NewRateFromFormatted("10-M")
	rateLimiterMiddleware := RateLimiter(rate)

	for i := 0; i < 11; i++ {
		rateLimiterMiddleware(c)
	}

	assert.Nil(t, err)
	assert.Equal(t, c.Writer.Status(), 429)
}
