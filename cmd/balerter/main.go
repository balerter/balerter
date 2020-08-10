package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	alertManager "github.com/balerter/balerter/internal/alert/manager"
	apiManager "github.com/balerter/balerter/internal/api/manager"
	"github.com/balerter/balerter/internal/config"
	coreStorageManager "github.com/balerter/balerter/internal/corestorage/manager"
	dsManager "github.com/balerter/balerter/internal/datasource/manager"
	"github.com/balerter/balerter/internal/logger"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/modules"
	chartModule "github.com/balerter/balerter/internal/modules/chart"
	httpModule "github.com/balerter/balerter/internal/modules/http"
	"github.com/balerter/balerter/internal/modules/kv"
	logModule "github.com/balerter/balerter/internal/modules/log"
	runtimeModule "github.com/balerter/balerter/internal/modules/runtime"
	"github.com/balerter/balerter/internal/runner"
	scriptsManager "github.com/balerter/balerter/internal/script/manager"
	uploadStorageManager "github.com/balerter/balerter/internal/upload_storage/manager"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

var (
	version = "undefined"

	loggerOptions []zap.Option // for testing purposes
)

const (
	defaultLuaModulesPath = "./?.lua;./modules/?.lua;./modules/?/init.lua"
)

func main() {
	configSource := flag.String("config", "config.yml", "Configuration source. Currently supports only path to yaml file and 'stdin'.")
	logLevel := flag.String("logLevel", "INFO", "Log level. ERROR, WARN, INFO or DEBUG")
	debug := flag.Bool("debug", false, "debug mode")
	once := flag.Bool("once", false, "once run scripts and exit")
	withScript := flag.String("script", "", "ignore all script sources and runs only one script. Meta-tag @ignore will be ignored")

	flag.Parse()

	msg, code := run(*configSource, *logLevel, *debug, *once, *withScript)

	log.Print(msg)
	os.Exit(code)
}

func run( //nolint:gocritic,gocyclo,funlen // Run main application
	configSource,
	logLevel string,
	debug,
	once bool,
	withScript string,
) (string, int) {
	lua.LuaPathDefault = defaultLuaModulesPath

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	wg := &sync.WaitGroup{}

	coreModules := make([]modules.Module, 0)

	if err := validateLogLevel(logLevel); err != nil {
		return err.Error(), 1
	}

	lgr, err := logger.New(logLevel, debug, loggerOptions...)
	if err != nil {
		return fmt.Sprintf("error init zap logger, %v", err), 1
	}

	metrics.SetVersion(version)

	lgr.Logger().Info("balerter start", zap.String("version", version))

	// Configuration
	cfg, err := config.New(configSource)
	if err != nil {
		return fmt.Sprintf("error init config, %v", err), 1
	}
	lgr.Logger().Debug("loaded configuration", zap.Any("config", cfg))

	if cfg.Global.LuaModulesPath != "" {
		lua.LuaPathDefault = cfg.Global.LuaModulesPath
	}

	lgr.Logger().Debug("lua modules path", zap.String("path", lua.LuaPathDefault))

	// Scripts sources
	lgr.Logger().Info("init scripts manager")
	scriptsMgr := scriptsManager.New()

	if withScript != "" {
		lgr.Logger().Info("rewrite script sources configuration", zap.String("filename", withScript))
		cfg.Scripts.Sources = config.ScriptsSources{
			File: []*config.ScriptSourceFile{
				{
					Name:          "cli-script",
					Filename:      withScript,
					DisableIgnore: true,
				},
			},
		}
	}

	if err = scriptsMgr.Init(cfg.Scripts.Sources); err != nil {
		return fmt.Sprintf("error init scripts manager, %v", err), 1
	}

	// datasources
	lgr.Logger().Info("init datasources manager")
	dsMgr := dsManager.New(lgr.Logger())
	if err = dsMgr.Init(cfg.DataSources); err != nil {
		return fmt.Sprintf("error init datasources manager, %v", err), 1
	}

	// upload storages
	lgr.Logger().Info("init upload storages manager")
	uploadStoragesMgr := uploadStorageManager.New(lgr.Logger())
	if err = uploadStoragesMgr.Init(cfg.Storages.Upload); err != nil {
		return fmt.Sprintf("error init upload storages manager, %v", err), 1
	}

	// core storages
	lgr.Logger().Info("init core storages manager")
	coreStoragesMgr, err := coreStorageManager.New(cfg.Storages.Core, lgr.Logger())
	if err != nil {
		return fmt.Sprintf("error create core storages manager, %v", err), 1
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
		return fmt.Sprintf("error get core storages engine for alert '%s', %v", cfg.Global.Storages.Alert, err), 1
	}
	alertMgr := alertManager.New(alertManagerStorageEngine, lgr.Logger())
	if err = alertMgr.Init(cfg.Channels); err != nil {
		return fmt.Sprintf("error init alert manager, %v", err), 1
	}
	coreModules = append(coreModules, alertMgr)

	if len(cfg.Global.SendStartNotification) > 0 {
		alertMgr.Send("", "", "Balerter Start", cfg.Global.SendStartNotification, nil, "")
	}

	for idx := range cfg.Channels.Email {
		if cfg.Channels.Email[idx].Port == "465" && cfg.Channels.Email[idx].Secure == "" {
			lgr.Logger().Info("secure port 465 with ssl for email channel " + cfg.Channels.Email[idx].Name)
		}
	}

	// ---------------------
	// |
	// | Core Modules
	// |
	// | KV
	// |
	kvEngine, err := coreStoragesMgr.Get(cfg.Global.Storages.KV)
	if err != nil {
		return fmt.Sprintf("error get kv storage engine '%s', %v", cfg.Global.Storages.KV, err), 1
	}
	lgr.Logger().Info("init kv storage", zap.String("engine", cfg.Global.Storages.KV))
	kvModule := kv.New(kvEngine)
	coreModules = append(coreModules, kvModule)

	// ---------------------
	// |
	// | API
	// |
	if cfg.Global.API.Address != "" {
		wg.Add(1)
		apis := apiManager.New(cfg.Global.API, alertManagerStorageEngine, kvEngine, lgr.Logger())
		go apis.Run(ctx, ctxCancel, wg)
	}

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
	// | Core Modules
	// |
	// | runtime
	// |
	runtimeMod := runtimeModule.New(logLevel, debug, once, withScript, configSource, lgr.Logger())
	coreModules = append(coreModules, runtimeMod)

	// ---------------------
	// |
	// | Runner
	// |
	lgr.Logger().Info("init runner")
	rnr := runner.New(cfg.Scripts.UpdateInterval, scriptsMgr, dsMgr, uploadStoragesMgr, coreModules, lgr.Logger())

	lgr.Logger().Info("run runner")
	go rnr.Watch(ctx, ctxCancel, wg, once)

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

	wg.Wait()

	dsMgr.Stop()

	if len(cfg.Global.SendStopNotification) > 0 {
		alertMgr.Send("", "", "Balerter Stop", cfg.Global.SendStopNotification, nil, "")
	}

	lgr.Logger().Info("terminate")

	return "", 0
}

func validateLogLevel(level string) error {
	for _, l := range []string{"ERROR", "WARN", "INFO", "DEBUG"} {
		if l == level {
			return nil
		}
	}

	return fmt.Errorf("wrong log level")
}
