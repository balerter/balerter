package test

import (
	"fmt"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type pair struct {
	main *script.Script
	test *script.Script
}

func splitScripts(scripts []*script.Script) (map[string]pair, error) {
	_scripts := make(map[string]*script.Script)
	_tests := make(map[string]*script.Script)
	for _, s := range scripts {
		if s.IsTest {
			_tests[s.Name] = s
		} else {
			_scripts[s.Name] = s
		}
	}

	pairs := make(map[string]pair)

	for name, t := range _tests {
		s, ok := _scripts[t.TestTarget]
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

// Result represents script test result
type Result struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Ok   bool   `json:"ok"`
}

// Run the test runner
func (rnr *Runner) Run() ([]modules.TestResult, bool, error) {
	var result []modules.TestResult
	ok := true

	ss, err := rnr.scriptsManager.GetWithTests()
	if err != nil {
		return nil, false, fmt.Errorf("error get scripts, %w", err)
	}

	pairs, err := splitScripts(ss)
	if err != nil {
		return nil, false, fmt.Errorf("error select tests, %w", err)
	}

	for name, pair := range pairs {
		result, err = rnr.runPair(result, name, pair)
		if err != nil {
			return nil, false, err
		}
	}

	for _, r := range result {
		if !r.Ok {
			ok = false
			break
		}
	}

	return result, ok, nil
}

func (rnr *Runner) runPair(result []modules.TestResult, name string, pair pair) ([]modules.TestResult, error) {
	rnr.logger.Debug("run test", zap.String("name", name))

	// run test file
	LTest := rnr.createLuaState(pair.test)
	err := LTest.DoString(string(pair.test.Body))
	if err != nil {
		LTest.Close()
		return nil, fmt.Errorf("error run test job, %w", err)
	}
	LTest.Close()

	// run main file
	LMain := rnr.createLuaState(pair.main)
	err = LMain.DoString(string(pair.main.Body))
	if err != nil {
		LMain.Close()
		return nil, fmt.Errorf("error run main job, %w", err)
	}
	LMain.Close()

	// collect datasources results
	results, err := rnr.dsManager.Result()
	if err != nil {
		return nil, fmt.Errorf("error get results from datasource manager, %w", err)
	}
	for _, r := range results {
		r.ScriptName = pair.test.Name
		result = append(result, r)
	}
	rnr.dsManager.Clean()

	// collect storages results
	results, err = rnr.storagesManager.Result()
	if err != nil {
		return nil, fmt.Errorf("error get results from storage manager, %w", err)
	}
	for _, r := range results {
		r.ScriptName = pair.test.Name
		result = append(result, r)
	}
	rnr.storagesManager.Clean()

	// collect errors from coreModules
	for _, mod := range rnr.coreModules {
		results, err = mod.Result()
		if err != nil {
			return nil, fmt.Errorf("error get results from '%s' module, %w", mod.Name(), err)
		}
		for _, r := range results {
			r.ScriptName = pair.test.Name
			result = append(result, r)
		}
		mod.Clean()
	}

	// total script result
	scriptResult := modules.TestResult{
		ScriptName: pair.test.Name,
		ModuleName: "RESULT",
		Message:    "PASS",
		Ok:         true,
	}

	for _, r := range result {
		if !r.Ok {
			scriptResult.Ok = false
			scriptResult.Message = "FAIL"
			break
		}
	}

	result = append(result, scriptResult)

	return result, nil
}

func (rnr *Runner) createLuaState(s *script.Script) *lua.LState {
	rnr.logger.Debug("create job", zap.String("name", s.Name))

	L := lua.NewState()

	for _, m := range rnr.coreModules {
		L.PreloadModule(m.Name(), m.GetLoader(s))
	}

	// Init storages
	for _, module := range rnr.storagesManager.Get() {
		moduleName := "storage." + module.Name()
		rnr.logger.Debug("add storage module", zap.String("name", moduleName))

		loader := module.GetLoader(s)
		L.PreloadModule(moduleName, loader)
	}

	// Init datasources
	for _, module := range rnr.dsManager.Get() {
		moduleName := "datasource." + module.Name()
		rnr.logger.Debug("add datasource module", zap.String("name", moduleName))

		loader := module.GetLoader(s)
		L.PreloadModule(moduleName, loader)
	}

	return L
}
