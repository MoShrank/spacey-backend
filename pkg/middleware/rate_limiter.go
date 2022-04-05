package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func LimitReachHandler(c *gin.Context) {
	httpconst.WriteLimitReached(c)
}

func RateLimiter(rate limiter.Rate) gin.HandlerFunc {
	store := memory.NewStore()
	limiter := limiter.New(store, rate)

	return mgin.NewMiddleware(
		limiter,
		mgin.WithLimitReachedHandler(LimitReachHandler),
	)
}
