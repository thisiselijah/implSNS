package handler

import (
	"backend/internal/models"
	"backend/internal/service"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProfileHandler 結構體
type ProfileHandler struct {
	profileService *service.ProfileService
}

// NewProfileHandler 是 ProfileHandler 的建構子
func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfileByUserID 處理獲取個人資料的請求
func (h *ProfileHandler) GetProfileByUserID(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	profile, err := h.profileService.GetProfileByUserID(uint(userID))
	if err != nil {
		// 檢查是否是 "profile not found" 類型的錯誤
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateBio 處理更新個人簡介的請求
func (h *ProfileHandler) UpdateBio(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var payload models.UpdateBioPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	updatedProfile, err := h.profileService.UpdateBio(uint(userID), payload.Bio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bio updated successfully",
		"profile": updatedProfile,
	})
}

// UpdateAvatar 處理更新頭像的請求
func (h *ProfileHandler) UpdateAvatar(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var payload models.UpdateAvatarPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	updatedProfile, err := h.profileService.UpdateAvatar(uint(userID), payload.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Avatar updated successfully",
		"profile": updatedProfile,
	})
}