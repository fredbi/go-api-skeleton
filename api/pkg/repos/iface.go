package repos

import (
	"github.com/jmoiron/sqlx"
)

//go:generate moq -out ./mocks/mocks_runnable.go -pkg mocks . RunnableRepo

// Iface exposes the capabilities of persistent repositories.
//
// Several repos may be exposed: just extend this with additional
// repo factories.
type Iface interface {
	Sample() SampleRepo
}

// Validatable knows how to validate a repo type
type Validatable interface {
	Validate() error
}

// RunnableRepo serves repositories and is manageable by a runtime to start, monitor and stop the repo.
//
// This version of the interface works against one single DB. Adopt DB(key string)
// to dispatch on multiple DBs.
type RunnableRepo interface {
	Iface

	Start() error
	Stop() error
	DB() *sqlx.DB
	// HealthCheck() error
}
