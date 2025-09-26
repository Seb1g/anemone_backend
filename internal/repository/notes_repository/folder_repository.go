package notes_repository

import (
	"anemone_notes/internal/model/notes_model"
	"context"
	"database/sql"
)

type FolderRepo struct {
	DB *sql.DB
}

func NewFolderRepo(db *sql.DB) *FolderRepo {
	return &FolderRepo{DB: db}
}

func (r *FolderRepo) CreateFolder(ctx context.Context, p *notes_model.Folder) (*notes_model.Folder, error) {
	q := `INSERT INTO notes_folder (user_id, title) VALUES ($1, $2) RETURNING *;`
	err := r.DB.QueryRowContext(ctx, q, p.UserID, p.Title).Scan(&p.ID, &p.UserID, &p.Title, &p.UpdatedAt, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *FolderRepo) GetAllFolders(ctx context.Context, id int) ([]*notes_model.Folder, error) {
	q := `SELECT * FROM notes_folder WHERE user_id=$1`
	rows, err := r.DB.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []*notes_model.Folder
	for rows.Next() {
		var p notes_model.Folder
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.UpdatedAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		folders = append(folders, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return folders, nil
}

func (r *FolderRepo) UpdateTitleFolder(ctx context.Context, id int, new_title string) (*notes_model.Folder, error) {
	q := `UPDATE notes_folder SET title=$1, updated_at=NOW() WHERE id=$2 RETURNING *;`
	var updatedFolder notes_model.Folder
	err := r.DB.QueryRowContext(ctx, q, new_title, id).Scan(&updatedFolder.ID, &updatedFolder.UserID, &updatedFolder.Title, &updatedFolder.UpdatedAt, &updatedFolder.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &updatedFolder, nil
}

func (r *FolderRepo) DeleteFolderByID(ctx context.Context, id int) error {
	q := `DELETE FROM notes_folder WHERE id=$1`
	_, err := r.DB.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	return nil
}