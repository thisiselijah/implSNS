// backend/internal/service/post_service.go
package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"log"
	"strconv"
	"time"
	"fmt"
)

type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository 
	feedRepo repository.FeedRepository // <--- 新增 feed repository
}


func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository, feedRepo repository.FeedRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
		feedRepo: feedRepo, // <--- 初始化 feed repository
	}
}

// CreatePost 處理創建貼文的邏輯
func (s *PostService) CreatePost(ctx context.Context, payload models.CreatePostPayload) (*models.Post, error) {
	post := &models.Post{
		AuthorID: payload.AuthorID,
		Content:  payload.Content,
		Media:    payload.Media,
		Tags:     payload.Tags,
		Location: payload.Location,
	}

	if err := s.postRepo.CreatePost(ctx, post); err != nil {
		log.Printf("Error creating post in service: %v", err)
		return nil, err
	}

	go s.fanOutToFollowers(post)

	// Fan-out logic would go here in a real application

	return post, nil
}

func (s *PostService) fanOutToFollowers(post *models.Post) {
	// 建立一個新的 context，因為原始的 HTTP request context 可能在 fan-out 完成前就結束了
	ctx := context.Background()

	// 1. 獲取發文者的粉絲列表
	authorIDUint, err := strconv.ParseUint(post.AuthorID, 10, 64)
	if err != nil {
		log.Printf("Fan-out failed: could not parse author ID '%s': %v", post.AuthorID, err)
		return
	}

	followers, err := s.userRepo.GetFollowers(uint(authorIDUint))
	if err != nil {
	    log.Printf("Error getting followers for user %s: %v", post.AuthorID, err)
	    return
	}
	
	if len(followers) == 0 {
		log.Printf("User %s has no followers to fan-out to.", post.AuthorID)
		return
	}
	
	// 設定 Feed 內容的存活時間 (TTL)，例如 90 天
	ttl := time.Now().Add(90 * 24 * time.Hour).Unix()

	var feedItems []models.UserFeedItem
	for _, follower := range followers {
		feedItem := models.UserFeedItem{
			// 使用 "USER#" 前綴來建立標準的 PK
			PK:           fmt.Sprintf("USER#%d", follower.ID), 
			// 使用貼文的創建時間作為 SK，方便排序
			SK:           post.CreatedAt,                   
			PostID:       post.PostID,
			AuthorID:     post.AuthorID,
			TTLTimestamp: ttl,
		}
		feedItems = append(feedItems, feedItem)
	}

	// 3. 使用 BatchWriteItem 進行批量寫入以提高效率
	if err := s.feedRepo.BatchAddToFeed(ctx, feedItems); err != nil {
		log.Printf("Failed to execute batch add to feed for post %s: %v", post.PostID, err)
	}

	log.Printf("Fanning out post %s to %d followers.", post.PostID, len(followers))
}

func (s *PostService) GetPostsByUserID(ctx context.Context, userID string, viewerID string) ([]models.PostFeedDTO, error) {
	// 1. 從 repository 獲取原始的貼文資料 (此處已按時間排序)
	posts, err := s.postRepo.GetPostsByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error getting posts from repo for user ID %s: %v", userID, err)
		return nil, err
	}

	if len(posts) == 0 {
		return []models.PostFeedDTO{}, nil
	}

	// 2. 如果瀏覽者已登入，則檢查其按讚狀態
	likedStatusMap := make(map[string]bool)
	if viewerID != "" && len(posts) > 0 {
		var postIDs []string
		for _, post := range posts {
			postIDs = append(postIDs, post.PostID)
		}
		// 調用已有的 repository 方法進行批量檢查
		likedStatusMap, err = s.postRepo.CheckIfPostsLikedBy(ctx, postIDs, viewerID)
		if err != nil {
			log.Printf("Could not check liked status for viewer %s: %v", viewerID, err)
			likedStatusMap = make(map[string]bool)
		}
	}

	// 3. 將 []models.Post 轉換為 []models.PostFeedDTO
	var feedDTOs []models.PostFeedDTO
	for _, post := range posts {
		// ... (省略獲取 authorName 的邏輯) ...
		var authorName string
		authorIDUint, _ := strconv.ParseUint(post.AuthorID, 10, 64)
		user, userErr := s.userRepo.GetUserByID(uint(authorIDUint))
		if userErr != nil {
			authorName = "User ID: " + post.AuthorID
		} else {
			authorName = user.Username
		}
		
		dto := models.PostFeedDTO{
			PostID:       post.PostID,
			AuthorID:     post.AuthorID,
			AuthorName:   authorName,
			Content:      post.Content,
			Media:        post.Media,
			Tags:         post.Tags,
			Location:     post.Location,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
			IsLiked:      likedStatusMap[post.PostID],
		}
		feedDTOs = append(feedDTOs, dto)
	}

	return feedDTOs, nil
}

