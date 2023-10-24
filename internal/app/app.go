package app

import (
	"context"

	"github.com/dissipative/opabinia/internal/infra/config"
	"github.com/dissipative/opabinia/internal/infra/logger"
)

const Name = "opabinia"

// App provides a set of methods to perform operations for the application.
type App struct {
	ctx    context.Context
	config *config.Config
	logger *logger.Logger
}

func NewApp() (*App, error) {
	conf, err := config.Load()
	if err != nil {
		return nil, err
	}

	l, err := logger.NewLogger(conf)
	if err != nil {
		return nil, err
	}

	return &App{
		ctx:    context.Background(),
		config: conf,
		logger: l,
	}, nil
}
