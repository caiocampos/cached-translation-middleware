package handler

import (
	"net/http"

	"cached-translation-middleware/internal/model"
	"cached-translation-middleware/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TranslationHandler handles HTTP requests for the /translate endpoint.
type TranslationHandler struct {
	middleware service.MiddlewareService
	logger     *zap.Logger
}

// NewTranslationHandler creates a new TranslationHandler.
func NewTranslationHandler(middleware service.MiddlewareService, logger *zap.Logger) *TranslationHandler {
	return &TranslationHandler{middleware: middleware, logger: logger}
}

// Translate handles POST /translate requests.
func (h *TranslationHandler) Translate(c *gin.Context) {
	var req model.TranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, source, err := h.middleware.Translate(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("translation error", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "translation service unavailable"})
		return
	}

	// Expose cache source as a response header (useful for debugging / observability)
	c.Header("X-Cache", string(source))
	c.JSON(http.StatusOK, resp)
}

// HealthCheck handles GET /health requests.
func (h *TranslationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
