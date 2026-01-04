package commonssmtp

import (
	"io"
	"time"

	smtp "github.com/emersion/go-smtp"
	"github.com/joukojo/go-mail-testserver/internal/httpapi"
)

type SmtpServer struct {
	storage    *httpapi.Storage
	backend    *backend
	SmtpServer *smtp.Server
}

type backend struct {
	store *httpapi.Storage
}

func (b *backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	// Allow any session for local testing.
	return &session{storage: b.store}, nil
}

type session struct {
	storage *httpapi.Storage
	from    string
	to      []string
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	s.to = nil
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = append(s.to, to)
	return nil
}

func (s *session) Data(r io.Reader) error {
	raw, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	s.storage.Add(msgFromRaw(s.from, s.to, raw))
	return nil
}

func msgFromRaw(s1 string, s2 []string, raw []byte) *httpapi.Message {
	return &httpapi.Message{
		From: s1,
		To:   s2,
		Body: string(raw),
		Raw:  raw,
	}
}

func (s *session) Reset()        {}
func (s *session) Logout() error { return nil }

// --- HTTP API ---

func NewSmtpServer(storage *httpapi.Storage, addr string) *SmtpServer {
	be := &backend{store: storage}

	s := smtp.NewServer(be)
	s.Addr = addr
	s.Domain = "localhost"
	s.AllowInsecureAuth = true // OK for local testing only
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 10 << 20 // 10 MB
	s.MaxRecipients = 50
	stmpServer := &SmtpServer{storage: storage, backend: be, SmtpServer: s}
	return stmpServer

}

func (s *SmtpServer) Start() error {
	return s.SmtpServer.ListenAndServe()
}
