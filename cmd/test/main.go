package main

import (
	"flag"
	"fmt"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"github.com/balerter/balerter/internal/config"
	coreStorageManager "github.com/balerter/balerter/internal/core_storage/manager"
	dsManagerTest "github.com/balerter/balerter/internal/datasource/manager/test"
	"github.com/balerter/balerter/internal/logger"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/mock"
	"github.com/balerter/balerter/internal/modules"
	chartModule "github.com/balerter/balerter/internal/modules/chart"
	httpModule "github.com/balerter/balerter/internal/modules/http"
	"github.com/balerter/balerter/internal/modules/kv"
	logModule "github.com/balerter/balerter/internal/modules/log"
	testModule "github.com/balerter/balerter/internal/modules/test"
	runnerTest "github.com/balerter/balerter/internal/runner/test"
	scriptsManager "github.com/balerter/balerter/internal/script/manager"
	uploadStorageManagerTest "github.com/balerter/balerter/internal/upload_storage/manager/test"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
)

var (
	version = "undefined"
)

var (
	configSource = flag.String("config", "config.yml", "Configuration source. Currently supports only path to yaml file.")
	logLevel     = flag.String("logLevel", "ERROR", "Log level. ERROR, WARN, INFO or DEBUG")
	debug        = flag.Bool("debug", false, "debug mode")
	asJson       = flag.Bool("json", false, "output json format")

	defaultLuaModulesPath = "./?.lua;./modules/?.lua;./modules/?/init.lua"
)

func main() {
	var resultsOutput io.Writer

	resultsOutput = os.Stdout

	lua.LuaPathDefault = defaultLuaModulesPath

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
	cfg, err := config.New(*configSource)
	if err != nil {
		lgr.Logger().Error("error init config", zap.Error(err))
		os.Exit(1)
	}
	lgr.Logger().Debug("loaded configuration", zap.Any("config", cfg))

	if cfg.Global.LuaModulesPath != "" {
		lua.LuaPathDefault = cfg.Global.LuaModulesPath
	}

	lgr.Logger().Debug("lua modules path", zap.String("path", lua.LuaPathDefault))

	// Scripts sources
	lgr.Logger().Info("init scripts manager")
	scriptsMgr := scriptsManager.New()
	if err := scriptsMgr.Init(cfg.Scripts.Sources); err != nil {
		lgr.Logger().Error("error init scripts manager", zap.Error(err))
		os.Exit(1)
	}

	// datasources
	lgr.Logger().Info("init datasources manager")
	dsMgr := dsManagerTest.New(lgr.Logger())
	if err := dsMgr.Init(cfg.DataSources); err != nil {
		lgr.Logger().Error("error init datasources manager", zap.Error(err))
		os.Exit(1)
	}

	// upload storages
	lgr.Logger().Info("init upload storages manager")
	uploadStoragesMgr := uploadStorageManagerTest.New(lgr.Logger())
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

	alertMgr := mock.New(alertManager.ModuleName(), alertManager.Methods(), lgr.Logger())
	coreModules = append(coreModules, alertMgr)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | KV
	// |
	kvEngine, err := coreStoragesMgr.Get("memory")
	if err != nil {
		lgr.Logger().Error("error get kv storage engine", zap.String("name", cfg.Global.Storages.KV), zap.Error(err))
		os.Exit(1)
	}
	lgr.Logger().Info("init kv storage", zap.String("engine", cfg.Global.Storages.KV))
	kvModule := kv.New(kvEngine)
	coreModules = append(coreModules, kvModule)

	// ---------------------
	// |
	// | API
	// |
	//wg.Add(1)
	//apis := apiManager.New(cfg.Global.API, alertManagerStorageEngine, kvEngine, lgr.Logger())
	//go apis.Run(ctx, ctxCancel, wg)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | Log
	// |
	logMod := mock.New("log", logModule.Methods(), lgr.Logger())
	coreModules = append(coreModules, logMod)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | Chart
	// |
	chartMod := chartModule.New(lgr.Logger()) // todo: mock it?
	coreModules = append(coreModules, chartMod)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | http
	// |
	httpMod := httpModule.New(lgr.Logger()) // todo: mock it
	coreModules = append(coreModules, httpMod)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | test
	// |
	testMod := testModule.New(dsMgr, uploadStoragesMgr, alertMgr, logMod, lgr.Logger())
	coreModules = append(coreModules, testMod)

	// ---------------------
	// |
	// | Runner
	// |
	lgr.Logger().Info("init runner")
	rnr := runnerTest.New(scriptsMgr, dsMgr, uploadStoragesMgr, alertMgr, logMod, coreModules, lgr.Logger())

	lgr.Logger().Info("run runner: tests")
	results, ok, err := rnr.Run()
	if err != nil {
		lgr.Logger().Error("error run tests", zap.Error(err))
		os.Exit(1)
	}

	if err := output(results, resultsOutput, *asJson); err != nil {
		lgr.Logger().Error("error output results", zap.Error(err))
		os.Exit(1)
	}

	lgr.Logger().Info("terminate")

	if !ok {
		os.Exit(2)
		return
	}

	os.Exit(0)
}

func validateLogLevel(level string) error {
	for _, l := range []string{"ERROR", "WARN", "INFO", "DEBUG"} {
		if l == level {
			return nil
		}
	}

	return fmt.Errorf("wrong log level")
}
