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
	res, err := s.Repo.CreateNote(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *PageService) GetPage(ctx context.Context, id int) (*notes_model.Page, error) {
	return s.Repo.GetOneNoteByID(ctx, id)
}

func (s *PageService) GetAllPages(ctx context.Context, user_id int) ([]*notes_model.Page, error) {
	return s.Repo.GetAll(ctx, user_id)
}

func (s *PageService) UpdateTitle(ctx context.Context, id int, new_title string) (*notes_model.Page, error) {
	return s.Repo.UpdateTitleByID(ctx, id,new_title);
}

func (s *PageService) UpdateContent(ctx context.Context, id int, new_content string) (*notes_model.Page, error) {
	return s.Repo.UpdateNoteByID(ctx, id, new_content)
}

func (s *PageService) GetAllNotesFromFolder(ctx context.Context, id int) ([]*notes_model.Page, error) {
	return s.Repo.GetAllNotesFromFolder(ctx, id)
}

func (s *PageService) AddNoteToFolder(ctx context.Context, noteID int, folderID int) (*notes_model.Page, error) {
	return s.Repo.AddNoteToFolder(ctx, noteID, folderID)
}

func (s *PageService) CencelingNoteFromFolder(ctx context.Context, noteID int) (*notes_model.Page, error) {
	return s.Repo.CencelingNoteFromFolder(ctx, noteID)
}

func (s *PageService) MarkDeletedNote(ctx context.Context, noteID int) error {
	return s.Repo.MarkDeletedNote(ctx, noteID)
}

func (s *PageService) UnmarkDeletedNote(ctx context.Context, noteID int) error {
	return s.Repo.UnmarkDeletedNote(ctx, noteID)
}

func (s *PageService) MarkDeletedMoreNotes(ctx context.Context, notesIDs[]int) error {
	return s.Repo.MarkDeletedMoreNotes(ctx, notesIDs)
}

func (s *PageService) UnmarkDeletedMoreNotes(ctx context.Context, notesIDs[]int) error {
	return s.Repo.UnmarkDeletedMoreNotes(ctx, notesIDs)
}

func (s *PageService) MarkDeletedAllNotes(ctx context.Context, userID int) error {
	return s.Repo.MarkDeletedAllNotes(ctx, userID)
}

func (s *PageService) UnmarkDeletedAllNotes(ctx context.Context, userID int) error {
	return s.Repo.UnmarkDeletedAllNotes(ctx, userID)
}

func (s *PageService) DeleteAllMarkNotes(ctx context.Context, userID int) error {
	return s.Repo.DeleteAllMarkNotes(ctx, userID)
}