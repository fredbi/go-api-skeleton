package migrations

import (
	"context"
	"embed"
	"errors"
	"path/filepath"

	"github.com/fredbi/go-api-skeleton/api/pkg/injected"
	"github.com/fredbi/go-trace/tracer"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed sql/*/*.sql
var embedMigrations embed.FS

// Migrator knows how to apply changes (migrations) to
// a versioned database schema.
type Migrator struct {
	rt injected.Iface
}

func New(rt injected.Iface) *Migrator {
	return &Migrator{
		rt: rt,
	}
}

func (m Migrator) Migrate(parentCtx context.Context) error {
	ctx, span, lg := tracer.StartSpan(parentCtx, m.rt)
	defer span.End()

	db := m.rt.DB().DB
	goose.SetBaseFS(embedMigrations)
	lg = lg.With(zap.String("module", "dbmigrator"))

	// TODO(fred): testability - add support for test environments
	dir := filepath.Join("sql", "default")

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	current, err := goose.EnsureDBVersionContext(ctx, db)
	if err != nil {
		return errors.Join(errors.New("could not ensure goose migration table"), err)
	}

	if rollForwardErr := goose.UpContext(ctx, db, dir); rollForwardErr != nil {
		// rollback a failed release back to when the deployment started
		lg.Error("failure",
			zap.String("action", "rollbacking to the initial state of deployment"),
			zap.Error(err),
		)

		if rollBackErr := goose.DownToContext(ctx, db, dir, current); rollBackErr != nil {
			lg.Error("encountered again an error while rollbacking",
				zap.String("action", "bailed"),
				zap.String("status", "this might leave your database in an inconsistent state"),
				zap.Error(err),
			)
			return errors.Join(errors.New("irrecoverable error"), err)
		}
	}

	return nil
}
