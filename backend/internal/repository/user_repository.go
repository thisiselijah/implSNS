package repository

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"backend/internal/models"
)

// UserRepository 介面定義了使用者資料的操作
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	// --- Profile ---
	GetUserProfileByUserID(userID string) (*models.UserProfile, error)
	UpdateUserProfile(profile *models.UserProfile) error
	CreateUserProfile(profile *models.UserProfile) error
	// --- Follow/Unfollow ---
	FollowUser(followerID, followedID string) error
	UnfollowUser(followerID, followedID string) error
	GetFollowers(userID string) ([]models.User, error)
	GetFollowing(userID string) ([]models.User, error)
}

// mysqlUserRepository 實現了 UserRepository 介面，用於 MySQL 資料庫
type mysqlUserRepository struct {
	db *sql.DB
}

// NewMySQLUserRepository 是 mysqlUserRepository 的建構子
func NewMySQLUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) GetAllUsers() ([]models.User, error) {
	ctx := context.Background()
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error querying all users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var id_uint uint
		if err := rows.Scan(&id_uint, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Printf("Error scanning user row: %v", err)
			continue
		}
		user.ID = strconv.FormatUint(uint64(id_uint), 10)
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration for all users: %v", err)
		return nil, err
	}

	return users, nil
}

// GetFollowers 獲取指定使用者的粉絲列表
func (r *mysqlUserRepository) GetFollowers(userID string) ([]models.User, error) {
	userIDNum, _ := strconv.ParseUint(userID, 10, 64)
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
		var id_uint uint
		if err := rows.Scan(&id_uint, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Printf("Error scanning follower row: %v", err)
			continue
		}
        user.ID = strconv.FormatUint(uint64(id_uint), 10)
		followers = append(followers, user)
	}

	return followers, nil
}

// GetFollowing 獲取指定使用者正在追蹤的列表
func (r *mysqlUserRepository) GetFollowing(userID string) ([]models.User, error) {
	userIDNum, _ := strconv.ParseUint(userID, 10, 64)
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
        var id_uint uint
		if err := rows.Scan(&id_uint, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Printf("Error scanning following row: %v", err)
			continue
		}
        user.ID = strconv.FormatUint(uint64(id_uint), 10)
		following = append(following, user)
	}

	return following, nil
}

// FollowUser 創建一個新的追蹤關係
func (r *mysqlUserRepository) FollowUser(followerID, followedID string) error {
	followerIDNum, _ := strconv.ParseUint(followerID, 10, 64)
	followedIDNum, _ := strconv.ParseUint(followedID, 10, 64)
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
	followerIDNum, _ := strconv.ParseUint(followerID, 10, 64)
	followedIDNum, _ := strconv.ParseUint(followedID, 10, 64)
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
	err := row.Scan(&id_uint, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error scanning user row for GetUserByEmail: %v", err)
		return nil, err
	}
	user.ID = strconv.FormatUint(uint64(id_uint), 10)
	return &user, nil
}

// GetUserByUsername 從 MySQL 資料庫中根據 username 查詢使用者
func (r *mysqlUserRepository) GetUserByUsername(username string) (*models.User, error) {
	ctx := context.Background()
	query := `SELECT id, username, email, password_hash, created_at, updated_at
			   FROM users WHERE username = ?`
	row := r.db.QueryRowContext(ctx, query, username)
	var user models.User
    var id_uint uint
	err := row.Scan(&id_uint, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error scanning user row for GetUserByUsername: %v", err)
		return nil, err
	}
    user.ID = strconv.FormatUint(uint64(id_uint), 10)
	return &user, nil
}

// GetUserByID 從 MySQL 資料庫中根據 ID 查詢使用者
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
	err = row.Scan(&id_uint, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error scanning user row for GetUserByID (ID: %d): %v", idNum, err)
		return nil, err
	}
	user.ID = strconv.FormatUint(uint64(id_uint), 10)
	return &user, nil
}

// GetUserProfileByUserID 根據 user_id 查詢使用者個人資料
func (r *mysqlUserRepository) GetUserProfileByUserID(userID string) (*models.UserProfile, error) {
	userIDNum, _ := strconv.ParseUint(userID, 10, 64)
	ctx := context.Background()
	query := `
        SELECT up.id, up.user_id, up.avatar_url, up.birth_date, up.bio, up.updated_at, u.username
        FROM user_profiles up
        JOIN users u ON up.user_id = u.id
        WHERE up.user_id = ?`
	row := r.db.QueryRowContext(ctx, query, userIDNum)

	var profile models.UserProfile
	var username string
    var dbUserID uint
	err := row.Scan(&profile.ID, &dbUserID, &profile.AvatarURL, &profile.BirthDate, &profile.Bio, &profile.UpdatedAt, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error scanning user profile row for GetUserProfileByUserID (UserID: %d): %v", userIDNum, err)
		return nil, err
	}
    profile.UserID = strconv.FormatUint(uint64(dbUserID), 10)
	profile.Username = username
	return &profile, nil
}

// UpdateUserProfile 更新使用者的個人資料
func (r *mysqlUserRepository) UpdateUserProfile(profile *models.UserProfile) error {
	ctx := context.Background()
    userIDNum, _ := strconv.ParseUint(profile.UserID, 10, 64)
	query := `UPDATE user_profiles SET avatar_url = ?, bio = ?, birth_date = ? WHERE user_id = ?`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error preparing statement for UpdateUserProfile: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, profile.AvatarURL, profile.Bio, profile.BirthDate, userIDNum)
	if err != nil {
		log.Printf("Error executing statement for UpdateUserProfile: %v", err)
		return err
	}
	return nil
}

// CreateUserProfile 創建一筆新的使用者個人資料
func (r *mysqlUserRepository) CreateUserProfile(profile *models.UserProfile) error {
	ctx := context.Background()
    userIDNum, _ := strconv.ParseUint(profile.UserID, 10, 64)
	query := `INSERT INTO user_profiles (user_id, avatar_url, bio, birth_date) VALUES (?, ?, ?, ?)`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error preparing statement for CreateUserProfile: %v", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, userIDNum, profile.AvatarURL, profile.Bio, profile.BirthDate)
	if err != nil {
		log.Printf("Error executing statement for CreateUserProfile: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for CreateUserProfile: %v", err)
	} else {
		profile.ID = uint(id)
	}
	return nil
}