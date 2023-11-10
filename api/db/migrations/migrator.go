package migrations

import (
	"context"
	"embed"

	"github.com/fredbi/go-api-skeleton/api/pkg/injected"
	"github.com/fredbi/go-trace/tracer"
	"github.com/fredbi/gooseplus"
)

//go:embed sql/*/*.sql
var embeddedMigrations embed.FS

// Migrator knows how to apply changes (migrations) to
// a versioned database schema.
type Migrator struct {
	rt injected.Iface
	*gooseplus.Migrator
}

func New(rt injected.Iface) *Migrator {
	zlg := rt.Logger().Zap()
	return &Migrator{
		rt: rt,
		Migrator: gooseplus.New(
			rt.DB().DB,
			gooseplus.WithDialect("postgres"),
			gooseplus.WithFS(embeddedMigrations),
			gooseplus.WithLogger(zlg),
			gooseplus.WithGlobalLock(true),
		),
	}
}

func (m Migrator) Migrate(parentCtx context.Context) error {
	ctx, span, _ := tracer.StartSpan(parentCtx, m.rt)
	defer span.End()

	// TODO: allow options in migrate, e.g. override the logger
	return m.Migrator.Migrate(ctx)
}
