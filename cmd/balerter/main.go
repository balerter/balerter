package main

import (
	"flag"
	"github.com/balerter/balerter/internal/config"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	version = "undefined"
)

func main() {
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Printf("error init zap logger, %v", err)
		os.Exit(1)
	}

	logger.Info("balerter start", zap.String("version", version))

	cfg := config.New()
	if err := cfg.Init(); err != nil {
		logger.Error("error init config", zap.Error(err))
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		logger.Error("error validate config", zap.Error(err))
		os.Exit(1)
	}

	//

	logger.Info("terminate")
}
