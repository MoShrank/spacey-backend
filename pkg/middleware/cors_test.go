package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSValidHost(t *testing.T) {
	tests := []struct {
		testName       string
		host           string
		method         string
		wantStatusCode int
		wantHost       string
	}{
		{
			"Preflight Valid Host",
			"http://localhost:3000",
			"OPTIONS",
			204,
			"http://localhost:3000",
		},
		{
			"Preflight Invalid Host",
			"http://localhost:3001",
			"OPTIONS",
			403,
			"",
		},
		{
			"Valid Host",
			"http://localhost:3000",
			"POST",
			200,
			"http://localhost:3000",
		},
		{
			"Invalid Host",
			"http://localhost:3001",
			"POST",
			403,
			"",
		},
		{
			"Valid Host",
			"https://www.spacey.moritz.dev",
			"POST",
			200,
			"https://www.spacey.moritz.dev",
		},
		{
			"Valid Host",
			"https://spacey.moritz.dev",
			"POST",
			200,
			"https://spacey.moritz.dev",
		},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(test.method, "/", nil)
			c.Request.Header.Set("Referer", test.host)

			CORSMiddleware()(c)

			assert.Equal(
				t,
				test.wantHost,
				c.Writer.Header().Get("Access-Control-Allow-Origin"),
			)
			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}

}
