package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"github.com/balerter/balerter/internal/config"
	dsManagerTest "github.com/balerter/balerter/internal/datasource/manager/test"
	"github.com/balerter/balerter/internal/logger"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/mock"
	"github.com/balerter/balerter/internal/modules"
	chartModule "github.com/balerter/balerter/internal/modules/chart"
	httpModule "github.com/balerter/balerter/internal/modules/http"
	"github.com/balerter/balerter/internal/modules/kv"
	logModule "github.com/balerter/balerter/internal/modules/log"
	runtimeModule "github.com/balerter/balerter/internal/modules/runtime"
	testModule "github.com/balerter/balerter/internal/modules/test"
	runnerTest "github.com/balerter/balerter/internal/runner/test"
	scriptsManager "github.com/balerter/balerter/internal/script/manager"
	uploadStorageManagerTest "github.com/balerter/balerter/internal/upload_storage/manager/test"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
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
	var resultsOutput io.Writer = os.Stdout

	lua.LuaPathDefault = defaultLuaModulesPath

	coreModules := make([]modules.ModuleTest, 0)

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
	if err := dsMgr.Init(&cfg.DataSources); err != nil {
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

	// ---------------------
	// |
	// | Core Modules
	// |
	// | AlertManager
	// |
	alertMgr := mock.New(alertManager.ModuleName(), alertManager.Methods(), lgr.Logger())
	coreModules = append(coreModules, alertMgr)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | KV
	// |
	kvModule := mock.New(kv.ModuleName(), kv.Methods(), lgr.Logger())
	coreModules = append(coreModules, kvModule)

	// ---------------------
	// |
	// | API
	// |
	// module is not used in the test environment

	// ---------------------
	// |
	// | Core Modules
	// |
	// | Log
	// |
	logMod := mock.New(logModule.ModuleName(), logModule.Methods(), lgr.Logger())
	coreModules = append(coreModules, logMod)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | Chart
	// |
	chartMod := mock.New(chartModule.ModuleName(), chartModule.Methods(), lgr.Logger())
	coreModules = append(coreModules, chartMod)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | http
	// |
	httpMod := mock.New(httpModule.ModuleName(), httpModule.Methods(), lgr.Logger())
	coreModules = append(coreModules, httpMod)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | runtime
	// |
	runtimeMod := mock.New(runtimeModule.ModuleName(), runtimeModule.Methods(), lgr.Logger())
	coreModules = append(coreModules, runtimeMod)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | test
	// |
	testMod := testModule.New(dsMgr, uploadStoragesMgr, coreModules, lgr.Logger())

	// ---------------------
	// |
	// | Runner
	// |
	lgr.Logger().Info("init runner")
	rnr := runnerTest.New(scriptsMgr, dsMgr, uploadStoragesMgr, testMod, coreModules, lgr.Logger())

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
