package util

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
)

func GetUrl(hostName, path string) string {
	return "http://" + hostName + "/" + path
}

func ProxyWithPath(targetUrl string) gin.HandlerFunc {
	remote, err := url.Parse(targetUrl)

	return func(c *gin.Context) {

		if err != nil {
			httpconst.WriteBadRequest(c, "failed to parse url")
			return
		}

		director := func(req *http.Request) {
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.Host = remote.Host
			req.URL.Path = remote.Path
		}

		proxy := getProxy(director)
		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

func Proxy(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = serviceName
			req.Host = serviceName
			req.URL.Path = c.Request.URL.Path
		}

		proxy := getProxy(director)
		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

func proxyErrorHandler(w http.ResponseWriter, r *http.Request, e error) {
	w.WriteHeader(http.StatusInternalServerError)
	jsonResponse := fmt.Sprintf(
		`{"error": "%d", "message": "%s"}`,
		http.StatusInternalServerError,
		"server error. cannot reach the target service",
	)
	// TODO the error message should be logged here
	w.Write([]byte(jsonResponse))
}

func getProxy(director func(req *http.Request)) *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{Director: director, ErrorHandler: proxyErrorHandler}
	return proxy
}
