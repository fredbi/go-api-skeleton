package keys

import (
	"time"

	"github.com/spf13/viper"
)

const (
	PGMaxIdleConns    = "config.maxIdleConns"
	PGMaxOpenConns    = "config.maxOpenConns"
	PGConnMaxLifetime = "config.connMaxLifetime"
	PGLogLevel        = "config.log.level"
	PGTraceEnabled    = "config.trace.enabled"
	PGURL             = "url"
	PGReplicas        = "replicas"
	PGUser            = "user"
	PGPassword        = "password"
	PGSet             = "config.set" //	plan_cache_mode: auto|force_custom_plan|force_generic_plan
	PGPingTimeout     = "config.pingTimeout"

	DefaultPGLogLevel = "info"
)

func SetDefaults(cfg *viper.Viper) {
	cfg.SetDefault(PGMaxIdleConns, 25)
	cfg.SetDefault(PGMaxOpenConns, 50)
	cfg.SetDefault(PGConnMaxLifetime, "5m")
	cfg.SetDefault(PGLogLevel, DefaultPGLogLevel)
	cfg.SetDefault(PGTraceEnabled, false)
	cfg.SetDefault(PGPingTimeout, 10*time.Second)
}
