package notes_model;

import "time";

type Folder struct {
	ID int
	UserID int
	Title string
	UpdatedAt time.Time
	CreatedAt time.Time
}