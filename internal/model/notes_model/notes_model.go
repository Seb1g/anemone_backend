package notes_model

import (
	"time"
	"database/sql"
)

type Page struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	IsDeleted bool
	FolderID  sql.NullInt64
	CreatedAt time.Time
	UpdatedAt time.Time
}
