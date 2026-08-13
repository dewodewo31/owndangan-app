package perf

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func BenchmarkPublicInvitation(b *testing.B) {
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		b.Skipf("db not available: %v", err)
	}

	repo := repository.NewEventRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetBySlug(ctx, "test-wedding")
	}
}

func BenchmarkListEvents(b *testing.B) {
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		b.Skipf("db not available: %v", err)
	}

	repo := repository.NewEventRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = repo.ListByUser(ctx, model.User{}.ID, 1, 20, "")
	}
}

func BenchmarkCreateGuest(b *testing.B) {
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		b.Skipf("db not available: %v", err)
	}

	repo := repository.NewGuestRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.Create(ctx, &model.Guest{
			EventID:  model.User{}.ID,
			Name:     fmt.Sprintf("Guest %d", i),
			Category: "family",
			Token:    fmt.Sprintf("token%d", i),
		})
	}
}

func BenchmarkLogin(b *testing.B) {
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		b.Skipf("db not available: %v", err)
	}

	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetByEmail(ctx, "test@example.com")
	}
}

func TestConcurrentRSVP(t *testing.T) {
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("db not available: %v", err)
	}

	repo := repository.NewRSVPRepository(db)
	ctx := context.Background()

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			err := repo.Create(ctx, &model.RSVP{
				GuestID:    model.User{}.ID,
				EventID:    model.User{}.ID,
				Attendance: "attending",
				GuestCount: 1,
			})
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	errCount := 0
	for err := range errors {
		t.Logf("Error: %v", err)
		errCount++
	}

	t.Logf("Concurrent RSVP: %d errors out of 100", errCount)
}

func TestDatabaseConnectionPool(t *testing.T) {
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("db not available: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = db.WithContext(ctx).Exec("SELECT 1")
		}()
	}

	wg.Wait()
	duration := time.Since(start)
	t.Logf("50 concurrent queries completed in %v", duration)
}
