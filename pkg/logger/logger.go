package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type zapAdapter struct {
	*zap.Logger
}

func (za zapAdapter) Sugar() SugaredLoggerInterface {
	return za.Logger.Sugar()
}

type SugaredLoggerInterface interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Sync() error
}

type LoggerInterface interface {
	Sugar() SugaredLoggerInterface
}

var newLoggerFunc = func(mode string, opts ...zap.Option) (LoggerInterface, error) {
	switch mode {
	case "dev":
		l, err := zap.NewDevelopment(opts...)
		if err != nil {
			return nil, fmt.Errorf("error creating new zap logger: %w", err)
		}
		return zapAdapter{l}, nil
	default:
		l, err := zap.NewProduction(opts...)
		if err != nil {
			return nil, fmt.Errorf("error creating new zap logger: %w", err)
		}
		return zapAdapter{l}, nil
	}
}

var (
	mu  sync.Mutex
	log SugaredLoggerInterface
)

func Init(mode string, opts ...zap.Option) {
	mu.Lock()
	defer mu.Unlock()

	base, err := newLoggerFunc(mode, opts...)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	log = base.Sugar()
}

func L() SugaredLoggerInterface {
	mu.Lock()
	current := log
	mu.Unlock()
	if current == nil {
		Init("dev")
		mu.Lock()
		current = log
		mu.Unlock()
	}
	return current
}
