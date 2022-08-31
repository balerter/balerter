package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	apiManager "github.com/balerter/balerter/internal/api/manager"
	channelsManager "github.com/balerter/balerter/internal/chmanager"
	"github.com/balerter/balerter/internal/cloud"
	"github.com/balerter/balerter/internal/config"
	"github.com/balerter/balerter/internal/coreapi"
	"github.com/balerter/balerter/internal/corestorage"
	coreStorageManager "github.com/balerter/balerter/internal/corestorage/manager"
	dsManager "github.com/balerter/balerter/internal/datasource/manager"
	"github.com/balerter/balerter/internal/logger"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/modules"
	alertModule "github.com/balerter/balerter/internal/modules/alert"
	chartModule "github.com/balerter/balerter/internal/modules/chart"
	"github.com/balerter/balerter/internal/modules/file"
	httpModule "github.com/balerter/balerter/internal/modules/http"
	"github.com/balerter/balerter/internal/modules/kv"
	logModule "github.com/balerter/balerter/internal/modules/log"
	"github.com/balerter/balerter/internal/modules/meta"
	runtimeModule "github.com/balerter/balerter/internal/modules/runtime"
	tlsModule "github.com/balerter/balerter/internal/modules/tls"
	"github.com/balerter/balerter/internal/runner"
	scriptsManager "github.com/balerter/balerter/internal/script/manager"
	"github.com/balerter/balerter/internal/service"
	uploadStorageManager "github.com/balerter/balerter/internal/upload_storage/manager"

	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

var (
	version = "undefined"

	loggerOptions []zap.Option // for testing purposes
)

const (
	defaultLuaModulesPath = "./?.lua;./modules/?.lua;./modules/?/init.lua;/modules/?.lua;/modules/?/init.lua"
)

func main() {
	fs := flag.NewFlagSet("fs", flag.ContinueOnError)
	cfg, flg, err := config.New(fs, os.Args[1:])
	if err != nil {
		log.Printf("error configuration load, %v", err)
		os.Exit(1)
	}

	msg, code := run(cfg, flg)

	log.Print(msg)
	os.Exit(code)
}

