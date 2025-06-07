// backend/internal/service/post_service.go
package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"log"
	"strconv"
	"errors"
)

// PostService 結構體新增 userRepo
type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository // <--- 新增 user repository
}

// NewPostService 的建構子，現在需要傳入 userRepo
func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo, // <--- 初始化 user repository
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

// GetPostsByUserID 獲取某位作者的所有貼文，並轉換為前端需要的 DTO 格式
func (s *PostService) GetPostsByUserID(ctx context.Context, userID string) ([]models.PostFeedDTO, error) {
	// 1. 從 repository 獲取原始的貼文資料
	posts, err := s.postRepo.GetPostsByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error getting posts from repo for user ID %s: %v", userID, err)
		return nil, err
	}

	if len(posts) == 0 {
		return []models.PostFeedDTO{}, nil
	}

	// 2. 將 []models.Post 轉換為 []models.PostFeedDTO
	var feedDTOs []models.PostFeedDTO
	for _, post := range posts {
		var authorName string
		// 將 string 型別的 AuthorID 轉換為 uint，以便查詢使用者資料
		authorIDUint, convErr := strconv.ParseUint(post.AuthorID, 10, 64)
		if convErr != nil {
			log.Printf("Error converting AuthorID string '%s' to uint for post '%s': %v. Using fallback name.", post.AuthorID, post.PostID, convErr)
			authorName = "User ID: " + post.AuthorID // 若轉換失敗，提供一個備用名稱
		} else {
			// 透過 userRepo 查詢作者的使用者名稱
			user, userErr := s.userRepo.GetUserByID(uint(authorIDUint))
			if userErr != nil {
				log.Printf("Error fetching user (ID: %d) for post '%s': %v. Using fallback name.", authorIDUint, post.PostID, userErr)
				authorName = "User ID: " + post.AuthorID // 若查詢失敗，提供一個備用名稱
			} else {
				authorName = user.Username
			}
		}

		// 建立 DTO 物件
		dto := models.PostFeedDTO{
			PostID:       post.PostID,
			AuthorID:     post.AuthorID,
			AuthorName:   authorName, // 填入查詢到的作者名稱
			Content:      post.Content,
			Media:        post.Media,
			Tags:         post.Tags,
			Location:     post.Location,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
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