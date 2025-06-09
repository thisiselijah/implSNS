// internal/service/recommendation_service.go
package service

import (
	"backend/internal/recommendation"
	"context"
)

// RecommendationService 結構
type RecommendationService struct {
	trendingRecommender *recommendation.TrendingRecommender
}

// NewRecommendationService 是 RecommendationService 的建構子
func NewRecommendationService(recommender *recommendation.TrendingRecommender) *RecommendationService {
	return &RecommendationService{
		trendingRecommender: recommender,
	}
}

// GenerateTrendingRecommendations 觸發熱門推薦的生成邏輯
func (s *RecommendationService) GenerateTrendingRecommendations(ctx context.Context) error {
	return s.trendingRecommender.GenerateRecommendations(ctx)
}