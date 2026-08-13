package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/repository"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var sharedDB *gorm.DB

func TestMain(m *testing.M) {
	dsn := "host=localhost port=5433 user=postgres password=password dbname=owndangan_test sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		os.Exit(0)
	}

	sqlDB, _ := db.DB()
	sqlDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";")

	if err := db.AutoMigrate(
		&model.User{}, &model.RefreshToken{}, &model.Package{}, &model.Transaction{},
		&model.Subscription{}, &model.Template{}, &model.Music{}, &model.Event{},
		&model.EventSection{}, &model.Guest{}, &model.RSVP{}, &model.GuestbookMessage{},
		&model.DigitalGift{}, &model.GalleryPhoto{}, &model.AuditLog{},
		&model.AnalyticsEvent{}, &model.WebhookIdempotency{},
	); err != nil {
		os.Exit(1)
	}

	sharedDB = db
	os.Exit(m.Run())
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	return sharedDB
}

func TestUserRepository_UniqueEmail(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	repo := repository.NewUserRepository(db)

	user1 := &model.User{Name: "User1", Email: "test@example.com", PasswordHash: "hash", Role: "user", Status: "active"}
	require.NoError(t, repo.Create(ctx, user1))

	user2 := &model.User{Name: "User2", Email: "test@example.com", PasswordHash: "hash", Role: "user", Status: "active"}
	err := repo.Create(ctx, user2)
	require.Error(t, err, "duplicate email should fail")
}

func TestUserRepository_ForeignKeyConstraint(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	pkg := &model.Package{Name: "Free", Code: "free_unique", Price: 0, GuestLimit: &[]int{50}[0], TemplateGroup: "standard", Features: datatypes.JSON("{}"), IsActive: true}
	require.NoError(t, db.WithContext(ctx).Create(pkg).Error)

	sub := &model.Subscription{UserID: uuid.New(), PackageID: pkg.ID, Status: "active", StartAt: time.Now(), ExpiresAt: &time.Time{}}
	err := db.WithContext(ctx).Create(sub).Error
	require.Error(t, err, "foreign key constraint should prevent creating subscription for non-existent user")
}

func TestPackageRepository_UniqueCode(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	repo := repository.NewPackageRepository(db)

	pkg1 := &model.Package{Name: "Pkg1", Code: "code1", Price: 10000, GuestLimit: &[]int{50}[0], TemplateGroup: "standard", Features: datatypes.JSON("{}"), IsActive: true}
	require.NoError(t, repo.Create(ctx, pkg1))

	pkg2 := &model.Package{Name: "Pkg2", Code: "code1", Price: 10000, GuestLimit: &[]int{50}[0], TemplateGroup: "standard", Features: datatypes.JSON("{}"), IsActive: true}
	err := repo.Create(ctx, pkg2)
	require.Error(t, err, "duplicate code should fail")
}

func TestEventRepository_UniqueSlug(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	repo := repository.NewEventRepository(db)

	user := &model.User{Name: "User", Email: "event@test.com", PasswordHash: "h", Role: "user", Status: "active"}
	require.NoError(t, db.WithContext(ctx).Create(user).Error)

	now1 := time.Now()
	event1 := &model.Event{UserID: user.ID, Title: "Event1", Slug: "my-wedding", GroomName: "A", BrideName: "B", WeddingDate: &now1, Status: "draft"}
	require.NoError(t, repo.Create(ctx, event1))

	now2 := time.Now()
	event2 := &model.Event{UserID: user.ID, Title: "Event2", Slug: "my-wedding", GroomName: "C", BrideName: "D", WeddingDate: &now2, Status: "draft"}
	err := repo.Create(ctx, event2)
	require.Error(t, err, "duplicate slug should fail")
}

func TestTransactionRepository_ForeignKey(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	user := &model.User{Name: "User", Email: "tx@test.com", PasswordHash: "h", Role: "user", Status: "active"}
	require.NoError(t, db.WithContext(ctx).Create(user).Error)

	pkg := &model.Package{Name: "TestPkg", Code: "txcode", Price: 50000, GuestLimit: &[]int{50}[0], TemplateGroup: "standard", Features: datatypes.JSON("{}"), IsActive: true}
	require.NoError(t, db.WithContext(ctx).Create(pkg).Error)

	txn := &model.Transaction{UserID: user.ID, PackageID: pkg.ID, OrderID: "ORDER-001", GrossAmount: 50000, Status: "pending"}
	require.NoError(t, db.WithContext(ctx).Create(txn).Error)

	var got model.Transaction
	require.NoError(t, db.WithContext(ctx).Preload("Package").Preload("User").Where("order_id = ?", "ORDER-001").First(&got).Error)
	require.Equal(t, "tx@test.com", got.User.Email)
	require.Equal(t, "TestPkg", got.Package.Name)
}

func TestDatabase_IndexesExist(t *testing.T) {
	db := setupTestDB(t)

	var userIdx []string
	db.Raw("SELECT indexname FROM pg_indexes WHERE tablename = 'users' AND indexname LIKE 'idx_users%'").Scan(&userIdx)
	require.NotEmpty(t, userIdx, "users table should have indexes")

	var txnIdx []string
	db.Raw("SELECT indexname FROM pg_indexes WHERE tablename = 'transactions' AND indexname LIKE 'idx_transactions%'").Scan(&txnIdx)
	require.NotEmpty(t, txnIdx, "transactions table should have indexes")

	var eventIdx []string
	db.Raw("SELECT indexname FROM pg_indexes WHERE tablename = 'events' AND indexname LIKE 'idx_events%'").Scan(&eventIdx)
	require.NotEmpty(t, eventIdx, "events table should have indexes")
}

func TestSubscriptionRepository_LinkedToTransaction(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	user := &model.User{Name: "User", Email: "sub@test.com", PasswordHash: "h", Role: "user", Status: "active"}
	require.NoError(t, db.WithContext(ctx).Create(user).Error)

	pkg := &model.Package{Name: "SubPkg", Code: "subcode", Price: 0, GuestLimit: &[]int{10}[0], TemplateGroup: "standard", Features: datatypes.JSON("{}"), IsActive: true}
	require.NoError(t, db.WithContext(ctx).Create(pkg).Error)

	txn := &model.Transaction{UserID: user.ID, PackageID: pkg.ID, OrderID: "ORDER-SUB-001", GrossAmount: 0, Status: "settlement"}
	require.NoError(t, db.WithContext(ctx).Create(txn).Error)

	now := time.Now()
	sub := &model.Subscription{UserID: user.ID, PackageID: pkg.ID, TransactionID: &txn.ID, Status: "active", StartAt: now, ExpiresAt: &now}
	subRepo := repository.NewSubscriptionRepository(db)
	require.NoError(t, subRepo.Create(ctx, sub))

	got, err := subRepo.GetByTransactionID(ctx, txn.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "active", got.Status)
}
