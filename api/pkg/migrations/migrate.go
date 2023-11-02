package migrations

import (
	"embed"
	"path/filepath"

	"github.com/fredbi/go-api-skeleton/api/pkg/injected"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*/*.sql
var embedMigrations embed.FS

func Migrate(rt injected.Iface) error {
	db := rt.DB().DB

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, filepath.Join("sql", "default")); err != nil {
		// TODO: a more convoluted operation is need to rollback a failed release,
		// back to when the deployment started.
		return err
	}

	return nil
}
