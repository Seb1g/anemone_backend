package notes_services

import (
	"anemone_notes/internal/model/notes_model"
	"anemone_notes/internal/repository/notes_repository"
	"context"
)

type PageService struct {
	Repo *notes_repository.PageRepo
}

func NewPageService(r *notes_repository.PageRepo) *PageService {
	return &PageService{Repo: r}
}

func (s *PageService) CreatePage(ctx context.Context, userID int, title, content string) (*notes_model.Page, error) {
	p := &notes_model.Page{UserID: userID, Title: title, Content: content}
	if err := s.Repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PageService) GetPage(ctx context.Context, id int) (*notes_model.Page, error) {
	return s.Repo.GetByID(ctx, id)
}
