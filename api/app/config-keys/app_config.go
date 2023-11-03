package keys

import (
	"time"

	"github.com/spf13/viper"
)

const (
	Scheme          = "scheme"
	Host            = "host"
	Port            = "port"
	TraceConfig     = "trace"
	MigrationConfig = "migrations"

	TraceEnabled     = "enabled"
	MigrationEnabled = "enabled"
	MigrationTimeout = "timeout"
)

func SetDefaults(cfg *viper.Viper) {
	cfg.SetDefault(Scheme, "http")
	cfg.SetDefault(Host, "0.0.0.0")
	cfg.SetDefault(Port, "8080")

	cfg.SetDefault(TraceConfig, map[string]interface{}{
		TraceEnabled: false,
	})
	cfg.SetDefault(MigrationConfig, map[string]interface{}{
		MigrationEnabled: true,
		MigrationTimeout: 10 * time.Minute,
	})
}
