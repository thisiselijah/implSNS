package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"database/sql"
	"errors"
	"log"
)

// ProfileService 結構體
type ProfileService struct {
	userRepo repository.UserRepository
}

// NewProfileService 是 ProfileService 的建構子
func NewProfileService(userRepo repository.UserRepository) *ProfileService {
	return &ProfileService{
		userRepo: userRepo,
	}
}

// GetProfileByUserID 獲取指定 UserID 的個人資料
func (s *ProfileService) GetProfileByUserID(userID string) (*models.UserProfile, error) {
    profile, err := s.userRepo.GetUserProfileByUserID(userID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            // 如果 profile 不存在，可以選擇回傳一個空的 profile，或建立一個預設的
            log.Printf("Profile for user ID %s not found, creating a default one.", userID)

            defaultProfile := &models.UserProfile{
                UserID: userID,
                // 可以設定預設值
                Bio: "This user is lazy and left nothing.",
            }

            if createErr := s.userRepo.CreateUserProfile(defaultProfile); createErr != nil {
                log.Printf("Failed to create default profile for user ID %s: %v", userID, createErr)
                return nil, errors.New("profile not found and failed to create default")
            }
            return defaultProfile, nil
        }
        // 其他資料庫錯誤
        log.Printf("Error getting profile from repo for user ID %s: %v", userID, err)
        return nil, err
    }
    return profile, nil
}

// UpdateBio 更新個人簡介
func (s *ProfileService) UpdateBio(userID string, bio string) (*models.UserProfile, error) {
    profile, err := s.GetProfileByUserID(userID) // 複用 GetProfileByUserID 確保 profile 存在
    if err != nil {
        return nil, err // 如果獲取或創建失敗，則直接返回錯誤
    }

    profile.Bio = bio
    if err := s.userRepo.UpdateUserProfile(profile); err != nil {
        log.Printf("Error updating bio in repo for user ID %s: %v", userID, err)
        return nil, errors.New("failed to update bio")
    }
    return profile, nil
}

// UpdateAvatar 更新頭像 URL
func (s *ProfileService) UpdateAvatar(userID string, AvatarURL string) (*models.UserProfile, error) {
    profile, err := s.GetProfileByUserID(userID) // 複用 GetProfileByUserID 確保 profile 存在
    if err != nil {
        return nil, err
    }

    profile.AvatarURL = AvatarURL
    if err := s.userRepo.UpdateUserProfile(profile); err != nil {
        log.Printf("Error updating avatar in repo for user ID %s: %v", userID, err)
        return nil, errors.New("failed to update avatar")
    }
    return profile, nil
}