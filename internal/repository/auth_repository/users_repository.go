package auth_repository;

import (
	"anemone_notes/internal/model/auth_model"
	"github.com/jmoiron/sqlx"
	"context"
)

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(ctx context.Context, u *auth_model.User) error {
	q := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at`
	return r.DB.QueryRowContext(ctx, q, u.Email, u.Password).Scan(&u.ID, &u.CreatedAt)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*auth_model.User, error) {
	var u auth_model.User
	q := `SELECT id, email, password, created_at FROM users WHERE email=$1`
	err := r.DB.QueryRowContext(ctx, q, email).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, email, newHash string) error {
	q := `UPDATE users SET password=$1 WHERE email=$2`
	_, err := r.DB.ExecContext(ctx, q, newHash, email)
	return err
}
