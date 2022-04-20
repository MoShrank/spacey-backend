package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNeedsBeta(t *testing.T) {
	tests := []struct {
		TestName   string
		Url        string
		WantStatus int
	}{
		{
			"Valid Request",
			"/notes?isBeta=true",
			http.StatusOK,
		},
		{
			"Missing Parameter",
			"/notes",
			http.StatusUnauthorized,
		},
		{
			"Invalid Parameter",
			"/notes?isBeta=Invalid",
			http.StatusUnauthorized,
		},
		{
			"Invalid Parameter",
			"/notes?isBeta=",
			http.StatusUnauthorized,
		},
		{
			"Not Beta User",
			"/notes?isBeta=False",
			http.StatusUnauthorized,
		},
	}

	for _, test := range tests {

		t.Run(test.TestName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"GET",
				test.Url,
				nil,
			)
			NeedsBeta()(c)

			assert.Equal(t, test.WantStatus, c.Writer.Status())
		})
	}
}
