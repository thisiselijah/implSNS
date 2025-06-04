package repository

import (
	"context"
	"database/sql"
	"log"

	"backend/internal/models"
)

// UserRepository 介面定義了使用者資料的操作
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error) // <--- 新增此方法
}

// mysqlUserRepository 實現了 UserRepository 介面，用於 MySQL 資料庫
type mysqlUserRepository struct {
	db *sql.DB
}

// NewMySQLUserRepository 是 mysqlUserRepository 的建構子
func NewMySQLUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
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
		user.ID = uint(id)
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
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
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
func (r *mysqlUserRepository) GetUserByID(id uint) (*models.User, error) {
	ctx := context.Background()
	query := `SELECT id, username, email, password_hash, created_at, updated_at
			   FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows { // 未找到使用者，返回錯誤以便上層處理
			return nil, err
		}
		log.Printf("Error scanning user row for GetUserByID (ID: %d): %v", id, err)
		return nil, err // 其他掃描錯誤
	}
	return &user, nil
}