package keys

import "github.com/spf13/viper"

const (
	LogLevel  = "log-level"
	AppConfig = "app"
)

// SetDefaults for CLI settings
func SetDefaults(cfg *viper.Viper) {
	cfg.SetDefault(AppConfig, make(map[string]interface{}))
	cfg.SetDefault(LogLevel, "info")
}
