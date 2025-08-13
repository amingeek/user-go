package user

import (
	"database/sql"
	"time"
)

type Repository interface {
	GetUserByPhoneNumber(phoneNumber string) (*User, error)
	RegisterUser(user *User) error
	GetUserByID(id int) (*User, error)
	ListUsers(limit, offset int, search string) ([]*User, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetUserByPhoneNumber(phoneNumber string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow("SELECT id, phone_number, registered_at FROM users WHERE phone_number = $1", phoneNumber).Scan(&user.ID, &user.PhoneNumber, &user.RegisteredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) RegisterUser(user *User) error {
	_, err := r.db.Exec("INSERT INTO users (phone_number, registered_at) VALUES ($1, $2)", user.PhoneNumber, time.Now())
	return err
}

func (r *PostgresRepository) GetUserByID(id int) (*User, error) {
	user := &User{}
	err := r.db.QueryRow("SELECT id, phone_number, registered_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.PhoneNumber, &user.RegisteredAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) ListUsers(limit, offset int, search string) ([]*User, error) {
	var users []*User
	query := "SELECT id, phone_number, registered_at FROM users"
	args := []interface{}{}

	if search != "" {
		query += " WHERE phone_number ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	// حالا، به درستی پارامترهای LIMIT و OFFSET را اضافه می‌کنیم
	// شماره پارامترها (placeholder) بستگی به استفاده از 'search' دارد.
	if search != "" {
		query += " LIMIT $2 OFFSET $3"
		args = append(args, limit, offset)
	} else {
		query += " LIMIT $1 OFFSET $2"
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.PhoneNumber, &user.RegisteredAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
