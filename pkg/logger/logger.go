package logger

import (
	"sync"

	"go.uber.org/zap"
)

type zapAdapter struct {
	*zap.Logger
}

func (za zapAdapter) Sugar() SugaredLoggerInterface {
	return za.Logger.Sugar()
}

// SugaredLoggerInterface описывает методы, используемые в коде
type SugaredLoggerInterface interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Sync() error
}

// LoggerInterface — возвращает SugaredLoggerInterface
// Фабрика для создания базового логгера
type LoggerInterface interface {
	Sugar() SugaredLoggerInterface
}

// newLoggerFunc может быть подменена в тестах
var newLoggerFunc = func(mode string, opts ...zap.Option) (LoggerInterface, error) {
	switch mode {
	case "dev":
		l, err := zap.NewDevelopment(opts...)
		if err != nil {
			return nil, err
		}
		return zapAdapter{l}, nil
	default:
		l, err := zap.NewProduction(opts...)
		if err != nil {
			return nil, err
		}
		return zapAdapter{l}, nil
	}
}

var (
	mu  sync.Mutex
	log SugaredLoggerInterface
)

// Init инициализирует или перезаписывает глобальный логгер
func Init(mode string, opts ...zap.Option) {
	mu.Lock()
	defer mu.Unlock()

	base, err := newLoggerFunc(mode, opts...)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	log = base.Sugar()
}

// L возвращает глобальный логгер, инициализируя его в dev по умолчанию
func L() SugaredLoggerInterface {
	// Double-checked locking to avoid deadlock
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