func run(
	cfg *config.Config,
	flg *config.Flags,
) (string, int) {
	lua.LuaPathDefault = defaultLuaModulesPath

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	wg := &sync.WaitGroup{}

	if err := validateLogLevel(flg.LogLevel); err != nil {
		return err.Error(), 1
	}

	lgr, err := logger.New(flg.LogLevel, flg.Debug, loggerOptions...)
	if err != nil {
		return fmt.Sprintf("error init zap logger, %v", err), 1
	}

	metrics.SetVersion(version)

	lgr.Logger().Info("balerter start", zap.String("version", version))

	lgr.Logger().Debug("loaded configuration", zap.Any("config", cfg), zap.Any("flags", flg))

	if cfg.LuaModulesPath != "" {
		lua.LuaPathDefault = cfg.LuaModulesPath
	}

	lgr.Logger().Debug("lua modules path", zap.String("path", lua.LuaPathDefault))

	if cfg.Cloud != nil {
		cloud.Init(cfg.Cloud, version, lgr.Logger())

		cloud.SendStart()
		defer cloud.SendStop()
	}

	// Scripts sources
	lgr.Logger().Info("init scripts manager")
	scriptsMgr := scriptsManager.New()

	// log if use a 465 port and an empty secure string for an email channel
	if cfg.Channels != nil {
		for idx := range cfg.Channels.Email {
			if cfg.Channels.Email[idx].Port == "465" && cfg.Channels.Email[idx].Secure == "" {
				lgr.Logger().Info("secure port 465 with ssl for email channel " + cfg.Channels.Email[idx].Name)
			}
		}
	}

	if err = scriptsMgr.Init(cfg.Scripts); err != nil {
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
	if err = uploadStoragesMgr.Init(cfg.StoragesUpload); err != nil {
		return fmt.Sprintf("error init upload storages manager, %v", err), 1
	}

	// core storages
	lgr.Logger().Info("init core storages manager")
	coreStoragesMgr, err := coreStorageManager.New(cfg.StoragesCore, lgr.Logger())
	if err != nil {
		return fmt.Sprintf("error create core storages manager, %v", err), 1
	}
	coreStorageAlert, err := coreStoragesMgr.Get(cfg.StorageAlert)
	if err != nil {
		return fmt.Sprintf("error get core storage: alert '%s', %v", cfg.StorageAlert, err), 1
	}
	coreStorageKV, err := coreStoragesMgr.Get(cfg.StorageKV)
	if err != nil {
		return fmt.Sprintf("error get core storage: kv '%s', %v", cfg.StorageKV, err), 1
	}

	// ChannelsManager
	lgr.Logger().Info("init channels manager")
	channelsMgr := channelsManager.New(lgr.Logger())
	if err = channelsMgr.Init(cfg.Channels, version); err != nil {
		return fmt.Sprintf("error init channels manager, %v", err), 1
	}

	coreModules := initCoreModules(coreStorageAlert, coreStorageKV, channelsMgr, lgr.Logger(), flg)

	if cfg.API != nil && cfg.API.CoreApi != nil && cfg.API.CoreApi.Address != "" {
		lgr.Logger().Info("init coreapi")

		lnCoreApi, errLnRunApi := net.Listen("tcp", cfg.API.CoreApi.Address)
		if errLnRunApi != nil {
			return fmt.Sprintf("error listen tcp address '%s', %v", cfg.API.CoreApi.Address, errLnRunApi), 1
		}
		defer lnCoreApi.Close()

		coreAPI := coreapi.New(dsMgr, coreModules, cfg.API.CoreApi.AuthToken, lgr.Logger())

		wg.Add(1)
		go coreAPI.Run(ctx, ctxCancel, wg, lnCoreApi)
	}

	// ---------------------
	// |
	// | Runner
	// |
	lgr.Logger().Info("init runner")

	var updateInterval int
	if cfg.Scripts != nil {
		updateInterval = cfg.Scripts.UpdateInterval
	}

	rnr, errCreateRunner := runner.New(
		time.Millisecond*time.Duration(updateInterval),
		scriptsMgr,
		dsMgr,
		uploadStoragesMgr,
		coreModules,
		flg.Script,
		cfg.System,
		flg.SafeMode,
		lgr.Logger(),
	)
	if errCreateRunner != nil {
		return fmt.Sprintf("error create runner, %v", errCreateRunner), 1
	}

	lgr.Logger().Info("run runner")
	go rnr.Watch(ctx, ctxCancel, flg.Once)

	// ---------------------
	// |
	// | API
	// |
	if cfg.API != nil {
		var ln net.Listener
		ln, err = net.Listen("tcp", cfg.API.Address)
		if err != nil {
			return fmt.Sprintf("error create api listener, %v", err), 1
		}
		apis := apiManager.New(cfg.API.Address, coreStorageAlert, coreStorageKV, channelsMgr, rnr, lgr.Logger())
		wg.Add(1)
		go apis.Run(ctx, ctxCancel, wg, ln)

		if cfg.API.ServiceAddress != "" {
			var ln net.Listener
			ln, err = net.Listen("tcp", cfg.API.ServiceAddress)
			if err != nil {
				return fmt.Sprintf("error create service listener, %v", err), 1
			}
			srv := service.New(lgr.Logger())
			wg.Add(1)
			go srv.Run(ctx, ctxCancel, wg, ln)
		}
	}

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

	lgr.Logger().Info("terminate")

	return "", 0
}

func initCoreModules(
	coreStorageAlert corestorage.CoreStorage,
	coreStorageKV corestorage.CoreStorage,
	chManager *channelsManager.ChannelsManager,
	lgr *zap.Logger,
	flg *config.Flags,
) []modules.Module {
	coreModules := make([]modules.Module, 0)

	alertMod := alertModule.New(coreStorageAlert.Alert(), chManager, lgr)
	coreModules = append(coreModules, alertMod)

	kvModule := kv.New(coreStorageKV.KV())
	coreModules = append(coreModules, kvModule)

	logMod := logModule.New(lgr)
	coreModules = append(coreModules, logMod)

	chartMod := chartModule.New(lgr)
	coreModules = append(coreModules, chartMod)

	if !flg.SafeMode {
		httpMod := httpModule.New(lgr)
		coreModules = append(coreModules, httpMod)
	}

	runtimeMod := runtimeModule.New(flg, lgr)
	coreModules = append(coreModules, runtimeMod)

	tlsMod := tlsModule.New()
	coreModules = append(coreModules, tlsMod)

	metaMod := meta.New(lgr)
	coreModules = append(coreModules, metaMod)

	fileMod := file.New(lgr)
	coreModules = append(coreModules, fileMod)

	return coreModules
}

func validateLogLevel(level string) error {
	for _, l := range []string{"ERROR", "INFO", "DEBUG"} {
		if l == level {
			return nil
		}
	}

	return fmt.Errorf("wrong log level")
}
