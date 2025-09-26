package notes_repository

import (
	"anemone_notes/internal/model/notes_model"
	"context"
	"database/sql"
	"github.com/lib/pq"
	"errors"
)

type PageRepo struct {
	DB *sql.DB
}

func NewPageRepo(db *sql.DB) *PageRepo {
	return &PageRepo{DB: db}
}

func (r *PageRepo) CreateNote(ctx context.Context, p *notes_model.Page) (*notes_model.Page, error) {
	q := `INSERT INTO pages (user_id, title, content) VALUES ($1, $2, $3) RETURNING *;`
	err := r.DB.QueryRowContext(ctx, q, p.UserID, p.Title, p.Content).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.IsDeleted, &p.FolderID, &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PageRepo) GetOneNoteByID(ctx context.Context, id int) (*notes_model.Page, error) {
	q := `SELECT * FROM pages WHERE id=$1;`
	var p notes_model.Page
	err := r.DB.QueryRowContext(ctx, q, id).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.IsDeleted, &p.FolderID, &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PageRepo) GetAll(ctx context.Context, user_id int) ([]*notes_model.Page, error) {
	q := `SELECT * FROM pages WHERE user_id=$1;`
	rows, err := r.DB.QueryContext(ctx, q, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []*notes_model.Page
	for rows.Next() {
		var p notes_model.Page
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.IsDeleted, &p.FolderID, &p.UpdatedAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *PageRepo) UpdateTitleByID(ctx context.Context, id int, new_title string) (*notes_model.Page, error) {
	q := `UPDATE pages SET title=$1, updated_at=NOW() WHERE id=$2 RETURNING *;`
	var updatedPage notes_model.Page
	err := r.DB.QueryRowContext(ctx, q, new_title, id).Scan(&updatedPage.ID, &updatedPage.UserID, &updatedPage.Title, &updatedPage.Content, &updatedPage.IsDeleted, &updatedPage.FolderID, &updatedPage.UpdatedAt, &updatedPage.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, err

	}
	return &updatedPage, nil
}

func (r *PageRepo) UpdateNoteByID(ctx context.Context, id int, new_content string) (*notes_model.Page, error) {
	q := `UPDATE pages SET content=$1, updated_at=NOW() WHERE id=$2 RETURNING *;`
	var updatedPage notes_model.Page
	err := r.DB.QueryRowContext(ctx, q, new_content, id).Scan(&updatedPage.ID, &updatedPage.UserID, &updatedPage.Title, &updatedPage.Content, &updatedPage.IsDeleted, &updatedPage.FolderID, &updatedPage.UpdatedAt, &updatedPage.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return &updatedPage, nil
}

func (r *PageRepo) GetAllNotesFromFolder(ctx context.Context, id int) ([]*notes_model.Page, error) {
	q := `SELECT * FROM pages WHERE folder_id=$1;`
	rows, err := r.DB.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*notes_model.Page
	for rows.Next() {
		var p notes_model.Page
		if err:= rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.IsDeleted, &p.FolderID, &p.UpdatedAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return notes, nil
}

func (r *PageRepo) AddNoteToFolder(ctx context.Context, noteID int, folderID int) (*notes_model.Page, error) {
	q := `UPDATE pages SET folder_id=$1, updated_at=NOW() WHERE id=$2 RETURNING *;`
	var updatedPage notes_model.Page
	// TODO: Возвращать помимо фолдер_айди еще и тайтл фолдера
	err := r.DB.QueryRowContext(ctx, q, folderID, noteID).Scan(&updatedPage.ID, &updatedPage.UserID, &updatedPage.Title, &updatedPage.Content, &updatedPage.IsDeleted, &updatedPage.FolderID, &updatedPage.UpdatedAt, &updatedPage.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedPage, nil
}

func (r *PageRepo) CencelingNoteFromFolder(ctx context.Context, noteID int) (*notes_model.Page, error) {
	q := `UPDATE pages SET folder_id=NULL, updated_at=NOW() WHERE id=$1 RETURNING *;`
	var updatedPage notes_model.Page
	err := r.DB.QueryRowContext(ctx, q, noteID).Scan(&updatedPage.ID, &updatedPage.UserID, &updatedPage.Title, &updatedPage.Content, &updatedPage.IsDeleted, &updatedPage.FolderID, &updatedPage.UpdatedAt, &updatedPage.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &updatedPage, nil
}

func (r *PageRepo) MarkDeletedNote(ctx context.Context, noteID int) error {
	q := `UPDATE pages SET is_deleted=true, updated_at=NOW() WHERE id=$1;`
	_, err := r.DB.ExecContext(ctx, q, noteID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PageRepo) UnmarkDeletedNote(ctx context.Context, noteID int) error {
	q := `UPDATE pages SET is_deleted=false, updated_at=NOW() WHERE id=$1;`
	_, err := r.DB.ExecContext(ctx, q, noteID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PageRepo) MarkDeletedMoreNotes(ctx context.Context, noteIDs[]int) error {
	q := `UPDATE pages SET is_deleted=true, updated_at=NOW() WHERE id=ANY($1);`
	_, err := r.DB.ExecContext(ctx, q, pq.Array(noteIDs))
	if err != nil {
		return err
	}
	return nil
}

func (r *PageRepo) UnmarkDeletedMoreNotes(ctx context.Context, noteIDs[]int) error {
	q := `UPDATE pages SET is_deleted=false, updated_at=NOW() WHERE id=ANY($1);`
	_, err := r.DB.ExecContext(ctx, q, pq.Array(noteIDs))
	if err != nil {
		return err
	}
	return nil
}

func (r *PageRepo) MarkDeletedAllNotes(ctx context.Context, userID int) error {
	q := `UPDATE pages SET is_deleted=true, updated_at=NOW() WHERE user_id=$1;`
	_, err := r.DB.ExecContext(ctx, q, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PageRepo) UnmarkDeletedAllNotes(ctx context.Context, userID int) error {
	q := `UPDATE pages SET is_deleted=false, updated_at=NOW() WHERE user_id=$1;`
	_, err := r.DB.ExecContext(ctx, q, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PageRepo) DeleteAllMarkNotes(ctx context.Context, userID int) error {
	q := `DELETE FROM pages WHERE user_id=$1 AND is_deleted=true;`
	_, err := r.DB.ExecContext(ctx, q, userID)
	if err != nil {
		return err
	}
	return nil
}