package database

import (
	"context"
	"fmt"
	"tegami-server/internal/config"
	"tegami-server/migrations"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Setup(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Silent),
		TranslateError: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to setup db: %w", err)
	}
	log.Info().Msg("connected to database")
	err = migrate(db)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	return db, nil
}

func migrate(gormDB *gorm.DB) error {
	db, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}

	res, err := provider.Up(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to migrate database")
		return err
	}

	log.Info().Int("total_amount", len(res)).Msg("migrated successfully")
	return nil
}
