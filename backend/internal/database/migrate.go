package database

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/owndangan/backend/internal/config"
	"github.com/owndangan/backend/internal/model"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type migration struct {
	version int64
	name    string
	upSQL   string
	downSQL string
}

var migrationRegex = regexp.MustCompile(`^(\d+)_(.+)\.up\.sql$`)

func RunMigrations(cfg config.DatabaseConfig, log zerolog.Logger) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	if err := ensureMigrationTable(db); err != nil {
		return fmt.Errorf("ensure migration table: %w", err)
	}

	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}

	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

	for _, m := range migrations {
		var applied int
		err := tx.Raw("SELECT 1 FROM migration_schema_versions WHERE version = ?", m.version).Row().Scan(&applied)
		if err == nil {
			continue
		}

		for _, stmt := range splitStatements(m.upSQL) {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			log.Debug().Str("stmt", stmt[:min(len(stmt), 100)]).Msg("executing statement")
			if err := tx.Exec(stmt).Error; err != nil {
				log.Error().Err(err).Str("stmt", stmt[:min(len(stmt), 200)]).Msg("migration statement failed")
				return fmt.Errorf("apply migration %d (%s): %w", m.version, m.name, err)
			}
		}

		if err := tx.Exec("INSERT INTO migration_schema_versions (version, dirty) VALUES (?, false)", m.version).Error; err != nil {
			return fmt.Errorf("record migration %d: %w", m.version, err)
		}

		log.Info().Int64("version", m.version).Str("name", m.name).Msg("migration applied")
	}

	log.Info().Int("count", len(migrations)).Msg("all migrations applied")
	return nil
}

func ensureMigrationTable(db *gorm.DB) error {
	return db.Exec("CREATE TABLE IF NOT EXISTS migration_schema_versions (version bigint primary key, dirty boolean)").Error
}

func splitStatements(sqlText string) []string {
	var statements []string
	var buf strings.Builder
	inDollar := false
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(sqlText); i++ {
		ch := sqlText[i]
		if inDollar {
			buf.WriteByte(ch)
			if i+1 < len(sqlText) && sqlText[i+1] == '$' {
				buf.WriteByte(sqlText[i+1])
				inDollar = false
				i++
			}
			continue
		}
		if ch == '$' && i+1 < len(sqlText) && sqlText[i+1] == '$' {
			buf.WriteString("$$")
			inDollar = true
			i++
			continue
		}
		if inSingleQuote {
			buf.WriteByte(ch)
			if ch == '\'' && !(i > 0 && sqlText[i-1] == '\\') {
				inSingleQuote = false
			}
			continue
		}
		if ch == '\'' {
			buf.WriteByte(ch)
			inSingleQuote = true
			continue
		}
		if inDoubleQuote {
			buf.WriteByte(ch)
			if ch == '"' {
				inDoubleQuote = false
			}
			continue
		}
		if ch == '"' {
			buf.WriteByte(ch)
			inDoubleQuote = true
			continue
		}
		if ch == ';' {
			statements = append(statements, buf.String())
			buf.Reset()
			continue
		}
		buf.WriteByte(ch)
	}
	if buf.Len() > 0 {
		statements = append(statements, buf.String())
	}
	return statements
}

func RollbackMigration(cfg config.DatabaseConfig, log zerolog.Logger) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	if err := ensureMigrationTable(db); err != nil {
		return fmt.Errorf("ensure migration table: %w", err)
	}

	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}
	if len(migrations) == 0 {
		return nil
	}

	last := migrations[len(migrations)-1]
	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	for _, stmt := range splitStatements(last.downSQL) {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := tx.Exec(stmt).Error; err != nil {
			log.Warn().Err(err).Msg("rollback statement error")
		}
	}
	tx.Exec("DELETE FROM migration_schema_versions WHERE version = ?", last.version)
	log.Info().Int64("version", last.version).Str("name", last.name).Msg("migration rolled back")
	return nil
}

func loadMigrations() ([]migration, error) {
	dirs := []string{"migrations", "../migrations", "../../migrations"}
	var allMigrations []migration

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
				continue
			}

			matches := migrationRegex.FindStringSubmatch(entry.Name())
			if len(matches) != 3 {
				continue
			}

			content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				continue
			}

			version, _ := strconv.ParseInt(matches[1], 10, 64)
			downPath := filepath.Join(dir, matches[1]+"_"+matches[2]+".down.sql")
			downContent, _ := os.ReadFile(downPath)

			allMigrations = append(allMigrations, migration{
				version: version,
				name:    matches[2],
				upSQL:   stripGooseMarkers(string(content)),
				downSQL: stripGooseMarkers(string(downContent)),
			})
		}
	}

	sort.Slice(allMigrations, func(i, j int) bool {
		return allMigrations[i].version < allMigrations[j].version
	})

	if len(allMigrations) == 0 {
		return nil, fmt.Errorf("no migration files found in any directory")
	}

	seen := make(map[int64]bool)
	var deduped []migration
	for _, m := range allMigrations {
		if !seen[m.version] {
			seen[m.version] = true
			deduped = append(deduped, m)
		}
	}

	return deduped, nil
}

func stripGooseMarkers(sqlText string) string {
	lines := strings.Split(sqlText, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "-- +goose") {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func AutoMigrateDB(cfg config.DatabaseConfig, log zerolog.Logger) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Package{},
		&model.Transaction{},
		&model.Subscription{},
		&model.Template{},
		&model.Music{},
		&model.Event{},
		&model.EventSection{},
		&model.Guest{},
		&model.RSVP{},
		&model.GuestbookMessage{},
		&model.DigitalGift{},
		&model.GalleryPhoto{},
		&model.LoveStory{},
		&model.AuditLog{},
		&model.AnalyticsEvent{},
		&model.WebhookIdempotency{},
	); err != nil {
		return err
	}

	log.Info().Msg("auto-migration completed")
	return nil
}
