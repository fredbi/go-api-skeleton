package keys

import "github.com/spf13/viper"

const (
	Scheme      = "scheme"
	Host        = "host"
	Port        = "port"
	TraceConfig = "trace"
)

func SetDefaults(cfg *viper.Viper) {
	cfg.SetDefault(Scheme, "http")
	cfg.SetDefault(Host, "0.0.0.0")
	cfg.SetDefault(Port, "8080")

	cfg.SetDefault(TraceConfig, map[string]interface{}{
		"enabled": false,
	})
}
