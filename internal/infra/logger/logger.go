package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

type Logger struct {
	*slog.Logger
	level slog.Level
}

type Config interface {
	LogLevel() []byte
	IsDevEnv() bool
}

func NewLogger(conf Config) (*Logger, error) {
	var lvl slog.Level

	err := lvl.UnmarshalText(conf.LogLevel())
	if err != nil {
		return nil, err
	}

	if conf.IsDevEnv() {
		return &Logger{
			slog.New(tint.NewHandler(os.Stdout, &tint.Options{
				Level:      lvl,
				TimeFormat: time.TimeOnly,
			})),
			lvl,
		}, nil
	}

	return &Logger{
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})),
		lvl,
	}, nil
}

func (l *Logger) IsDebug() bool {
	return l.level == slog.LevelDebug
}
