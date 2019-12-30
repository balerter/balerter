package main

import (
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	version = "undefined"
)

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Printf("error init zap logger, %v", err)
		os.Exit(1)
	}

	logger.Info("balerter start", zap.String("version", version))

	//

	logger.Info("terminate")
}
