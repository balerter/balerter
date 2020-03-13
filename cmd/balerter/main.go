package main

import (
	"context"
	"flag"
	"fmt"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	apiManager "github.com/balerter/balerter/internal/api/manager"
	"github.com/balerter/balerter/internal/config"
	coreStorageManager "github.com/balerter/balerter/internal/core_storage/manager"
	dsManager "github.com/balerter/balerter/internal/datasource/manager"
	"github.com/balerter/balerter/internal/logger"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/modules"
	chartModule "github.com/balerter/balerter/internal/modules/chart"
	httpModule "github.com/balerter/balerter/internal/modules/http"
	"github.com/balerter/balerter/internal/modules/kv"
	logModule "github.com/balerter/balerter/internal/modules/log"
	"github.com/balerter/balerter/internal/runner"
	scriptsManager "github.com/balerter/balerter/internal/script/manager"
	uploadStorageManager "github.com/balerter/balerter/internal/upload_storage/manager"
	lua "github.com/yuin/gopher-lua"
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
	LuaLDir := "./modules"
	lua.LuaPathDefault = "./?.lua;" + LuaLDir + "/?.lua;" + LuaLDir + "/?/init.lua"

	ctx, ctxCancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	coreModules := make([]modules.Module, 0)

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

	metrics.SetVersion(version)

	lgr.Logger().Info("balerter start", zap.String("version", version))

	// Configuration
	cfg := config.New()
	if err := cfg.Init(*configSource); err != nil {
		lgr.Logger().Error("error init config", zap.Error(err))
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

	// upload storages
	lgr.Logger().Info("init upload storages manager")
	uploadStoragesMgr := uploadStorageManager.New(lgr.Logger())
	if err := uploadStoragesMgr.Init(cfg.Storages.Upload); err != nil {
		lgr.Logger().Error("error init upload storages manager", zap.Error(err))
		os.Exit(1)
	}

	// core storages
	lgr.Logger().Info("init core storages manager")
	coreStoragesMgr, err := coreStorageManager.New(cfg.Storages.Core, lgr.Logger())
	if err != nil {
		lgr.Logger().Error("error create core storages manager", zap.Error(err))
		os.Exit(1)
	}

	// ---------------------
	// |
	// | Core Modules
	// |
	// | AlertManager
	// |
	lgr.Logger().Info("init alert manager")

	alertManagerStorageEngine, err := coreStoragesMgr.Get(cfg.Global.Storages.Alert)
	if err != nil {
		lgr.Logger().Error("error get core storages engine for alert", zap.String("name", cfg.Global.Storages.Alert), zap.Error(err))
		os.Exit(1)
	}
	alertMgr := alertManager.New(alertManagerStorageEngine, lgr.Logger())
	if err := alertMgr.Init(cfg.Channels); err != nil {
		lgr.Logger().Error("error init alert manager", zap.Error(err))
		os.Exit(1)
	}
	coreModules = append(coreModules, alertMgr)

	if len(cfg.Global.SendStartNotification) > 0 {
		alertMgr.Send("", "", "Balerter Start", cfg.Global.SendStartNotification, nil, "")
	}

	// ---------------------
	// |
	// | API
	// |
	wg.Add(1)
	apis := apiManager.New(cfg.Global.API, alertMgr, lgr.Logger())
	go apis.Run(ctx, ctxCancel, wg)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | KV
	// |
	kvEngine, err := coreStoragesMgr.Get(cfg.Global.Storages.KV)
	if err != nil {
		lgr.Logger().Error("error get kv storage engine", zap.String("name", cfg.Global.Storages.KV), zap.Error(err))
		os.Exit(1)
	}
	lgr.Logger().Info("init kv storage", zap.String("engine", cfg.Global.Storages.KV))
	kvModule := kv.New(kvEngine)
	coreModules = append(coreModules, kvModule)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | Log
	// |
	logMod := logModule.New(lgr.Logger())
	coreModules = append(coreModules, logMod)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | Chart
	// |
	chartMod := chartModule.New(lgr.Logger())
	coreModules = append(coreModules, chartMod)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | http
	// |
	httpMod := httpModule.New(lgr.Logger())
	coreModules = append(coreModules, httpMod)

	// ---------------------
	// |
	// | Runner
	// |
	lgr.Logger().Info("init runner")
	rnr := runner.New(scriptsMgr, dsMgr, uploadStoragesMgr, coreModules, lgr.Logger())

	lgr.Logger().Info("run runner")
	go rnr.Watch(ctx, wg)

	// ---------------------
	// |
	// | Shutdown
	// |
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

	if len(cfg.Global.SendStopNotification) > 0 {
		alertMgr.Send("", "", "Balerter Stop", cfg.Global.SendStopNotification, nil, "")
	}

	lgr.Logger().Info("terminate")
}

func validateLogLevel(level string) error {
	for _, l := range []string{"ERROR", "WARN", "INFO", "DEBUG"} {
		if l == level {
			return nil
		}
	}

	return fmt.Errorf("wrong log level")
}
