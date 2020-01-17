package main

import (
	"context"
	"flag"
	"fmt"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"github.com/balerter/balerter/internal/config"
	dsManager "github.com/balerter/balerter/internal/datasource/manager"
	"github.com/balerter/balerter/internal/logger"
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
	logLevel     = flag.String("logLevel", "INFO", "Log level. ERROR, WARN, INFO or DEBUG")
	debug        = flag.Bool("debug", false, "debug mode")
)

func main() {
	flag.Parse()

	if err := validateLogLevel(*logLevel); err != nil {
		log.Print(err)
		os.Exit(1)
	}

	lgr, err := logger.New(*logLevel, *debug)
	if err != nil {
		log.Printf("error init zap logger, %v", err)
		os.Exit(1)
	}

	lgr.Logger().Info("balerter start", zap.String("version", version))

	// Configuration
	cfg := config.New()
	if err := cfg.Init(*configSource); err != nil {
		lgr.Logger().Error("error init config", zap.Error(err))
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		lgr.Logger().Error("error validate config", zap.Error(err))
		os.Exit(1)
	}
	lgr.Logger().Debug("loaded configuration", zap.Any("config", cfg))

	// Scripts sources
	lgr.Logger().Info("init scripts manager")
	scriptsMgr := scriptsManager.New()
	if err := scriptsMgr.Init(cfg.Scripts.Sources); err != nil {
		lgr.Logger().Error("error init scripts manager", zap.Error(err))
		os.Exit(1)
	}

	// datasources
	lgr.Logger().Info("init datasources manager")
	dsMgr := dsManager.New(lgr.Logger())
	if err := dsMgr.Init(cfg.DataSources); err != nil {
		lgr.Logger().Error("error init datasources manager", zap.Error(err))
		os.Exit(1)
	}

	// alert
	lgr.Logger().Info("init alert manager")
	alertMgr := alertManager.New(lgr.Logger())
	if err := alertMgr.Init(cfg.Channels); err != nil {
		lgr.Logger().Error("error init alert manager", zap.Error(err))
		os.Exit(1)
	}

	if err := alertMgr.Send("BAlerter *START*", cfg.Global.SendStartNotification); err != nil {
		lgr.Logger().Error("error send start notification", zap.Error(err))
	}

	// Runner
	lgr.Logger().Info("init runner")
	rnr := runner.New(scriptsMgr, dsMgr, alertMgr, cfg.Scripts.Sources.UpdateInterval, lgr.Logger())

	ctx, ctxCancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	lgr.Logger().Info("run runner")
	go rnr.Watch(ctx, wg)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGTERM)

	var sig os.Signal

	select {
	case sig = <-ch:
		lgr.Logger().Info("got os signal", zap.String("signal", sig.String()))
		ctxCancel()
	case <-ctx.Done():
	}

	rnr.Stop()
	dsMgr.Stop()

	wg.Wait()

	if err := alertMgr.Send("BAlerter *STOP*", cfg.Global.SendStartNotification); err != nil {
		lgr.Logger().Error("error send stop notification", zap.Error(err))
	}

	lgr.Logger().Info("terminate")
}

func validateLogLevel(level string) error {
	if level != "ERROR" && level != "WARN" && level != "INFO" && level != "DEBUG" {
		return fmt.Errorf("wrong log level")
	}

	return nil
}
