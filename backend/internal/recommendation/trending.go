// internal/recommendation/trending.go
package recommendation

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"fmt"
	"log"
	"sort"
	"time"
)

const (
	likeWeight           = 1.0
	commentWeight        = 0
	lookbackDays         = 7
	trendingAlgorithmKey = "trending-v1.0" // 定義一個常數作為演算法金鑰
	maxTrendingPosts     = 100             // 儲存前 100 篇熱門貼文
)

// TrendingRecommender 包含演算法所需的依賴
type TrendingRecommender struct {
	postRepo           repository.PostRepository
	userRepo           repository.UserRepository
	recommendationRepo repository.RecommendationRepository
}

// NewTrendingRecommender 是 TrendingRecommender 的建構子
func NewTrendingRecommender(postRepo repository.PostRepository, userRepo repository.UserRepository, recoRepo repository.RecommendationRepository) *TrendingRecommender {
	return &TrendingRecommender{
		postRepo:           postRepo,
		userRepo:           userRepo,
		recommendationRepo: recoRepo,
	}
}

// TrendingPost 是一個臨時結構，用於排序
type TrendingPost struct {
	PostID string
	Score  float64
}

// GenerateRecommendations 執行主要的推薦邏輯
func (r *TrendingRecommender) GenerateRecommendations(ctx context.Context) error {
	// --- 1. 獲取近期貼文並計算分數 ---
	log.Println("Fetching recent posts for trending calculation...")
	allRecentPosts, err := r.postRepo.GetRecentPosts(ctx, lookbackDays)
	if err != nil {
		log.Printf("Error getting recent posts for trending recommendations: %v", err)
		return fmt.Errorf("could not get recent posts: %w", err)
	}

	if len(allRecentPosts) == 0 {
		log.Println("No recent posts found to generate recommendations.")
		return nil
	}

	var trendingList []TrendingPost
	for _, post := range allRecentPosts {
		score := float64(post.LikeCount)*likeWeight + float64(post.CommentCount)*commentWeight
		trendingList = append(trendingList, TrendingPost{PostID: post.PostID, Score: score})
	}

	// --- 2. 排序得到全域熱門列表 ---
	sort.Slice(trendingList, func(i, j int) bool {
		return trendingList[i].Score > trendingList[j].Score
	})
	log.Printf("Calculation complete. Found %d trending posts.", len(trendingList))

	// --- 3. 準備全域熱門列表以供儲存 ---
	// 核心重構：不再遍歷所有使用者。
	// 我們建立一個單一的全域推薦列表。

	var globalTrendingItems []models.UserRecommendationItem
	// 限制儲存的貼文數量為 maxTrendingPosts
	numToSave := len(trendingList)
	if numToSave > maxTrendingPosts {
		numToSave = maxTrendingPosts
	}

	for i := 0; i < numToSave; i++ {
		trendingPost := trendingList[i]

		// 排序鍵基於分數，確保順序，並包含 PostID 以保證唯一性。
		uniqueSortKey := fmt.Sprintf("%010.2f#%s", trendingPost.Score, trendingPost.PostID)

		recItem := models.UserRecommendationItem{
			PK:               fmt.Sprintf("TRENDING#%s", trendingAlgorithmKey), // 使用一個常數 PK 代表全域列表
			SK:               uniqueSortKey,                                   // SK 用於按分數排序
			PostID:           trendingPost.PostID,
			AlgorithmVersion: trendingAlgorithmKey,
			GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
			// GSI 相關鍵在此查詢模式下不再需要。
		}
		globalTrendingItems = append(globalTrendingItems, recItem)
	}

	// --- 4. 將單一的全域列表批量寫入 DynamoDB ---
	if len(globalTrendingItems) > 0 {
		log.Printf("Saving %d global trending posts to DynamoDB...", len(globalTrendingItems))
		err = r.recommendationRepo.SaveRecommendations(ctx, globalTrendingItems)
		if err != nil {
			return fmt.Errorf("failed to save recommendations: %w", err)
		}
	}

	log.Println("Successfully generated and saved global trending recommendations.")
	return nil
}