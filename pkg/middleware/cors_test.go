package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	tests := []struct {
		testName       string
		host           string
		method         string
		wantStatusCode int
		wantOrigin     string
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
			"http://google.com",
			"OPTIONS",
			401,
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
			"http://google.com",
			"POST",
			401,
			"",
		},
		{
			"Valid Host",
			"https://spacey-learn.com",
			"POST",
			200,
			"https://spacey-learn.com",
		},
	}

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(test.method, "/", nil)
			c.Request.Header.Set("Origin", test.host)

			CORSMiddleware("spacey-learn.com")(c)

			assert.Equal(
				t,
				test.wantOrigin,
				c.Writer.Header().Get("Access-Control-Allow-Origin"),
			)
			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}

}
