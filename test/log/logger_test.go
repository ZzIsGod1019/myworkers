package test

import (
	"myworkers/logger"
	"strconv"
	"testing"
)

func TestLogger(t *testing.T) {
	logger.LoggerInit()

	for i := 0; i < 10240; i++ {
		logger.Trace("info", "en", "hello world"+strconv.Itoa(i))
	}
}
