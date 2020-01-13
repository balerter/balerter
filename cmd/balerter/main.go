package main

import (
	"context"
	"flag"
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/runner"
	scriptsManager "github.com/balerter/balerter/internal/script/manager"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	version = "undefined"
)

var (
	configSource = flag.String("config", "config.yml", "Configuration source. Currently supports only path to yaml file.")
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
	if err := cfg.Init(*configSource); err != nil {
		logger.Error("error init config", zap.Error(err))
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		logger.Error("error validate config", zap.Error(err))
		os.Exit(1)
	}
	logger.Debug("loaded configuration", zap.Any("config", cfg))

	logger.Info("init scripts manager")
	scriptsMgr := scriptsManager.New()
	if err := scriptsMgr.Init(cfg.Scripts.Sources); err != nil {
		logger.Error("error init scripts manager", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("init runner")
	rnr := runner.New(scriptsMgr, cfg.Scripts.Sources.UpdateInterval, logger)

	ctx, ctxCancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	logger.Info("run runner")
	go rnr.Watch(ctx, wg)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGTERM)

	var sig os.Signal

	select {
	case sig = <-ch:
		logger.Info("got os signal", zap.String("signal", sig.String()))
		ctxCancel()
	case <-ctx.Done():
	}

	rnr.Stop()

	wg.Wait()

	logger.Info("terminate")
}
