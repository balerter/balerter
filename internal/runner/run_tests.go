package runner

import (
	"fmt"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strings"
)

type pair struct {
	main *script.Script
	test *script.Script
}

func splitScripts(scripts []*script.Script) (map[string]pair, error) {
	_scripts := make(map[string]*script.Script)
	_tests := make(map[string]*script.Script)
	for _, s := range scripts {
		if strings.HasSuffix(s.Name, "_test") {
			_tests[s.Name] = s
		} else {
			_scripts[s.Name] = s
		}
	}

	pairs := make(map[string]pair)

	for name, t := range _tests {
		s, ok := _scripts[strings.TrimSuffix(name, "_test")]
		if !ok {
			return nil, fmt.Errorf("main script for test '%s' not found", name)
		}

		pairs[name] = pair{
			main: s,
			test: t,
		}
	}

	return pairs, nil
}

func (rnr *Runner) RunTests() []error {
	var errors []error

	ss, err := rnr.scriptsManager.Get()
	if err != nil {
		errors = append(errors, fmt.Errorf("error get scripts, %w", err))
		return errors
	}

	pairs, err := splitScripts(ss)
	if err != nil {
		errors = append(errors, fmt.Errorf("error select tests, %w", err))
		return errors
	}

	for name, pair := range pairs {
		rnr.logger.Debug("run test", zap.String("name", name))

		jTest := newJob(pair.test, rnr.logger)
		LTest := rnr.createLuaTestingState(jTest)
		err := LTest.DoString(string(jTest.script.Body))
		if err != nil {
			errors = append(errors, fmt.Errorf("error run test job, %w", err))
			LTest.Close()
			continue
		}
		LTest.Close()

		jMain := newJob(pair.main, rnr.logger)
		LMain := rnr.createLuaTestingState(jMain)
		err = LMain.DoString(string(jMain.script.Body))
		if err != nil {
			errors = append(errors, fmt.Errorf("error run main job, %w", err))
			LMain.Close()
			continue
		}
		LMain.Close()

		// todo check test expectations
	}

	errors = append(errors, rnr.dsManager.Errors()...)

	return errors
}

// todo refactoring duplicated code
func (rnr *Runner) createLuaTestingState(j *Job) *lua.LState {
	rnr.logger.Debug("create job", zap.String("name", j.name))

	L := lua.NewState()

	for _, m := range rnr.coreModules {
		L.PreloadModule(m.Name(), m.GetLoader(j.script))
	}

	// Init storages
	for _, module := range rnr.storagesManager.Get() {
		moduleName := "storage." + module.Name()
		rnr.logger.Debug("add storage module", zap.String("name", moduleName))

		loader := module.GetLoader(j.script)
		L.PreloadModule(moduleName, loader)
	}

	// Init datasources
	for _, module := range rnr.dsManager.GetMocks() {
		moduleName := "datasource." + module.Name()
		rnr.logger.Debug("add datasource module", zap.String("name", moduleName))

		loader := module.GetLoader(j.script)
		L.PreloadModule(moduleName, loader)
	}

	return L
}
