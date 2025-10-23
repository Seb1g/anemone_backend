package smtp_server

import (
	"anemone_notes/internal/config"
	"anemone_notes/internal/model/mail_model"
	"anemone_notes/internal/repository/mail_repository"
	"encoding/base64"
	"errors"
	"github.com/emersion/go-smtp"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"log"
	"maps"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

type Server struct {
	cfg  *config.Config
	repo *mail_repository.MailRepository
}

func NewServer(cfg *config.Config, repo *mail_repository.MailRepository) *Server {
	return &Server{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *Server) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		repo:   s.repo,
		domain: s.cfg.DomainName,
	}, nil
}

func (s *Server) Start() {
	srv := smtp.NewServer(s)

	srv.Addr = ":" + s.cfg.SMTPPort
	srv.Domain = s.cfg.DomainName
	srv.ReadTimeout = 10 * time.Second
	srv.WriteTimeout = 10 * time.Second
	srv.MaxMessageBytes = 1024 * 1024
	srv.MaxRecipients = 50
	srv.AllowInsecureAuth = true

	log.Printf("INFO: Starting SMTP server at %s for domain %s", srv.Addr, srv.Domain)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("FATAL: Failed to start SMTP server: %v", err)
	}
}

type Session struct {
	repo      *mail_repository.MailRepository
	domain    string
	from      string
	rcptTo    []string
	addressID int
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !strings.HasSuffix(to, "@"+s.domain) {
		log.Printf("SMTP RCPT: domain mismatch. Recipient: %s", to)
		return errors.New("invalid recipient domain")
	}

	addr, err := s.repo.FindAddressByString(to)
	if err != nil {
		log.Printf("SMTP RCPT: address not found: %s", to)
		return errors.New("address does not exist")
	}

	s.rcptTo = append(s.rcptTo, to)
	s.addressID = addr.ID
	return nil
}

func applyDecoding(r io.Reader, headers mail.Header) io.Reader {
	cte := strings.ToLower(headers.Get("Content-Transfer-Encoding"))

	switch cte {
	case "base64":
		// Декодер Base64
		return base64.NewDecoder(base64.StdEncoding, r)
	case "quoted-printable":
		// Декодер Quoted-Printable
		return quotedprintable.NewReader(r)
	default:
		// Неизвестная или 7bit/8bit (без кодирования)
		return r
	}
}

func extractAndSanitizeHTML(msg *mail.Message) (string, error) {
	p := bluemonday.UGCPolicy()

	styleRe := regexp.MustCompile(`^[a-zA-Z0-9\s\:\;\#\(\)\-\,\.%]*$`)
	p.AllowAttrs("style").Matching(styleRe).OnElements("p", "span", "div", "td", "th")

	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)

	bodyReader := applyDecoding(msg.Body, msg.Header)

	if err != nil || (!strings.HasPrefix(mediaType, "multipart/") && mediaType != "text/html") {
		bodyBytes, _ := io.ReadAll(bodyReader)
		return p.Sanitize(string(bodyBytes)), nil
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(bodyReader, params["boundary"])
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading MIME part: %v", err)
				continue
			}

			mailHeader := make(mail.Header)
			maps.Copy(mailHeader, part.Header)

			partMsg := &mail.Message{
				Header: mailHeader,
				Body:   part,
			}

			html, err := extractAndSanitizeHTML(partMsg)
			if err == nil && html != "" {
				return html, nil
			}
		}
	}

	if mediaType == "text/html" {
		bodyBytes, err := io.ReadAll(bodyReader)
		if err != nil {
			return "", err
		}

		cleanHTML := p.Sanitize(string(bodyBytes))
		return cleanHTML, nil
	}

	return "", nil
}

func (s *Session) Data(r io.Reader) error {
	if len(s.rcptTo) == 0 {
		return errors.New("no recipients")
	}

	msg, err := mail.ReadMessage(r)
	if err != nil {
		log.Printf("SMTP DATA: could not read message: %v", err)
		return err
	}

	cleanedHTMLBody, err := extractAndSanitizeHTML(msg)
	if err != nil {
		log.Printf("SMTP DATA: could not extract/sanitize HTML: %v", err)
		return errors.New("failed to process message body")
	}

	subject := msg.Header.Get("Subject")

	newEmail := &mail_model.Email{
		AddressID:  s.addressID,
		Sender:     s.from,
		Recipients: s.rcptTo,
		Subject:    subject,
		Body:       cleanedHTMLBody,
	}

	if err := s.repo.SaveEmail(newEmail); err != nil {
		log.Printf("SMTP DATA: failed to save email for address ID %d: %v", s.addressID, err)
		return errors.New("internal server error")
	}

	log.Printf("SMTP DATA: saved email for %s", strings.Join(s.rcptTo, ", "))
	return nil
}

func (s *Session) Reset() {
	s.from = ""
	s.rcptTo = nil
	s.addressID = 0
}

func (s *Session) Logout() error {
	return nil
}
