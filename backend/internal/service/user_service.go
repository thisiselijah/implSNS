package service

import (
	"backend/internal/repository"
	"backend/internal/models"
	"errors"
)

// UserService 結構體
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService 是 UserService 的建構子
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetFollowers 獲取粉絲列表
func (s *UserService) GetFollowers(userID string) ([]models.User, error) {
    return s.userRepo.GetFollowers(userID)
}

// GetFollowing 獲取正在追蹤的列表
func (s *UserService) GetFollowing(userID string) ([]models.User, error) {
    return s.userRepo.GetFollowing(userID)
}

// FollowUser 處理追蹤使用者的邏輯
func (s *UserService) FollowUser(followerID, followedID string) error {
    if followerID == followedID {
        return errors.New("user cannot follow themselves")
    }
    return s.userRepo.FollowUser(followerID, followedID)
}

// UnfollowUser 處理取消追蹤使用者的邏輯
func (s *UserService) UnfollowUser(followerID, followedID string) error {
    return s.userRepo.UnfollowUser(followerID, followedID)
}