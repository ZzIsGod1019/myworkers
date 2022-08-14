package main

import (
	"myworkers/logger"
	"strconv"
)

func main() {
	logger.LoggerInit()

	for i := 0; i < 10240; i++ {
		logger.Trace("info", "en", "hello world"+strconv.Itoa(i))
	}
}
