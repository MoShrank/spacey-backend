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
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	var deck entity.DeckReq
	if err := h.validator.ValidateJSON(c, &deck); err != nil {
		return
	}

	deckRes, err := h.deckUseCase.CreateDeck(userID.(string), &deck)
	if err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteCreated(c, deckRes)

}

func (h *DeckHandler) GetDeck(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	deckID := c.Param("id")
	if deckID == "" {
		httpconst.WriteBadRequest(c)
		return
	}

	deck, err := h.deckUseCase.GetDeck(userID.(string), deckID)
	if err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteSuccess(c, deck)
}

func (h *DeckHandler) GetDecks(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	decks, err := h.deckUseCase.GetDecks(userID.(string))
	if err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteSuccess(c, decks)

}

func (h *DeckHandler) UpdateDeck(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	var deck entity.DeckReq
	if err := h.validator.ValidateJSON(c, &deck); err != nil {
		return
	}

	deckID := c.Param("id")
	if deckID == "" {
		httpconst.WriteBadRequest(c)
		return
	}

	deckRes, err := h.deckUseCase.UpdateDeck(userID.(string), deckID, &deck)
	if err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteSuccess(c, deckRes)
}

func (h *DeckHandler) DeleteDeck(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok || userID == "" {
		httpconst.WriteUnauthorized(c)
		return
	}

	deckID := c.Param("id")
	if deckID == "" {
		httpconst.WriteBadRequest(c)
		return
	}

	if err := h.deckUseCase.DeleteDeck(userID.(string), deckID); err != nil {
		httpconst.WriteDatabaseError(c)
		return
	}

	httpconst.WriteSuccess(c, map[string]string{"Message": "Deck deleted"})
}
