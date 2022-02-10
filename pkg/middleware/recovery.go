package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			httpconst.WriteInternalServerError(c, err)
		}
		httpconst.WriteInternalServerError(c, http.StatusText(http.StatusInternalServerError))
	})
}
