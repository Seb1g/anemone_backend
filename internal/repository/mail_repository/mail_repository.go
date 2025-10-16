package mail_repository

import (
	"anemone_notes/internal/model/mail_model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type MailRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *MailRepository {
	return &MailRepository{db: db}
}

func (r *MailRepository) CreateTempAddress(domain string) (*mail_model.TempAddress, error) {
	addrStr := fmt.Sprintf("%s@%s", uuid.New().String()[:10], domain)
	addr := &mail_model.TempAddress{Address: addrStr}
	query := `INSERT INTO temp_addresses (address) VALUES ($1) RETURNING id, created_at`
	err := r.db.QueryRow(query, addr.Address).Scan(&addr.ID, &addr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func (r *MailRepository) FindAddressByString(address string) (*mail_model.TempAddress, error) {
	var addr mail_model.TempAddress
	query := `SELECT id, address, created_at FROM temp_addresses WHERE address = $1`
	err := r.db.Get(&addr, query, address)
	return &addr, err
}

func (r *MailRepository) SaveEmail(email *mail_model.Email) error {
	query := `INSERT INTO emails (address_id, sender, recipients, subject, body)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query,
		email.AddressID,
		email.Sender,
		pq.Array(email.Recipients),
		email.Subject,
		email.Body,
	)
	return err
}

func (r *MailRepository) GetEmailsForAddress(addressID int) ([]mail_model.Email, error) {
	var emails []mail_model.Email
	query := `SELECT id, sender, recipients, subject, body, received_at FROM emails WHERE address_id = $1 ORDER BY received_at DESC`
	err := r.db.Select(&emails, query, addressID)
	return emails, err
}
