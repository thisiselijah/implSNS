// backend/internal/service/post_service.go
package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"log"
)

// PostService 結構體
type PostService struct {
	postRepo repository.PostRepository
}

// NewPostService 是 PostService 的建構子
func NewPostService(postRepo repository.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
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

	// Fan-out logic would go here in a real application

	return post, nil
}

// GetPostsByAuthorID 獲取某位作者的所有貼文
func (s *PostService) GetPostsByUserID(ctx context.Context, userID string) ([]models.Post, error) {
	return s.postRepo.GetPostsByUserID(ctx, userID)
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