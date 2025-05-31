package repository

import (
	"context"      // 推薦用於資料庫操作，以控制超時和取消
	"database/sql" // 標準 SQL 套件
	"log"          // 簡單日誌

	"backend/internal/models" // 引入你的 User 模型
)
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
}

// mysqlUserRepository 實現了 UserRepository 介面，用於 MySQL 資料庫
type mysqlUserRepository struct {
	db *sql.DB // 資料庫連線池
}

// NewMySQLUserRepository 是 mysqlUserRepository 的建構子
// 它接收一個 *sql.DB 連線池，並返回 UserRepository 介面類型
func NewMySQLUserRepository(db *sql.DB) UserRepository { // Implements UserRepository from repository.go
	return &mysqlUserRepository{db: db}
}

// CreateUser 將新使用者儲存到 MySQL 資料庫
func (r *mysqlUserRepository) CreateUser(user *models.User) error { // Signature matches UserRepository interface
	ctx := context.Background() // 或者從參數傳入 ctx

	query := `INSERT INTO users (username, email, password_hash, created_at, updated_at)
			   VALUES (?, ?, ?, ?, ?)`

	stmt, err := r.db.PrepareContext(ctx, query) //
	if err != nil { //
		log.Printf("Error preparing statement for CreateUser: %v", err) //
		return err //
	}
	defer stmt.Close() //

	result, err := stmt.ExecContext(ctx, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt) //
	if err != nil { //
		log.Printf("Error executing statement for CreateUser: %v", err) //
		return err //
	}

	id, err := result.LastInsertId() //
	if err != nil { //
		log.Printf("Error getting last insert ID for CreateUser: %v", err) //
	} else {
		user.ID = uint(id) // 假設 User ID 是 uint
	}

	return nil //
}

// GetUserByEmail 從 MySQL 資料庫中根據 email 查詢使用者
func (r *mysqlUserRepository) GetUserByEmail(email string) (*models.User, error) { // Signature matches UserRepository interface
	ctx := context.Background() // 或者從參數傳入 ctx

	query := `SELECT id, username, email, password_hash, created_at, updated_at
			   FROM users WHERE email = ?`
	
	row := r.db.QueryRowContext(ctx, query, email) //

	var user models.User //
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt) //
	if err != nil { //
		if err == sql.ErrNoRows { //
			return nil, err // 未找到使用者
		}
		log.Printf("Error scanning user row for GetUserByEmail: %v", err) //
		return nil, err //
	}

	return &user, nil //
}

// GetUserByUsername 從 MySQL 資料庫中根據 username 查詢使用者
func (r *mysqlUserRepository) GetUserByUsername(username string) (*models.User, error) { // Signature matches UserRepository interface
	ctx := context.Background() // 或者從參數傳入 ctx

	query := `SELECT id, username, email, password_hash, created_at, updated_at
			   FROM users WHERE username = ?`

	row := r.db.QueryRowContext(ctx, query, username) //

	var user models.User //
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt) //
	if err != nil { //
		if err == sql.ErrNoRows { //
			return nil, err // 未找到使用者
		}
		log.Printf("Error scanning user row for GetUserByUsername: %v", err) //
		return nil, err //
	}

	return &user, nil //
}