// UpdatePost 處理更新貼文的邏輯
func (s *PostService) UpdatePost(ctx context.Context, payload models.UpdatePostPayload) (*models.Post, error) {
	// 1. 先獲取原始貼文，以確認其存在並取得完整 Key
	existingPost, err := s.postRepo.GetPostByID(ctx, payload.PostID)
	if err != nil {
		return nil, err // Post not found
	}

	// 2. 更新欄位
	existingPost.Content = payload.Content
	// ... 更新其他允許的欄位

	// 3. 呼叫 repo 進行更新
	if err := s.postRepo.UpdatePost(ctx, existingPost); err != nil {
		return nil, err
	}

	return existingPost, nil
}

// DeletePost 處理刪除貼文的邏輯
func (s *PostService) DeletePost(ctx context.Context, payload models.DeletePostPayload) error {
	// 為了更可靠地刪除，我們先根據 postID 查詢貼文，以獲取完整的 SK
	post, err := s.postRepo.GetPostByID(ctx, payload.PostID)
	if err != nil {
		log.Printf("Cannot delete post: post with ID %s not found. Error: %v", payload.PostID, err)
		return err // Or return a specific "not found" error
	}

	// 檢查操作者是否有權限刪除 (可選，但推薦)
	// if post.AuthorID != payload.AuthorID {
	// 	log.Printf("User %s is not authorized to delete post %s owned by %s", payload.AuthorID, payload.PostID, post.AuthorID)
	// 	return models.ErrForbidden // 您需要定義這個錯誤
	// }

	// 使用從查詢中得到的 AuthorID, PostID 和 CreatedAt 來刪除
	return s.postRepo.DeletePost(ctx, post.AuthorID, post.PostID, post.CreatedAt)
}

func (s *PostService) LikePost(ctx context.Context, postID, userID string) error {
	// 1. Fetch the post first to get its full details (including PK and SK)
	post, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		log.Printf("LikePost failed: could not find post with ID %s. Error: %v", postID, err)
		return errors.New("post not found")
	}

	// 2. Pass the full post object to the repository method
	if err := s.postRepo.AddLike(ctx, post, userID); err != nil {
		log.Printf("Error liking post in service: %v", err)
		return err
	}
	return nil
}

// UnlikePost 處理取消按讚的邏輯
func (s *PostService) UnlikePost(ctx context.Context, postID, userID string) error {

	posts, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		log.Printf("UnlikePost failed: could not find post with ID %s. Error: %v", postID, err)
		return errors.New("post not found")
	}

	if err := s.postRepo.RemoveLike(ctx, posts, userID); err != nil {
		log.Printf("Error unliking post in service: %v", err)
		return err
	}
	return nil
}

// CreateComment 處理新增評論的邏輯
func (s *PostService) CreateComment(ctx context.Context, payload models.CreateCommentPayload) (*models.Comment, error) {
	// 透過 userRepo 查詢作者的使用者名稱
	authorIDUint, _ := strconv.ParseUint(payload.AuthorID, 10, 64)
	user, userErr := s.userRepo.GetUserByID(uint(authorIDUint))
	if userErr != nil {
		return nil, errors.New("author not found")
	}

	comment := &models.Comment{
		PostID:     payload.PostID,
		AuthorID:   payload.AuthorID,
		AuthorName: user.Username, // 填入使用者名稱
		Content:    payload.Content,
	}

	posts, err := s.postRepo.GetPostByID(ctx, payload.PostID)
	if err != nil {
		log.Printf("CreateComment failed: could not find post with ID %s. Error: %v", payload.PostID, err)
		return nil, errors.New("post not found")
	}

	if err := s.postRepo.CreateComment(ctx, posts, comment); err != nil {
		log.Printf("Error creating comment in service: %v", err)
		return nil, err
	}
	return comment, nil
}

// DeleteComment 處理刪除評論的邏輯
func (s *PostService) DeleteComment(ctx context.Context, postID, commentSK, userID string) error {
	// 1. 獲取評論，以進行授權檢查
	comment, err := s.postRepo.GetCommentBySK(ctx, postID, commentSK)
	if err != nil {
		return err
	}

	// 2. 授權檢查：只有評論者本人可以刪除
	if comment.AuthorID != userID {
		return errors.New("user not authorized to delete this comment")
	}

	posts, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		log.Printf("DeleteComment failed: could not find post with ID %s. Error: %v", postID, err)
		return errors.New("post not found")
	}

	// 3. 執行刪除
	err = s.postRepo.DeleteComment(ctx, posts, commentSK)
	if err != nil {
		log.Printf("Error deleting comment in service: %v", err)
		return err
	}

	return nil

}
