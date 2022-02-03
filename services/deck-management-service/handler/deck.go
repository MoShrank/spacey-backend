package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/deck-management-service/entity"
)

type DeckHandler struct {
	logger      logger.LoggerInterface
	deckUseCase entity.DeckUseCaseInterface
	validator   validator.ValidatorInterface
}

type DeckHandlerInterface interface {
	GetDeck(c *gin.Context)
	GetDecks(c *gin.Context)
	CreateDeck(c *gin.Context)
	UpdateDeck(c *gin.Context)
	DeleteDeck(c *gin.Context)
}

func NewDeckHandler(
	loggerObj logger.LoggerInterface,
	deckUseCase entity.DeckUseCaseInterface,
	validator validator.ValidatorInterface,
) DeckHandlerInterface {
	return &DeckHandler{
		logger:      loggerObj,
		deckUseCase: deckUseCase,
		validator:   validator,
	}
}

func (h *DeckHandler) CreateDeck(c *gin.Context) {
	userID := c.Request.URL.Query().Get("userID")

	if userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	var deck entity.DeckReq
	if err := h.validator.ValidateJSON(c, &deck); err != nil {
		return
	}

	deckRes, err := h.deckUseCase.CreateDeck(userID, &deck)
	if err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteCreated(c, deckRes)

}

func (h *DeckHandler) GetDeck(c *gin.Context) {
	userID := c.Request.URL.Query().Get("userID")

	if userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	deckID := c.Param("deckID")
	if deckID == "" {
		httpconst.WriteBadRequest(c, "Deck ID is required")
		return
	}

	deck, err := h.deckUseCase.GetDeck(userID, deckID)
	if err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteSuccess(c, deck)
}

func (h *DeckHandler) GetDecks(c *gin.Context) {
	userID := c.Request.URL.Query().Get("userID")

	if userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	decks, err := h.deckUseCase.GetDecks(userID)
	if err != nil {
		httpconst.WriteBadRequest(c, err.Error())
		return
	}

	httpconst.WriteSuccess(c, decks)
}

func (h *DeckHandler) UpdateDeck(c *gin.Context) {
	userID := c.Request.URL.Query().Get("userID")

	if userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	var deck entity.DeckReq
	if err := h.validator.ValidateJSON(c, &deck); err != nil {
		return
	}

	deckID := c.Param("deckID")
	if deckID == "" {
		httpconst.WriteBadRequest(c, "Deck ID is required")
		return
	}

	deckRes, err := h.deckUseCase.UpdateDeck(userID, deckID, &deck)
	if err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteSuccess(c, deckRes)
}

func (h *DeckHandler) DeleteDeck(c *gin.Context) {
	userID := c.Request.URL.Query().Get("userID")

	if userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	deckID := c.Param("deckID")
	if deckID == "" {
		httpconst.WriteBadRequest(c, "Deck ID is required")
		return
	}

	if err := h.deckUseCase.DeleteDeck(userID, deckID); err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteSuccess(c, map[string]string{"Message": "Deck deleted"})
}
