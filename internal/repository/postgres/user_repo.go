package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"user-go/internal/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	query := `
		INSERT INTO users (phone, created_at)
		VALUES ($1, $2)
		RETURNING id
	`

	var id string
	err := r.pool.QueryRow(ctx, query, user.Phone, user.CreatedAt).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = id
	return &user, nil
}

func (r *UserRepository) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	query := `
		SELECT id, phone, created_at 
		FROM users 
		WHERE phone = $1
	`

	var user model.User
	err := r.pool.QueryRow(ctx, query, phone).Scan(&user.ID, &user.Phone, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context, page, pageSize int, search string) ([]model.User, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, phone, created_at
		FROM users
		WHERE phone ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	searchPattern := "%" + search + "%"
	rows, err := r.pool.Query(ctx, query, searchPattern, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Phone, &user.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM users WHERE phone ILIKE $1`
	var total int
	err = r.pool.QueryRow(ctx, countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users count: %w", err)
	}

	return users, total, nil
}
