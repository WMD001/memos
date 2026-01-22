package vercel

import (
	"context"
	"log/slog"
	"os"

	"github.com/pkg/errors"

	"github.com/usememos/memos/internal/profile"
	"github.com/usememos/memos/store"
	"github.com/usememos/memos/store/db/mysql"
	"github.com/usememos/memos/store/db/postgres"
	"github.com/usememos/memos/store/db/sqlite"
)

func NewVercelDBDriver(ctx context.Context) (store.Driver, error) {
	profileConfig := &profile.Profile{
		Driver: getEnv("MEMOS_DRIVER", "sqlite"),
		DSN:    getEnv("MEMOS_DSN", ""),
	}

	if profileConfig.DSN == "" && profileConfig.Driver == "sqlite" {
		profileConfig.DSN = "/tmp/memos_prod.db"
	}

	if err := profileConfig.Validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate profile")
	}

	var driver store.Driver
	var err error

	switch profileConfig.Driver {
	case "sqlite":
		driver, err = sqlite.NewDB(profileConfig)
	case "mysql":
		driver, err = mysql.NewDB(profileConfig)
	case "postgres":
		driver, err = postgres.NewDB(profileConfig)
	default:
		return nil, errors.New("unknown db driver")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to create db driver")
	}

	slog.Info("Database initialized", "driver", profileConfig.Driver)
	return driver, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
