package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) GetByPhone(phone string) (*User, error) {
	var user User
	err := r.pool.QueryRow(context.Background(),
		"SELECT phone, registration_date FROM users WHERE phone=$1", phone).
		Scan(&user.Phone, &user.RegistrationDate)

	if err != nil {
		// اگر ردیف پیدا نشد، خطای استاندارد repository.ErrUserNotFound را برگردان
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) Create(phone string) (*User, error) {
	now := time.Now()
	_, err := r.pool.Exec(context.Background(),
		"INSERT INTO users (phone, registration_date) VALUES ($1, $2)", phone, now)
	if err != nil {
		return nil, err
	}
	return &User{Phone: phone, RegistrationDate: now}, nil
}

func (r *PostgresUserRepository) List(offset, limit int, search string) ([]User, error) {
	rows, err := r.pool.Query(context.Background(),
		"SELECT phone, registration_date FROM users WHERE phone ILIKE $1 ORDER BY registration_date DESC OFFSET $2 LIMIT $3",
		"%"+search+"%", offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Phone, &u.RegistrationDate); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PostgresUserRepository) UpdatePhone(oldPhone string, newPhone string) error {
	cmdTag, err := r.pool.Exec(context.Background(),
		"UPDATE users SET phone=$1 WHERE phone=$2", newPhone, oldPhone)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *PostgresUserRepository) Delete(phone string) error {
	cmdTag, err := r.pool.Exec(context.Background(),
		"DELETE FROM users WHERE phone=$1", phone)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}
