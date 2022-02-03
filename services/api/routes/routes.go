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
			req.URL.Scheme = c.Request.URL.Scheme
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

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.JSONMiddleware())

	authMiddleware := auth.NewJWT(cfg)

	userServiceHostName := cfg.GetUserServiceHostName()
	configServiceHostName := "config-service"

	router.GET("/config/frontend", proxyWithPath(getUrl(configServiceHostName, "config/frontend")))

	userGroup := router.Group("/user")
	{
		router.GET(
			"/user",
			middleware.Auth(authMiddleware),
			proxyWithPath(getUrl(userServiceHostName, "user")),
		)
		userGroup.POST("", proxyWithPath(getUrl(userServiceHostName, "user")))
		userGroup.POST("/login", proxyWithPath(getUrl(userServiceHostName, "login")))
		userGroup.GET("/logout", proxyWithPath(getUrl(userServiceHostName, "logout")))
		userGroup.DELETE("", proxyWithPath(getUrl(userServiceHostName, "users")))
		userGroup.PUT("/password", proxyWithPath(getUrl(userServiceHostName, "password")))
	}

	deckServiceHostName := cfg.GetDeckServiceHostName()
	deckGroup := router.Group("/decks").Use(middleware.Auth(authMiddleware))
	{
		deckGroup.GET("", proxyWithPath(getUrl(deckServiceHostName, "decks")))
		deckGroup.POST("", proxyWithPath(getUrl(deckServiceHostName, "decks")))
		deckGroup.PUT("/:deckID", proxy(deckServiceHostName))
		deckGroup.DELETE("/:deckID", proxy(deckServiceHostName))

		deckGroup.POST("/:deckID/cards", proxy(deckServiceHostName))
		deckGroup.PUT("/:deckID/cards/:id", proxy(deckServiceHostName))
		deckGroup.DELETE("/:deckID/cards/:id", proxy(deckServiceHostName))

		deckGroup.GET("/public", proxyWithPath(getUrl(deckServiceHostName, "/public")))
		deckGroup.POST("/public", handler.CopyPublicDeck)
	}

	learningGroup := router.Group("/learning").Use(middleware.Auth(authMiddleware))
	{
		learningGroup.POST("/session", handler.CreateLearningSession)
		learningGroup.GET("/session/:id", handler.GetLearningSession)
		learningGroup.PUT("/session/:id", handler.UpdateLearningSession)
		learningGroup.GET("/test/:id", handler.TestCard)
	}

	reminder := router.Group("/reminder").Use(middleware.Auth(authMiddleware))
	{
		reminder.POST("", handler.CreateReminder)
		reminder.GET("/:id", handler.GetReminder)
		reminder.PUT("/:id", handler.UpdateReminder)
		reminder.DELETE("/:id", handler.DeleteReminder)
	}

	statistics := router.Group("/statistics").Use(middleware.Auth(authMiddleware))
	{
		statistics.GET("", handler.GetStatistics)
	}
}
