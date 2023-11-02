package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	configkeys "github.com/fredbi/go-api-skeleton/api/app/config-keys"
	"github.com/fredbi/go-api-skeleton/api/pkg/injected"
	"github.com/fredbi/go-api-skeleton/api/pkg/repos"
	"github.com/fredbi/go-api-skeleton/api/pkg/repos/pgrepo"
	"github.com/fredbi/go-trace/log"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var _ injected.Iface = &runtime{}

type (
	// Server is a simple http app server. Simple, but extensible.
	//
	// For a more capable wrapper (with trace & metrics exporters, graceful shutdown...),
	// checkout github.com/casualjim/go-appserver or similar alternatives.
	//
	// We can use this wrapper to initialize other components, e.g. auth middleware etc.
	Server struct {
		*runtime
		appName string
	}

	// runtime holds all dependencies to be injected into the app.
	runtime struct {
		zlg    *zap.Logger
		logger log.Factory
		cfg    *viper.Viper
		db     *sqlx.DB
		repo   repos.RunnableRepo // this version of the runtime works against one single repo
		router chi.Router
	}
)

// NewServer builds a new instance of the API server.
//
// There is a 2-stage warming-up: (i) Init() (ii) Start().
// Stages are separated to be able to extend the capabilities with hot reload.
func NewServer(appName string, logger *zap.Logger, cfg *viper.Viper) *Server {
	host, _ := os.Hostname()

	return &Server{
		runtime: &runtime{
			zlg: logger,
			logger: log.NewFactory(logger).With(
				zap.String("service", appName),
				zap.String("host", host),
			),
			cfg:    cfg,
			router: chi.NewRouter(),
		},
		appName: appName,
	}
}

func (r *runtime) Logger() log.Factory {
	return r.logger
}

func (r *runtime) DB() *sqlx.DB {
	return r.db
}

func (r *runtime) Config() *viper.Viper {
	return r.cfg
}

func (r *runtime) Repos() repos.Iface {
	return r.repo
}

func (r *runtime) App() chi.Router {
	return r.router
}

// Init warms up the server by preparing things such as
// database connection pools, dependent API clients, etc.
func (s *Server) Init() error {
	lg := s.logger.Bg()

	if err := s.migrateDB(); err != nil {
		// NOTE: this should only run short-lived migration jobs.
		//
		// For heavier DB operations, use a different lane, in order
		// to avoid timeouts during deployment.
		lg.Error("could not apply DB migrations", zap.Error(err))

		return err
	}

	if err := s.initDB(); err != nil {
		lg.Error("could not initialize DB connection pool", zap.Error(err))

		return err
	}

	if err := s.initRouter(); err != nil {
		lg.Error("could not initialize API router", zap.Error(err))

		return err
	}

	// other initializations...

	return nil
}

func (s *Server) migrateDB() error {
	return nil // TODO
}

func (s *Server) initDB() error {
	lg := s.logger.Bg()
	repo := pgrepo.NewRepository(s.appName, s.zlg, s.logger, s.cfg)

	if err := repo.Start(); err != nil {
		lg.Error("could not connect to the DB", zap.Error(err))

		return err
	}

	s.repo = repo

	return nil
}

func (s *Server) initRouter() error {
	// apply middleware etc here
	s.router.Use(log.Requests(s.logger))

	// register handlers
	// TODO

	// handlers.Register(s.runtime)
	return nil
}

// Start the API server.
//
// This version starts an http server.
//
// It can be easily extended to start multiple listeners
// and suppot multiple protocols (e.g. grpc, ...).
//
// TODO: https configuration
func (s *Server) Start() error {
	s.router.Route("/", func(r chi.Router) {
		// TODO
	})

	scheme := s.cfg.GetString(configkeys.Scheme)

	if scheme != "http" {
		return fmt.Errorf("unsupported scheme. At the moment, only %v is supported, but got %q", "http", scheme)
	}

	host := s.cfg.GetString(configkeys.Host)
	port := s.cfg.GetUint(configkeys.Port)
	hostPort := fmt.Sprintf("%s:%d", host, port)

	s.logger.Bg().Info("listening",
		zap.String("scheme", scheme),
		zap.String("addr", hostPort),
	)

	// TODO: craft more elaborate server, with various timeouts set properly
	return http.ListenAndServe(hostPort, s.router) //#nosec
}

// Stop the API server.
func (s *Server) Stop() error {
	var failed bool
	jointErr := errors.New("error while stopping the API server")

	if err := s.stopDB(); err != nil {
		failed = true
		jointErr = errors.Join(jointErr, err)
	}

	// other components to be stopped ...

	if failed {
		return jointErr
	}

	return nil
}

func (s *Server) stopDB() error {
	if s.repo == nil {
		return nil
	}

	return s.repo.Stop()
}
