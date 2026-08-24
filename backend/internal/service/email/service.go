package email

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/owndangan/backend/internal/config"
	"github.com/rs/zerolog"
	"gopkg.in/gomail.v2"
)

const sendTimeout = 10 * time.Second

var retryBackoff = []time.Duration{time.Minute, 5 * time.Minute}

type Service struct {
	from     string
	fromName string
	log      zerolog.Logger

	timeout time.Duration
	backoff []time.Duration
	sleep   func(time.Duration)
	send    func(to, subject, htmlBody string) error
}

func NewService(cfg config.SMTPConfig, log zerolog.Logger) *Service {
	s := &Service{
		from:     cfg.From,
		fromName: cfg.FromName,
		log:      log,
		timeout:  sendTimeout,
		backoff:  retryBackoff,
		sleep:    time.Sleep,
	}
	s.send = func(to, subject, htmlBody string) error {
		m := gomail.NewMessage()
		m.SetHeader("From", m.FormatAddress(s.from, s.fromName))
		m.SetHeader("To", to)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", htmlBody)
		d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
		return d.DialAndSend(m)
	}
	return s
}

func (s *Service) Send(to, subject, htmlBody string) error {
	done := make(chan error, 1)
	go func() { done <- s.send(to, subject, htmlBody) }()
	select {
	case err := <-done:
		return err
	case <-time.After(s.timeout):
		return errors.New("email send timed out")
	}
}

func (s *Service) SendWithRetry(to, subject, htmlBody string) error {
	var lastErr error
	for attempt := 0; attempt < len(s.backoff)+1; attempt++ {
		if attempt > 0 {
			s.sleep(s.backoff[attempt-1])
		}
		lastErr = s.Send(to, subject, htmlBody)
		if lastErr == nil {
			return nil
		}
		s.log.Warn().Err(lastErr).Int("attempt", attempt+1).Str("to", to).Msg("email send failed")
	}
	return lastErr
}

// SendAsync never blocks or fails the caller; delivery errors are logged.
func (s *Service) SendAsync(to, subject, htmlBody string) {
	go func() {
		if err := s.SendWithRetry(to, subject, htmlBody); err != nil {
			s.log.Error().Err(err).Str("to", to).Str("subject", subject).Msg("email delivery failed permanently")
		}
	}()
}

type WelcomeData struct {
	Name     string
	LoginURL string
}

type PaymentSuccessData struct {
	Name       string
	PlanName   string
	Amount     string
	ExpiryDate string
}

type ExpiryReminderData struct {
	Name       string
	PlanName   string
	ExpiryDate string
	RenewURL   string
}

type ExpiredData struct {
	Name     string
	PlanName string
	RenewURL string
}

type RSVPData struct {
	OwnerName   string
	GuestName   string
	Invitation  string
	Attendance  string
	GuestCount  int
	SubmittedAt string
}

type GuestbookData struct {
	OwnerName   string
	GuestName   string
	Message     string
	Invitation  string
	SubmittedAt string
}

var (
	welcomeTpl = template.Must(template.New("welcome").Parse(
		`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif">
<h1>Halo {{.Name}}!</h1>
<p>Akun Owndangan Anda berhasil dibuat. Mulai buat undangan pernikahan digital pertama Anda sekarang.</p>
<a href="{{.LoginURL}}">Masuk ke Dashboard</a>
</body></html>`))

	paymentSuccessTpl = template.Must(template.New("payment_success").Parse(
		`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif">
<h1>Halo {{.Name}}!</h1>
<p>Pembayaran Anda untuk paket <strong>{{.PlanName}}</strong> sebesar <strong>Rp{{.Amount}}</strong> telah berhasil.</p>
<p>Langganan Anda aktif hingga {{.ExpiryDate}}.</p>
</body></html>`))

	expiryReminderTpl = template.Must(template.New("expiry_reminder").Parse(
		`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif">
<h1>Halo {{.Name}}!</h1>
<p>Langganan <strong>{{.PlanName}}</strong> Anda akan berakhir pada {{.ExpiryDate}}.</p>
<a href="{{.RenewURL}}">Perpanjang Sekarang</a>
</body></html>`))

	expiredTpl = template.Must(template.New("expired").Parse(
		`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif">
<h1>Halo {{.Name}}!</h1>
<p>Langganan <strong>{{.PlanName}}</strong> Anda telah berakhir. Beberapa fitur kini terkunci.</p>
<a href="{{.RenewURL}}">Perpanjang Sekarang</a>
</body></html>`))

	rsvpTpl = template.Must(template.New("rsvp_notification").Parse(
		`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif">
<h1>Halo {{.OwnerName}}!</h1>
<p><strong>{{.GuestName}}</strong> mengirimkan konfirmasi kehadiran untuk undangan <strong>{{.Invitation}}</strong>.</p>
<p>Kehadiran: {{.Attendance}}<br/>Jumlah tamu: {{.GuestCount}}<br/>Waktu: {{.SubmittedAt}}</p>
</body></html>`))

	guestbookTpl = template.Must(template.New("guestbook_notification").Parse(
		`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif">
<h1>Halo {{.OwnerName}}!</h1>
<p><strong>{{.GuestName}}</strong> menulis di buku tamu undangan <strong>{{.Invitation}}</strong>:</p>
<blockquote>{{.Message}}</blockquote>
<p>Waktu: {{.SubmittedAt}}</p>
</body></html>`))
)

func render(t *template.Template, data any) (string, error) {
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	return b.String(), nil
}

func RenderWelcome(d WelcomeData) (string, error) {
	return render(welcomeTpl, d)
}

func RenderPaymentSuccess(d PaymentSuccessData) (string, error) {
	return render(paymentSuccessTpl, d)
}

func RenderExpiryReminder(d ExpiryReminderData) (string, error) {
	return render(expiryReminderTpl, d)
}

func RenderExpired(d ExpiredData) (string, error) {
	return render(expiredTpl, d)
}

func RenderRSVP(d RSVPData) (string, error) {
	return render(rsvpTpl, d)
}

func RenderGuestbook(d GuestbookData) (string, error) {
	return render(guestbookTpl, d)
}
