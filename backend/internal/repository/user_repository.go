package repository

import (
	"context"
	"database/sql"
	"log"
	"strconv" // <-- 新增 import

	"backend/internal/models"
)

// UserRepository 介面定義了使用者資料的操作
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	// --- Profile ---
	GetUserProfileByUserID(userID string) (*models.UserProfile, error) 
	UpdateUserProfile(profile *models.UserProfile) error
	CreateUserProfile(profile *models.UserProfile) error
	// --- Follow/Unfollow ---
	FollowUser(followerID, followedID string) error // <-- 修改為 string
	UnfollowUser(followerID, followedID string) error // <-- 修改為 string
	GetFollowers(userID string) ([]models.User, error) // <-- 修改為 string
	GetFollowing(userID string) ([]models.User, error) // <-- 修改為 string
}

// mysqlUserRepository 實現了 UserRepository 介面，用於 MySQL 資料庫
type mysqlUserRepository struct {
	db *sql.DB
}

// NewMySQLUserRepository 是 mysqlUserRepository 的建構子
func NewMySQLUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
}

// GetFollowers 獲取指定使用者的粉絲列表
func (r *mysqlUserRepository) GetFollowers(userID string) ([]models.User, error) {
	userIDNum, err := strconv.ParseUint(userID, 10, 64)
	ctx := context.Background()
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.updated_at
		FROM users u
		INNER JOIN follows f ON u.id = f.follower_id
		WHERE f.followed_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, userIDNum)
	if err != nil {
		log.Printf("Error querying followers for user ID %d: %v", userIDNum, err)
		return nil, err
	}
	defer rows.Close()

	var followers []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Printf("Error scanning follower row: %v", err)
			continue
		}
		followers = append(followers, user)
	}

	return followers, nil
}

// GetFollowing 獲取指定使用者正在追蹤的列表
func (r *mysqlUserRepository) GetFollowing(userID string) ([]models.User, error) {
	userIDNum, err := strconv.ParseUint(userID, 10, 64)
	ctx := context.Background()
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.updated_at
		FROM users u
		INNER JOIN follows f ON u.id = f.followed_id
		WHERE f.follower_id = ?`

	rows, err := r.db.QueryContext(ctx, query, userIDNum)
	if err != nil {
		log.Printf("Error querying following for user ID %d: %v", userIDNum, err)
		return nil, err
	}
	defer rows.Close()

	var following []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Printf("Error scanning following row: %v", err)
			continue
		}
		following = append(following, user)
	}

	return following, nil
}

// FollowUser 創建一個新的追蹤關係
func (r *mysqlUserRepository) FollowUser(followerID, followedID string) error {
	followerIDNum, err := strconv.ParseUint(followerID, 10, 64)
	followedIDNum, err := strconv.ParseUint(followedID, 10, 64)
	ctx := context.Background()
	query := "INSERT INTO follows (follower_id, followed_id) VALUES (?, ?)"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error preparing statement for FollowUser: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, followerIDNum, followedIDNum)
	if err != nil {
		log.Printf("Error executing statement for FollowUser: %v", err)
		return err
	}
	return nil
}

// UnfollowUser 移除一個追蹤關係
func (r *mysqlUserRepository) UnfollowUser(followerID, followedID string) error {
	followerIDNum, _:= strconv.ParseUint(followerID, 10, 64)
	followedIDNum, _:= strconv.ParseUint(followedID, 10, 64)
	ctx := context.Background()
	query := "DELETE FROM follows WHERE follower_id = ? AND followed_id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error preparing statement for UnfollowUser: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, followerIDNum, followedIDNum)
	if err != nil {
		log.Printf("Error executing statement for UnfollowUser: %v", err)
		return err
	}
	return nil
}
// CreateUser 將新使用者儲存到 MySQL 資料庫
func (r *mysqlUserRepository) CreateUser(user *models.User) error {
	ctx := context.Background()
	query := `INSERT INTO users (username, email, password_hash, created_at, updated_at)
			   VALUES (?, ?, ?, ?, ?)`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error preparing statement for CreateUser: %v", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Printf("Error executing statement for CreateUser: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for CreateUser: %v", err)
	} else {
		user.ID = strconv.FormatInt(id, 10)
	}
	return nil
}

// GetUserByEmail 從 MySQL 資料庫中根據 email 查詢使用者
func (r *mysqlUserRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx := context.Background()
	query := `SELECT id, username, email, password_hash, created_at, updated_at
			   FROM users WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)
	var user models.User
	var id_uint uint
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	user.ID = strconv.FormatUint(uint64(id_uint), 10) // <-- 手動轉換
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error scanning user row for GetUserByEmail: %v", err)
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 從 MySQL 資料庫中根據 username 查詢使用者
func (r *mysqlUserRepository) GetUserByUsername(username string) (*models.User, error) {
	ctx := context.Background()
	query := `SELECT id, username, email, password_hash, created_at, updated_at
			   FROM users WHERE username = ?`
	row := r.db.QueryRowContext(ctx, query, username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error scanning user row for GetUserByUsername: %v", err)
		return nil, err
	}
	return &user, nil
}

