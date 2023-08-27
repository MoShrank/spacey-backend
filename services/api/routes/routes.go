package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/pkg/middleware"
	"github.com/moshrank/spacey-backend/services/api/handler"
	"github.com/moshrank/spacey-backend/services/api/store"
	"github.com/moshrank/spacey-backend/services/api/util"
	limiter "github.com/ulule/limiter/v3"
)

func CreateRoutes(router *gin.Engine, cfg config.ConfigInterface, db db.DatabaseInterface) {
	router.GET("/ping", handler.Ping)

	domain := cfg.GetDomain()
	router.Use(middleware.CORSMiddleware(domain))
	router.Use(middleware.JSONMiddleware())

	userStore := store.NewStore(db)
	jwt := auth.NewJWT(cfg)
	authMiddleware := middleware.Auth(jwt, cfg)
	verifyEmailMiddleware := middleware.NeedsEmailVerified(userStore, jwt, cfg)

	userServiceHostName := cfg.GetUserServiceHostName()
	configServiceHostName := cfg.GetConfigServiceHostName()
	router.GET(
		"/config/frontend",
		util.ProxyWithPath(util.GetUrl(configServiceHostName, "config/frontend")),
	)

	userRateLimit := cfg.GetUserRateLimit()
	rate, err := limiter.NewRateFromFormatted(fmt.Sprintf("%d-M", userRateLimit))
	if err != nil {
		panic(err)
	}
	rateLimiterMiddleware := middleware.RateLimiter(rate)

	userGroup := router.Group("/user")
	{
		router.GET(
			"/user",
			authMiddleware,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "user")),
		)
		userGroup.POST(
			"",
			rateLimiterMiddleware,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "user")),
		)
		userGroup.POST(
			"validate",
			rateLimiterMiddleware,
			authMiddleware,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "validate")),
		)
		userGroup.GET(
			"validate",
			rateLimiterMiddleware,
			authMiddleware,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "validate")),
		)
		userGroup.POST(
			"/login",
			rateLimiterMiddleware,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "login")),
		)
		userGroup.GET(
			"/logout",
			authMiddleware,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "logout")),
		)
	}

	deckServiceHostName := cfg.GetDeckServiceHostName()
	deckGroup := router.Group("/decks").Use(authMiddleware, verifyEmailMiddleware)
	{
		deckGroup.GET("", util.ProxyWithPath(util.GetUrl(deckServiceHostName, "decks")))
		deckGroup.POST("", util.ProxyWithPath(util.GetUrl(deckServiceHostName, "decks")))
		deckGroup.PUT("/:deckID", util.Proxy(deckServiceHostName))
		deckGroup.DELETE("/:deckID", util.Proxy(deckServiceHostName))

		deckGroup.POST("/:deckID/card", util.Proxy(deckServiceHostName))
		deckGroup.POST("/:deckID/cards", util.Proxy(deckServiceHostName))
		deckGroup.PUT("/:deckID/cards/:id", util.Proxy(deckServiceHostName))
		deckGroup.DELETE("/:deckID/cards/:id", util.Proxy(deckServiceHostName))

		deckGroup.GET("/public", util.ProxyWithPath(util.GetUrl(deckServiceHostName, "/public")))
	}

	learningServiceHostName := cfg.GetLearningServiceHostName()
	learningGroup := router.Group("/learning").Use(authMiddleware, verifyEmailMiddleware)
	{
		learningGroup.POST(
			"/session",
			util.ProxyWithPath(util.GetUrl(learningServiceHostName, "session")),
		)
		learningGroup.PUT(
			"/session",
			util.ProxyWithPath(util.GetUrl(learningServiceHostName, "session")),
		)
		learningGroup.POST(
			"/event",
			util.ProxyWithPath(util.GetUrl(learningServiceHostName, "event")),
		)
		learningGroup.GET(
			"/events",
			util.ProxyWithPath(util.GetUrl(learningServiceHostName, "events")),
		)
		learningGroup.POST(
			"/probabilities",
			util.ProxyWithPath(util.GetUrl(learningServiceHostName, "probabilities")),
		)
	}

	cardGenerationServiceHostName := cfg.GetCardGenerationServiceHostName()

	cardGenerationGroup := router.Group("/notes").
		Use(authMiddleware, verifyEmailMiddleware, middleware.NeedsBeta())
	{
		cardGenerationGroup.GET(
			"",
			util.Proxy(cardGenerationServiceHostName),
		)
		cardGenerationGroup.POST(
			"",
			util.Proxy(cardGenerationServiceHostName),
		)
		cardGenerationGroup.PUT(
			"/:noteID",
			util.Proxy(cardGenerationServiceHostName),
		)
		cardGenerationGroup.POST(
			"/:noteID/cards",
			util.Proxy(cardGenerationServiceHostName),
		)
		cardGenerationGroup.POST(
			"/:noteID/card",
			util.Proxy(cardGenerationServiceHostName),
		)
	}

	webContentGroup := router.Group("/post").
		Use(authMiddleware, verifyEmailMiddleware)
	{
		webContentGroup.GET(
			"",
			util.Proxy(cardGenerationServiceHostName),
		)
		webContentGroup.POST(
			"",
			util.Proxy(cardGenerationServiceHostName),
		)
		webContentGroup.DELETE(
			"/:id",
			util.Proxy(cardGenerationServiceHostName),
		)
		webContentGroup.GET(
			"/:id/answer",
			util.Proxy(cardGenerationServiceHostName),
		)
		webContentGroup.GET(
			"/search",
			util.Proxy(cardGenerationServiceHostName),
		)
	}

}
