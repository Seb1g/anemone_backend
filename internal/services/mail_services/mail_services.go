package mail_services

import (
	"anemone_notes/internal/model/mail_model"
	"anemone_notes/internal/repository/mail_repository"
	"anemone_notes/internal/utils"
)

type MailService struct {
	repo       *mail_repository.MailRepository
	jwtManager *utils.JWTManager
	domain     string
}

func New(repo *mail_repository.MailRepository, jwtManager *utils.JWTManager, domain string) *MailService {
	return &MailService{
		repo:       repo,
		jwtManager: jwtManager,
		domain:     domain,
	}
}

type GeneratedAddressResponse struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

func (s *MailService) GenerateAddressAndToken() (*GeneratedAddressResponse, error) {
	addr, err := s.repo.CreateTempAddress(s.domain)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtManager.Generate(addr.ID, addr.Address)
	if err != nil {
		return nil, err
	}

	return &GeneratedAddressResponse{
		Address: addr.Address,
		Token:   token,
	}, nil
}

func (s *MailService) GetInboxForAddress(addressID int) ([]mail_model.Email, error) {
	return s.repo.GetEmailsForAddress(addressID)
}
