package file

import (
	"os"

	"github.com/balerter/balerter/internal/modules"

	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

// ModuleName returns the module name
func ModuleName() string {
	return "file"
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"save",
		"load",
	}
}

type File struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *File {
	f := &File{
		logger: logger,
	}

	return f
}

// Name returns the module name
func (f *File) Name() string {
	return ModuleName()
}

// GetLoader returns the lua loader
func (f *File) GetLoader(_ modules.Job) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"save": f.save,
				"load": f.load,
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			luaState.Push(mod)
			return 1
		}
	}()
}

// Stop the module
func (f *File) Stop() error {
	return nil
}

func (f *File) save(luaState *lua.LState) int {
	filename := luaState.Get(1).String()
	data := luaState.Get(2).String() // nolint:gomnd // param position

	errWrite := os.WriteFile(filename, []byte(data), 0644)
	if errWrite != nil {
		luaState.Push(lua.LString(errWrite.Error()))
		return 1
	}

	return 0
}

func (f *File) load(luaState *lua.LState) int {
	filename := luaState.Get(1).String()

	data, errRead := os.ReadFile(filename)
	if errRead != nil {
		luaState.Push(lua.LString(""))
		luaState.Push(lua.LString(errRead.Error()))
		return 2
	}

	luaState.Push(lua.LString(data))

	return 1
}
