package utils

import (
	"time"

	"github.com/Sirupsen/logrus"
)

// Time any function in the repository -
// Usage - defer utils.TimeTrack(time.Now(), "Filename.go-FunctionName",log)
func TimeTrack(start time.Time, name string, log *logrus.Logger) {
	elapsed := time.Since(start)
	log.Info("TimeTrack : ", name, " took ", elapsed, "\n")
}
