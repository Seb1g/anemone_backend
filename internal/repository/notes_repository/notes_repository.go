package notes_repository

import (
	"anemone_notes/internal/model/notes_model"
	"context"
	"database/sql"
)

type PageRepo struct {
	DB *sql.DB
}

func NewPageRepo(db *sql.DB) *PageRepo {
	return &PageRepo{DB: db}
}

func (r *PageRepo) Create(ctx context.Context, p *notes_model.Page) error {
	q := `INSERT INTO pages (user_id, title, content) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.DB.QueryRowContext(ctx, q, p.UserID, p.Title, p.Content).Scan(&p.ID, &p.CreatedAt)
}

func (r *PageRepo) GetByID(ctx context.Context, id int) (*notes_model.Page, error) {
	q := `SELECT id, user_id, title, content, created_at FROM pages WHERE id=$1`
	var p notes_model.Page
	err := r.DB.QueryRowContext(ctx, q, id).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
