package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(ctx context.Context, databaseUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		return nil, fmt.Errorf("open PostgresSQL connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying database pool: %w", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping PostgresSQL: %w", err)
	}

	return db, nil
}
