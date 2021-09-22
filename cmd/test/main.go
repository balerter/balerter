package main

import (
	"flag"
	"fmt"
	"github.com/balerter/balerter/internal/modules/api"
	"github.com/balerter/balerter/internal/modules/meta"
	"io"
	"log"
	"os"

	"github.com/balerter/balerter/internal/config"
	dsManagerTest "github.com/balerter/balerter/internal/datasource/manager/test"
	"github.com/balerter/balerter/internal/logger"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/mock"
	"github.com/balerter/balerter/internal/modules"
	alertModule "github.com/balerter/balerter/internal/modules/alert"
	chartModule "github.com/balerter/balerter/internal/modules/chart"
	httpModule "github.com/balerter/balerter/internal/modules/http"
	"github.com/balerter/balerter/internal/modules/kv"
	logModule "github.com/balerter/balerter/internal/modules/log"
	runtimeModule "github.com/balerter/balerter/internal/modules/runtime"
	testModule "github.com/balerter/balerter/internal/modules/test"
	tlsModule "github.com/balerter/balerter/internal/modules/tls"
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

func run(cfg *config.Config, flg *config.Flags) (string, int) {
	var resultsOutput io.Writer = os.Stdout

	lua.LuaPathDefault = defaultLuaModulesPath

	coreModules := make([]modules.ModuleTest, 0)

	var err error

	if err = validateLogLevel(flg.LogLevel); err != nil {
		return err.Error(), 1
	}

	lgr, err := logger.New(flg.LogLevel, flg.Debug)
	if err != nil {
		return err.Error(), 1
	}

	metrics.SetVersion(version)

	lgr.Logger().Info("balerter start", zap.String("version", version))

	// Configuration
	lgr.Logger().Debug("loaded configuration", zap.Any("config", cfg), zap.Any("flags", flg))

	if cfg.LuaModulesPath != "" {
		lua.LuaPathDefault = cfg.LuaModulesPath
	}

	lgr.Logger().Debug("lua modules path", zap.String("path", lua.LuaPathDefault))

	// Scripts sources
	lgr.Logger().Info("init scripts manager")
	scriptsMgr := scriptsManager.New()
	if err = scriptsMgr.Init(cfg.Scripts); err != nil {
		return fmt.Sprintf("error init scripts sources, %v", err), 1
	}

	// datasources
	lgr.Logger().Info("init datasources manager")
	dsMgr := dsManagerTest.New(lgr.Logger())
	if err = dsMgr.Init(cfg.DataSources); err != nil {
		return fmt.Sprintf("error init datasources manager, %v", err), 1
	}

	// upload storages
	lgr.Logger().Info("init upload storages manager")
	uploadStoragesMgr := uploadStorageManagerTest.New(lgr.Logger())
	if err = uploadStoragesMgr.Init(cfg.StoragesUpload); err != nil {
		return fmt.Sprintf("error init upload storages manager, %v", err), 1
	}

	// ---------------------
	// |
	// | Core Modules
	// |
	// | AlertManager
	// |
	alertMgr := mock.New(alertModule.ModuleName(), alertModule.Methods(), lgr.Logger())
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
	// | API
	// |
	apiMod := mock.New(api.ModuleName(), api.Methods(), lgr.Logger())
	coreModules = append(coreModules, apiMod)

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
	// | tls
	// |
	tlsMod := mock.New(tlsModule.ModuleName(), tlsModule.Methods(), lgr.Logger())
	coreModules = append(coreModules, tlsMod)

	// ---------------------
	// |
	// | Core Modules
	// |
	// | meta
	// |
	metaMod := mock.New(meta.ModuleName(), meta.Methods(), lgr.Logger())
	coreModules = append(coreModules, metaMod)

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
	coreModules = append(coreModules, testMod)

	// ---------------------
	// |
	// | Runner
	// |
	lgr.Logger().Info("init runner")
	rnr := runnerTest.New(scriptsMgr, dsMgr, uploadStoragesMgr, coreModules, lgr.Logger())

	lgr.Logger().Info("run runner: tests")
	results, ok, err := rnr.Run()
	if err != nil {
		lgr.Logger().Error("error run tests", zap.Error(err))
		os.Exit(1)
	}

	if err := output(results, resultsOutput, flg.AsJSON); err != nil {
		return fmt.Sprintf("error out results, %v", err), 1
	}

	lgr.Logger().Info("terminate")

	if !ok {
		return "", 1
	}

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
