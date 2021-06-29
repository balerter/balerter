package kv

import (
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
)

// ModuleName returns the module name
func ModuleName() string {
	return "kv"
}

// Methods returns module methods
func Methods() []string {
	return []string{
		"all",
		"get",
		"put",
		"delete",
		"upsert",
	}
}

// KV represents the KV core module
type KV struct {
	engine coreStorage.KV
}

// New creates new KV core module
func New(engine coreStorage.KV) *KV {
	kv := &KV{
		engine: engine,
	}

	return kv
}

// Name returns the module name
func (kv *KV) Name() string {
	return ModuleName()
}

// GetLoader returns the lua loader
func (kv *KV) GetLoader(_ *script.Script) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"all":    kv.all,
				"get":    kv.get,
				"put":    kv.put,
				"delete": kv.delete,
				"upsert": kv.upsert,
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			luaState.Push(mod)
			return 1
		}
	}()
}

// Stop the module
func (kv *KV) Stop() error {
	return nil
}

func (kv *KV) all(luaState *lua.LState) int {
	data, err := kv.engine.All()
	if err != nil {
		luaState.Push(lua.LString(""))
		luaState.Push(lua.LString(err.Error()))
		return 2
	}

	res := &lua.LTable{}

	for key, value := range data {
		res.RawSetString(key, lua.LString(value))
	}

	luaState.Push(res)
	luaState.Push(lua.LNil)

	return 2
}

func (kv *KV) get(luaState *lua.LState) int {
	varName := luaState.Get(1).String()

	val, err := kv.engine.Get(varName)
	if err != nil {
		luaState.Push(lua.LString(""))
		luaState.Push(lua.LString(err.Error()))
		return 2
	}

	luaState.Push(lua.LString(val))
	luaState.Push(lua.LNil)

	return 2
}

func (kv *KV) put(luaState *lua.LState) int {
	varName := luaState.Get(1).String()
	varVal := luaState.Get(2).String()

	err := kv.engine.Put(varName, varVal)
	if err != nil {
		luaState.Push(lua.LString(err.Error()))
		return 1
	}

	return 0
}

func (kv *KV) upsert(luaState *lua.LState) int {
	varName := luaState.Get(1).String()
	varVal := luaState.Get(2).String()

	err := kv.engine.Upsert(varName, varVal)
	if err != nil {
		luaState.Push(lua.LString(err.Error()))
		return 1
	}

	return 0
}

func (kv *KV) delete(luaState *lua.LState) int {
	varName := luaState.Get(1).String()

	err := kv.engine.Delete(varName)
	if err != nil {
		luaState.Push(lua.LString(err.Error()))
		return 1
	}

	return 0
}
