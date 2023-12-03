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

	/*
		Middleware Setup
	*/
	domain := cfg.GetDomain()
	cors := middleware.CORSMiddleware(domain)
	json := middleware.JSONMiddleware()

	router.Use(cors)

	jwt := auth.NewJWT(cfg)
	auth := middleware.Auth(jwt, cfg)

	userStore := store.NewStore(db)
	emailVerified := middleware.NeedsEmailVerified(userStore, jwt, cfg)

	userRateLimit := cfg.GetUserRateLimit()
	rate, err := limiter.NewRateFromFormatted(fmt.Sprintf("%d-M", userRateLimit))
	if err != nil {
		panic(err)
	}
	rateLimit := middleware.RateLimiter(rate)
	/* ----- */

	jsonEndpoints := router.Group("/")
	jsonEndpoints.Use(json)

	configServiceHostName := cfg.GetConfigServiceHostName()
	jsonEndpoints.GET(
		"/config/frontend",
		util.ProxyWithPath(util.GetUrl(configServiceHostName, "config/frontend")),
	)

	userServiceHostName := cfg.GetUserServiceHostName()
	userGroup := jsonEndpoints.Group("/user")
	{
		router.GET(
			"/user",
			auth,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "user")),
		)
		userGroup.POST(
			"",
			rateLimit,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "user")),
		)
		userGroup.POST(
			"validate",
			rateLimit,
			auth,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "validate")),
		)
		userGroup.GET(
			"validate",
			rateLimit,
			auth,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "validate")),
		)
		userGroup.POST(
			"/login",
			rateLimit,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "login")),
		)
		userGroup.GET(
			"/logout",
			auth,
			util.ProxyWithPath(util.GetUrl(userServiceHostName, "logout")),
		)
	}

	deckServiceHostName := cfg.GetDeckServiceHostName()
	deckGroup := jsonEndpoints.Group("/decks").Use(auth, emailVerified)
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
	learningGroup := jsonEndpoints.Group("/learning").Use(auth, emailVerified)
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
	cardGenerationGroup := jsonEndpoints.Group("/notes").
		Use(auth, emailVerified, middleware.NeedsBeta())
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

	webContentGroup := jsonEndpoints.Group("/post").
		Use(auth, emailVerified)
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

	pdfGroup := jsonEndpoints.Group("/pdf").
		Use(auth, emailVerified)
	{
		pdfGroup.GET(
			"",
			util.Proxy(cardGenerationServiceHostName),
		)

		pdfGroup.GET(
			"/:id/search",
			util.Proxy(cardGenerationServiceHostName),
		)
		pdfGroup.DELETE(
			"/:id",
			util.Proxy(cardGenerationServiceHostName),
		)
	}

	fileUploadGroup := router.Group("/").Use(auth, emailVerified)
	{
		fileUploadGroup.POST(
			"pdf",
			util.Proxy(cardGenerationServiceHostName),
		)
	}

}
