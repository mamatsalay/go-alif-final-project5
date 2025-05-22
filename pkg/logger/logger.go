package logger

import "go.uber.org/zap"

var log *zap.Logger

func Init(mode string) {
	var err error

	if mode == "dev" {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}

	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

func L() *zap.Logger {
	if log == nil {
		Init("dev")
	}

	return log
}

func Sugar() *zap.SugaredLogger {
	return L().Sugar()
}