// GetUserByID 從 MySQL 資料庫中根據 ID 查詢使用者 // <--- 新增此方法的實作
func (r *mysqlUserRepository) GetUserByID(id string) (*models.User, error) {
	idNum, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	query := `SELECT id, username, email, password_hash, created_at, updated_at
			   FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, idNum)

	var id_uint uint
	var user models.User
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows { // 未找到使用者，返回錯誤以便上層處理
			return nil, err
		}
		log.Printf("Error scanning user row for GetUserByID (ID: %d): %v", idNum, err)
		return nil, err // 其他掃描錯誤
	}
	user.ID = strconv.FormatUint(uint64(id_uint), 10) // 將掃描出的數字 ID 轉為字串
	return &user, nil
}

// GetUserProfileByUserID 根據 user_id 查詢使用者個人資料
func (r *mysqlUserRepository) GetUserProfileByUserID(userID string) (*models.UserProfile, error) {
	userIDNum, _ := strconv.ParseUint(userID, 10, 64)
    ctx := context.Background()
    // 加入 JOIN 查詢 username
    query := `
        SELECT up.id, up.user_id, up.avatar_url, up.birth_date, up.bio, up.updated_at, u.username
        FROM user_profiles up
        JOIN users u ON up.user_id = u.id
        WHERE up.user_id = ?`
    row := r.db.QueryRowContext(ctx, query, userID)

    var profile models.UserProfile
    var username string
    err := row.Scan(&profile.ID, &profile.UserID, &profile.AvatarURL, &profile.BirthDate, &profile.Bio, &profile.UpdatedAt, &username)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err // 回傳錯誤，讓 service 層知道沒有找到 profile
        }
        log.Printf("Error scanning user profile row for GetUserProfileByUserID (UserID: %d): %v", userIDNum, err)
        return nil, err
    }
    profile.Username = username // 你需要在 models.UserProfile 結構中加上 Username 欄位
    return &profile, nil
}

// UpdateUserProfile 更新使用者的個人資料
func (r *mysqlUserRepository) UpdateUserProfile(profile *models.UserProfile) error {
	ctx := context.Background()
	// 注意：這裡我們假設 profile 物件中包含了所有需要更新的欄位
	// updatedAt 會由資料庫自動更新
	query := `UPDATE user_profiles SET avatar_url = ?, bio = ?, birth_date = ? WHERE user_id = ?`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error preparing statement for UpdateUserProfile: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, profile.AvatarURL, profile.Bio, profile.BirthDate, profile.UserID)
	if err != nil {
		log.Printf("Error executing statement for UpdateUserProfile: %v", err)
		return err
	}
	return nil
}

// CreateUserProfile 創建一筆新的使用者個人資料
func (r *mysqlUserRepository) CreateUserProfile(profile *models.UserProfile) error {
	ctx := context.Background()
	query := `INSERT INTO user_profiles (user_id, avatar_url, bio, birth_date) VALUES (?, ?, ?, ?)`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error preparing statement for CreateUserProfile: %v", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, profile.UserID, profile.AvatarURL, profile.Bio, profile.BirthDate)
	if err != nil {
		log.Printf("Error executing statement for CreateUserProfile: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for CreateUserProfile: %v", err)
		// 即使無法獲取 ID，資料已插入，所以不一定需要回傳錯誤
	} else {
		profile.ID = uint(id)
	}
	return nil
}