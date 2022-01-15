package testingutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type HTTPTestingContext struct {
	*gin.Context
}

func NewTestingContext(responseRecorder *httptest.ResponseRecorder, method, path, body string,
) *HTTPTestingContext {
	c, _ := gin.CreateTestContext(responseRecorder)
	c.Request, _ = http.NewRequest(
		method,
		path,
		bytes.NewBuffer([]byte(body)),
	)

	return &HTTPTestingContext{c}
}

func (c *HTTPTestingContext) AddQueryParameter(name, value string) *HTTPTestingContext {
	q := c.Request.URL.Query()
	q.Add(name, value)
	c.Request.URL.RawQuery = q.Encode()

	return c
}
