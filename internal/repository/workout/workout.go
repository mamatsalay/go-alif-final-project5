package workout

import (
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type WorkoutRepositoryParams struct {
	dig.In

	Log *zap.SugaredLogger
}
