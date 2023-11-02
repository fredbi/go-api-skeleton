package commands

import (
	"fmt"

	"github.com/fredbi/go-api-skeleton/api/app"
	configkeys "github.com/fredbi/go-api-skeleton/api/cmd/app-name/commands/config-keys" // CHANGE_ME
	"github.com/fredbi/go-cli/cli"
	"github.com/fredbi/go-cli/cli/injectable"
	"github.com/fredbi/go-cli/config"
	"github.com/fredbi/go-trace/log"
	"github.com/fredbi/go-trace/tracer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// the name of this app.
const appName = "app-name" // CHANGE_ME

func Root() *cli.Command {
	// resolve configuration file, env & CLI flags
	cfg, err := config.LoadWithSecrets("",
		config.WithSearchParentDir(true),
		config.WithMute(true),
	)
	if err != nil {
		cli.Die("unable to load configuration file: %v", err)
	}
	configkeys.SetDefaults(cfg)

	// root structured zap logger
	zlog, closer := log.MustGetLogger(appName,
		log.WithLevel(cfg.GetString("log.level")),
	)
	zlog.Info("starting app", zap.String("app_name", appName))

	return cli.NewCommand(
		&cobra.Command{
			Use:          appName,
			Short:        fmt.Sprintf("Serve API for %s", appName),
			Long:         "expose a REST JSON API over http/https",
			RunE:         root,
			SilenceUsage: true,
			Args:         cobra.NoArgs,
			PersistentPreRun: func(_ *cobra.Command, _ []string) {
				// set a distinctive key in traces
				tracer.RegisterPrefix(appName)
			},
			PersistentPostRun: func(_ *cobra.Command, _ []string) {
				// sync the logger upon exit
				closer()
			},
		},
		cli.WithFlag("log-level", "info", "controls logging verbosity",
			cli.FlagIsPersistent(),
			cli.BindFlagToConfig(configkeys.LogLevel),
		),
		// versioning based on module version
		cli.WithAutoVersion(),
		// inject dependencies in the command's context
		cli.WithConfig(cfg),
		cli.WithInjectables(
			injectable.NewZapLogger(zlog),
		),
	)
}

func root(c *cobra.Command, _ []string) error {
	// resolve dependencies and start the server
	ctx := c.Context()
	zlg := injectable.ZapLoggerFromContext(ctx, zap.NewNop)
	cfg := injectable.ConfigFromContext(ctx, func() *viper.Viper {
		cfg := viper.New()
		configkeys.SetDefaults(cfg)

		return cfg
	})

	server := app.NewServer(appName, zlg, cfg.Sub(configkeys.AppConfig))

	// connect to the database, etc.
	if err := server.Init(); err != nil {
		zlg.Error("an error occured while warming-up this instance",
			zap.String("outcome", "server not started"),
			zap.Error(err),
		)

		return err
	}

	defer func() {
		err := server.Stop()
		if err != nil {
			zlg.Warn("an error occured while stopping this instance",
				zap.Error(err),
			)
		}
	}()

	// serve API endpoints
	return server.Start()
}
