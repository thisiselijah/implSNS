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
	likeWeight    = 1.0
	commentWeight = 0
	lookbackDays  = 7
)

// TrendingRecommender 包含演算法所需的依賴
type TrendingRecommender struct {
	postRepo         repository.PostRepository
	userRepo         repository.UserRepository
	recommendationRepo repository.RecommendationRepository // <-- 新增依賴
}

// NewTrendingRecommender 是 TrendingRecommender 的建構子
func NewTrendingRecommender(postRepo repository.PostRepository, userRepo repository.UserRepository, recoRepo repository.RecommendationRepository) *TrendingRecommender {
	return &TrendingRecommender{
		postRepo:         postRepo,
		userRepo:         userRepo,
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

	// --- 3. 為每個使用者產生個人化推薦 (此處仍使用模擬使用者列表) ---
	// 在實際應用中，應從 userRepo 獲取所有活躍使用者
	allUsers := []models.User{{ID: 1}, {ID: 10}}
	
	var allRecommendations []models.UserRecommendationItem
	for _, user := range allUsers {
		userID := fmt.Sprintf("%d", user.ID)

		// 模擬獲取該使用者的觀看紀錄
		seenPosts := make(map[string]bool)
		if userID == "10" {
			seenPosts["postABC"] = true
			seenPosts["postGHI"] = true
		}

		// 過濾已看過或自己發的貼文
		var userRecs []models.UserRecommendationItem
		for _, trendingPost := range trendingList {
			if len(userRecs) >= 20 { // 每位使用者最多推薦 20 篇
				break
			}
			if _, found := seenPosts[trendingPost.PostID]; found {
				continue // 過濾看過的
			}

			// 這裡應檢查是否為使用者自己的貼文 (需要 Post 的 AuthorID)
			// ...

			recItem := models.UserRecommendationItem{
				PK:               fmt.Sprintf("USER#%s", userID),
				SK:               trendingPost.Score, // SK 為分數，用於排序
				GSI1PK:           "trending-v1.0",
				GSI1SK:           fmt.Sprintf("USER#%s", userID),
				PostID:           trendingPost.PostID,
				AlgorithmVersion: "trending-v1.0",
				GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
			}
			userRecs = append(userRecs, recItem)
		}
		allRecommendations = append(allRecommendations, userRecs...)
		log.Printf("Generated %d recommendations for USER#%s.", len(userRecs), userID)
	}

	// --- 4. 將所有推薦結果批量寫入 DynamoDB ---
	if len(allRecommendations) > 0 {
		log.Printf("Saving %d recommendations to DynamoDB...", len(allRecommendations))
		err = r.recommendationRepo.SaveRecommendations(ctx, allRecommendations)
		if err != nil {
			return fmt.Errorf("failed to save recommendations: %w", err)
		}
	}

	return nil
}