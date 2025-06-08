// internal/recommendation/trending.go
package recommendation

import (
	"context"
	"log"
	"sort"
	"time"
	"fmt"
	"backend/internal/models"
	"backend/internal/repository" 
)

const (
	likeWeight    = 1.0
	commentWeight = 1.5
	lookbackDays  = 7
)

// TrendingRecommender 包含演算法所需的依賴
type TrendingRecommender struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

// TrendingPost 是一個臨時結構，用於排序
type TrendingPost struct {
	PostID string
	Score  float64
}

// GenerateRecommendations 執行主要的推薦邏輯
func (r *TrendingRecommender) GenerateRecommendations(ctx context.Context) error {
	// --- 1. 獲取近期貼文並計算分數 ---
	// 實際應用中，postRepo 應有一個方法來獲取近期貼文
	// allRecentPosts, err := r.postRepo.GetRecentPosts(ctx, lookbackDays)
	// 此處為模擬：
	log.Println("模擬：獲取所有近期貼文...")
	allRecentPosts, _ := r.postRepo.GetPostsByUserID(ctx, "1") // 模擬從某使用者獲取
	anotherUserPosts, _ := r.postRepo.GetPostsByUserID(ctx, "2")
	allRecentPosts = append(allRecentPosts, anotherUserPosts...)


	var trendingList []TrendingPost
	for _, post := range allRecentPosts {
		// post.LikeCount 和 post.CommentCount 已被反正規化，可直接使用
		score := float64(post.LikeCount)*likeWeight + float64(post.CommentCount)*commentWeight
		trendingList = append(trendingList, TrendingPost{PostID: post.PostID, Score: score})
	}

	// --- 2. 排序得到全域熱門列表 ---
	sort.Slice(trendingList, func(i, j int) bool {
		return trendingList[i].Score > trendingList[j].Score
	})
	log.Printf("計算完成，總共有 %d 篇熱門貼文。", len(trendingList))


	// --- 3. 為每個使用者產生個人化推薦 ---
	// allUsers, err := r.userRepo.GetAllUsers(ctx) // 應有此方法
	// 模擬使用者列表：
	allUsers := []models.User{{ID: 1}, {ID: 10}} 
	
	for _, user := range allUsers {
		userID := fmt.Sprintf("%d", user.ID)

		// 獲取該使用者的觀看紀錄 (模擬)
		// seenPosts, _ := r.seenHistoryRepo.GetSeenPosts(ctx, userID)
		// 根據 SeenHistory.json 的範例，USER#10 看過 postABC 和 postGHI
		seenPosts := make(map[string]bool)
		if userID == "10" {
			seenPosts["postABC"] = true
			seenPosts["postGHI"] = true
		}


		// 過濾已看過或自己發的貼文
		var recommendations []models.UserRecommendationItem
		for _, trendingPost := range trendingList {
			// 如果已經推薦滿了，就跳出
			if len(recommendations) >= 20 {
				break
			}
			
			// 檢查是否看過
			if _, found := seenPosts[trendingPost.PostID]; found {
				continue // 已看過，跳過
			}

			// 這裡應檢查是否為使用者自己的貼文 (需要 Post 的 AuthorID)
			// ...

			// 加入推薦列表
			recItem := models.UserRecommendationItem{
				PK:               fmt.Sprintf("USER#%s", userID),
				SK:               fmt.Sprintf("%f", trendingPost.Score), // 使用分數作為 SK，方便排序
				GSI1PK:           "trending-v1.0", // 演算法版本
				GSI1SK:           fmt.Sprintf("USER#%s", userID),
				PostID:           trendingPost.PostID,
				AlgorithmVersion: "trending-v1.0",
				GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
			}
			recommendations = append(recommendations, recItem)
		}

		// --- 4. 將結果寫入 DynamoDB ---
		// err := r.recommendationRepo.SaveRecommendations(ctx, recommendations)
		log.Printf("為 USER#%s 產生了 %d 條推薦。", userID, len(recommendations))
		// (模擬寫入)
	}

	return nil
}