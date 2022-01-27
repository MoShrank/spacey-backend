package routes

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/services/api/handler"
)

func getUrl(hostName, path string) string {
	return "http://" + hostName + "/" + path
}

func proxy(targetUrl string) gin.HandlerFunc {
	remote, err := url.Parse(targetUrl)

	return func(c *gin.Context) {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to parse url",
			})
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

func CreateRoutes(router *gin.Engine, cfg config.ConfigInterface) {
	router.GET("/ping", handler.Ping)

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.JSONMiddleware())

	authMiddleware := auth.NewJWT(cfg)

	userServiceHostName := cfg.GetUserServiceHostName()
	configServiceHostName := "config-service"

	router.GET("/config/frontend", proxy(getUrl(configServiceHostName, "config/frontend")))

	userGroup := router.Group("/user")
	{
		router.GET(
			"/user",
			middleware.Auth(authMiddleware),
			proxy(getUrl(userServiceHostName, "user")),
		)
		userGroup.POST("", proxy(getUrl(userServiceHostName, "user")))
		userGroup.POST("/login", proxy(getUrl(userServiceHostName, "login")))
		userGroup.GET("/logout", proxy(getUrl(userServiceHostName, "logout")))
		userGroup.DELETE("", proxy(getUrl(userServiceHostName, "users")))
		userGroup.PUT("/password", proxy(getUrl(userServiceHostName, "password")))
	}

	deckServiceHostName := cfg.GetDeckServiceHostName()
	deckGroup := router.Group("/deck").Use(middleware.Auth(authMiddleware))
	{
		deckGroup.GET("", handler.GetDecks)
		deckGroup.POST("", proxy(getUrl(deckServiceHostName, "deck")))
		deckGroup.PUT("/:id", handler.UpdateDeck)
		deckGroup.DELETE("/:id", handler.DeleteDeck)

		deckGroup.POST("/:id/card", handler.CreateCard)
		deckGroup.PUT("/:id/card/:card_id", handler.UpdateCard)
		deckGroup.DELETE("/:id/card/:card_id", handler.DeleteCard)

		deckGroup.GET("/public", proxy(getUrl(deckServiceHostName, "/public")))
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
