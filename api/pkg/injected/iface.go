package injected

import (
	"github.com/fredbi/go-api-skeleton/api/pkg/repos"
	"github.com/fredbi/go-trace/log"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

//go:generate moq -out ./mocks/runtime_iface.go -pkg mocks . Iface

// Iface is the interface that provides the application context.
//
// This context contains the state that is global to an application instance.
// It provides insulation, so that multiple contexts can actually run in parallel in the same process
//
// Typical extensions of this app environment are: healthcheck instrumentation, authentication & authorization tooling,
// API clients used to process incoming requests, etc.
type Iface interface {
	App() chi.Router
	Logger() log.Factory
	Config() *viper.Viper
	Repos() repos.Iface
	DB() *sqlx.DB
}
