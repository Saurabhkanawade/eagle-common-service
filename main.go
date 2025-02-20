package main

import (
	"github.com/Saurabhkanawade/eagle-common-service/logwriter"
)

func main() {
	logger := logwriter.NewLoggers()

	logger.Info("Welcome to the eagle-common-service.....")
}
