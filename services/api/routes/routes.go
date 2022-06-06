package routes

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/services/api/handler"
	limiter "github.com/ulule/limiter/v3"
)

func getUrl(hostName, path string) string {
	return "http://" + hostName + "/" + path
}

func proxyWithPath(targetUrl string) gin.HandlerFunc {
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

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

func proxy(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {

		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = serviceName
			req.Host = serviceName
			req.URL.Path = c.Request.URL.Path
		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

func CreateRoutes(router *gin.Engine, cfg config.ConfigInterface) {
	router.GET("/ping", handler.Ping)

	router.Use(middleware.CORSMiddleware(cfg.GetDomain()))
	router.Use(middleware.JSONMiddleware())

	jwt := auth.NewJWT(cfg)
	authMiddleware := middleware.Auth(jwt, cfg)

	userServiceHostName := cfg.GetUserServiceHostName()
	configServiceHostName := "config-service"
	router.GET("/config/frontend", proxyWithPath(getUrl(configServiceHostName, "config/frontend")))

	rate, err := limiter.NewRateFromFormatted("10-M")
	if err != nil {
		panic(err)
	}
	rateLimiterMiddleware := middleware.RateLimiter(rate)

	userGroup := router.Group("/user")
	{
		router.GET(
			"/user",
			authMiddleware,
			proxyWithPath(getUrl(userServiceHostName, "user")),
		)
		userGroup.POST(
			"",
			rateLimiterMiddleware,
			proxyWithPath(getUrl(userServiceHostName, "user")),
		)
		userGroup.POST(
			"/login",
			rateLimiterMiddleware,
			proxyWithPath(getUrl(userServiceHostName, "login")),
		)
		userGroup.GET(
			"/logout",
			authMiddleware,
			proxyWithPath(getUrl(userServiceHostName, "logout")),
		)
	}

	deckServiceHostName := cfg.GetDeckServiceHostName()
	deckGroup := router.Group("/decks").Use(authMiddleware)
	{
		deckGroup.GET("", proxyWithPath(getUrl(deckServiceHostName, "decks")))
		deckGroup.POST("", proxyWithPath(getUrl(deckServiceHostName, "decks")))
		deckGroup.PUT("/:deckID", proxy(deckServiceHostName))
		deckGroup.DELETE("/:deckID", proxy(deckServiceHostName))

		deckGroup.POST("/:deckID/card", proxy(deckServiceHostName))
		deckGroup.POST("/:deckID/cards", proxy(deckServiceHostName))
		deckGroup.PUT("/:deckID/cards/:id", proxy(deckServiceHostName))
		deckGroup.DELETE("/:deckID/cards/:id", proxy(deckServiceHostName))

		deckGroup.GET("/public", proxyWithPath(getUrl(deckServiceHostName, "/public")))
	}

	learningServiceHostName := cfg.GetLearningServiceHostName()
	learningGroup := router.Group("/learning").Use(authMiddleware)
	{
		learningGroup.POST("/session", proxyWithPath(getUrl(learningServiceHostName, "session")))
		learningGroup.PUT("/session", proxyWithPath(getUrl(learningServiceHostName, "session")))
		learningGroup.POST("/event", proxyWithPath(getUrl(learningServiceHostName, "event")))
		learningGroup.GET("/events", proxyWithPath(getUrl(learningServiceHostName, "events")))
		learningGroup.POST(
			"/probabilities",
			proxyWithPath(getUrl(learningServiceHostName, "probabilities")),
		)
	}

	cardGenerationServiceHostName := cfg.GetCardGenerationServiceHostName()

	router.GET(
		"/notes",
		proxy(cardGenerationServiceHostName),
	)

	cardGenerationGroup := router.Group("/notes").
		Use(authMiddleware, middleware.NeedsBeta())
	{
		cardGenerationGroup.POST(
			"",
			proxy(cardGenerationServiceHostName),
		)
		cardGenerationGroup.PUT(
			"/:noteID",
			proxy(cardGenerationServiceHostName),
		)
		cardGenerationGroup.POST(
			"/:noteID/cards",
			proxy(cardGenerationServiceHostName),
		)
	}
}
