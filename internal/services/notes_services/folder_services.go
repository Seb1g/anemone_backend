package notes_services;

import (
	"anemone_notes/internal/model/notes_model"
	"anemone_notes/internal/repository/notes_repository"
	"context"
)

type FolderService struct {
	FolderRepo *notes_repository.FolderRepo
}

func NewFolderService(fr *notes_repository.FolderRepo) *FolderService {
	return &FolderService{FolderRepo: fr}
}

func (s *FolderService) CreateFolder(ctx context.Context, user_id int, title string) (*notes_model.Folder, error) {
	p := &notes_model.Folder{UserID: user_id, Title: title}
	res, err := s.FolderRepo.CreateFolder(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *FolderService) GetAllFolders(ctx context.Context, id int) ([]*notes_model.Folder, error) {
	return s.FolderRepo.GetAllFolders(ctx, id)
}

func (s *FolderService) UpdateTitleFolder(ctx context.Context, id int, newTitle string) (*notes_model.Folder, error) {
	return s.FolderRepo.UpdateTitleFolder(ctx, id, newTitle)
}

func (s *FolderService) DeleteFolders(ctx context.Context, id int) error {
	return s.FolderRepo.DeleteFolderByID(ctx, id)
}
