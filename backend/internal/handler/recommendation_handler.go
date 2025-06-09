// internal/handler/recommendation_handler.go
package handler

import (
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecommendationHandler 結構
type RecommendationHandler struct {
	recoService *service.RecommendationService
}

// NewRecommendationHandler 是 RecommendationHandler 的建構子
func NewRecommendationHandler(recoService *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		recoService: recoService,
	}
}

// GenerateTrending 處理觸發生成熱門推薦的請求
func (h *RecommendationHandler) GenerateTrending(c *gin.Context) {
	err := h.recoService.GenerateTrendingRecommendations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate trending recommendations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trending recommendation generation process started."})
}