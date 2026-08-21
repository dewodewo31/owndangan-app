package email

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func newTestService(send func(to, subject, html string) error) (*Service, *[]time.Duration) {
	var sleeps []time.Duration
	s := &Service{
		from:     "noreply@owndangan.com",
		fromName: "Owndangan",
		log:      zerolog.Nop(),
		timeout:  50 * time.Millisecond,
		backoff:  retryBackoff,
		sleep:    func(d time.Duration) { sleeps = append(sleeps, d) },
		send:     send,
	}
	return s, &sleeps
}

func TestSendWithRetry_GivesUpAfterThreeAttempts(t *testing.T) {
	calls := 0
	s, sleeps := newTestService(func(to, subject, html string) error {
		calls++
		return errors.New("smtp down")
	})

	err := s.SendWithRetry("user@example.com", "subj", "<p>hi</p>")
	if err == nil {
		t.Fatal("want error after exhausting retries")
	}
	if calls != 3 {
		t.Fatalf("attempts = %d, want 3", calls)
	}
	if len(*sleeps) != 2 {
		t.Fatalf("sleeps = %d, want 2", len(*sleeps))
	}
	if (*sleeps)[0] != time.Minute || (*sleeps)[1] != 5*time.Minute {
		t.Fatalf("backoff = %v, want [1m 5m]", *sleeps)
	}
}

func TestSendWithRetry_SucceedsOnSecondAttempt(t *testing.T) {
	calls := 0
	s, _ := newTestService(func(to, subject, html string) error {
		calls++
		if calls == 1 {
			return errors.New("transient")
		}
		return nil
	})

	if err := s.SendWithRetry("user@example.com", "subj", "<p>hi</p>"); err != nil {
		t.Fatalf("want nil error, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("attempts = %d, want 2", calls)
	}
}

func TestSend_TimesOutWhenSenderBlocks(t *testing.T) {
	s, _ := newTestService(func(to, subject, html string) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	if err := s.Send("user@example.com", "subj", "<p>hi</p>"); err == nil {
		t.Fatal("want timeout error")
	}
}

func TestTemplates_EscapeUserContent(t *testing.T) {
	html, err := RenderPaymentSuccess(PaymentSuccessData{
		Name:       `<script>alert(1)</script>`,
		PlanName:   "Premium",
		Amount:     "299.000",
		ExpiryDate: "2026-01-01",
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(html, "<script>") {
		t.Fatal("unescaped script tag in rendered template")
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Fatal("expected escaped user content")
	}
	if !strings.Contains(html, "Premium") {
		t.Fatal("plan name missing from output")
	}
}

func TestRenderAll_NoError(t *testing.T) {
	if _, err := RenderWelcome(WelcomeData{Name: "A", LoginURL: "https://x"}); err != nil {
		t.Fatal(err)
	}
	if _, err := RenderExpiryReminder(ExpiryReminderData{Name: "A", PlanName: "P", ExpiryDate: "d", RenewURL: "u"}); err != nil {
		t.Fatal(err)
	}
	if _, err := RenderExpired(ExpiredData{Name: "A", PlanName: "P", RenewURL: "u"}); err != nil {
		t.Fatal(err)
	}
}
