package logger

import "go.uber.org/zap"

var log *zap.SugaredLogger

func Init(mode string) {
	var baseLogger *zap.Logger
	var err error

	switch mode {
	case "dev":
		baseLogger, err = zap.NewDevelopment()
	default:
		baseLogger, err = zap.NewProduction()
	}

	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	log = baseLogger.Sugar()
}

func L() *zap.SugaredLogger {
	if log == nil {
		Init("dev")
	}
	return log
}
