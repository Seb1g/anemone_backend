package notes_model

import "time"

type User struct {
	ID       int
	Email    string
	Password string
}

type Page struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
}
