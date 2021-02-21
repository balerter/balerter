package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/balerter/balerter/internal/alert"
	"github.com/balerter/balerter/internal/config/scripts/sources"
	"github.com/balerter/balerter/internal/config/scripts/sources/file"
	"github.com/balerter/balerter/internal/corestorage"
	alertModule "github.com/balerter/balerter/internal/modules/alert"
	"github.com/balerter/balerter/internal/service"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	apiManager "github.com/balerter/balerter/internal/api/manager"
	channelsManager "github.com/balerter/balerter/internal/chmanager"
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

func run(
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
		cfg.Scripts.Sources = sources.Sources{
			File: []*file.File{
				{
					Name:          "cli-script",
					Filename:      withScript,
					DisableIgnore: true,
				},
			},
		}
	}

	// log if use a 465 port and an empty secure string for an email channel
	for idx := range cfg.Channels.Email {
		if cfg.Channels.Email[idx].Port == "465" && cfg.Channels.Email[idx].Secure == "" {
			lgr.Logger().Info("secure port 465 with ssl for email channel " + cfg.Channels.Email[idx].Name)
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
	coreStorageAlert, err := coreStoragesMgr.Get(cfg.Global.Storages.Alert)
	if err != nil {
		return fmt.Sprintf("error get core storage: alert '%s', %v", cfg.Global.Storages.Alert, err), 1
	}
	coreStorageKV, err := coreStoragesMgr.Get(cfg.Global.Storages.KV)
	if err != nil {
		return fmt.Sprintf("error get core storage: kv '%s', %v", cfg.Global.Storages.KV, err), 1
	}

	// ChannelsManager
	lgr.Logger().Info("init channels manager")
	channelsMgr := channelsManager.New(lgr.Logger())
	if err = channelsMgr.Init(cfg.Channels); err != nil {
		return fmt.Sprintf("error init channels manager, %v", err), 1
	}
	// TODO: pass channels manager...

	// ---------------------
	// |
	// | API
	// |
	if cfg.Global.API.Address != "" {
		var ln net.Listener
		ln, err = net.Listen("tcp", cfg.Global.API.Address)
		if err != nil {
			return fmt.Sprintf("error create api listener, %v", err), 1
		}
		apis := apiManager.New(cfg.Global.API, coreStorageAlert, coreStorageKV, channelsMgr, lgr.Logger())
		wg.Add(1)
		go apis.Run(ctx, ctxCancel, wg, ln)
	}

	// ---------------------
	// |
	// | Service
	// |
	if cfg.Global.Service.Address != "" {
		var ln net.Listener
		ln, err = net.Listen("tcp", cfg.Global.Service.Address)
		if err != nil {
			return fmt.Sprintf("error create service listener, %v", err), 1
		}
		srv := service.New(lgr.Logger())
		wg.Add(1)
		go srv.Run(ctx, ctxCancel, wg, ln)
	}

	coreModules := initCoreModules(coreStorageAlert, coreStorageKV, channelsMgr, lgr.Logger(), logLevel, debug, once, withScript, configSource)

	if len(cfg.Global.SendStartNotification) > 0 {
		channelsMgr.Send(nil, "Balerter start", &alert.Options{
			Channels: cfg.Global.SendStartNotification,
		})
	}

	// ---------------------
	// |
	// | Runner
	// |
	lgr.Logger().Info("init runner")
	rnr := runner.New(cfg.Scripts.UpdateInterval, scriptsMgr, dsMgr, uploadStoragesMgr, coreModules, lgr.Logger())

	lgr.Logger().Info("run runner")
	go rnr.Watch(ctx, ctxCancel, once)

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
	coreStoragesMgr.Stop()

	if len(cfg.Global.SendStopNotification) > 0 {
		channelsMgr.Send(nil, "Balerter stop", &alert.Options{
			Channels: cfg.Global.SendStopNotification,
		})
	}

	lgr.Logger().Info("terminate")

	return "", 0
}

func initCoreModules(
	coreStorageAlert corestorage.CoreStorage,
	coreStorageKV corestorage.CoreStorage,
	chManager *channelsManager.ChannelsManager,
	logger *zap.Logger,
	logLevel string,
	debug bool,
	once bool,
	withScript string,
	configSource string,
) []modules.Module {
	coreModules := make([]modules.Module, 0)

	alertMod := alertModule.New(coreStorageAlert.Alert(), chManager, logger)
	coreModules = append(coreModules, alertMod)

	kvModule := kv.New(coreStorageKV.KV())
	coreModules = append(coreModules, kvModule)

	logMod := logModule.New(logger)
	coreModules = append(coreModules, logMod)

	chartMod := chartModule.New(logger)
	coreModules = append(coreModules, chartMod)

	httpMod := httpModule.New(logger)
	coreModules = append(coreModules, httpMod)

	runtimeMod := runtimeModule.New(logLevel, debug, once, withScript, configSource, logger)
	coreModules = append(coreModules, runtimeMod)

	return coreModules
}

func validateLogLevel(level string) error {
	for _, l := range []string{"ERROR", "WARN", "INFO", "DEBUG"} {
		if l == level {
			return nil
		}
	}

	return fmt.Errorf("wrong log level")
}
