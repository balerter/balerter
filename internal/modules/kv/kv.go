package kv

import (
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
)

func ModuleName() string {
	return "kv"
}

func Methods() []string {
	return []string{
		"get",
		"put",
		"delete",
		"upsert",
	}
}

type KV struct {
	engine coreStorage.CoreStorage
}

func New(engine coreStorage.CoreStorage) *KV {
	kv := &KV{
		engine: engine,
	}

	return kv
}

func (kv *KV) Name() string {
	return ModuleName()
}

func (kv *KV) GetLoader(_ *script.Script) lua.LGFunction {
	return func() lua.LGFunction {
		return func(luaState *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"get":    kv.get,
				"put":    kv.put,
				"delete": kv.delete,
				"upsert": kv.upsert,
			}

			mod := luaState.SetFuncs(luaState.NewTable(), exports)

			luaState.Push(mod)
			return 1 //nolint:mnd // its arguments count, not magic number
		}
	}()
}

func (kv *KV) Stop() error {
	return nil
}

func (kv *KV) get(luaState *lua.LState) int {
	varName := luaState.Get(1).String() //nolint:mnd // its argument index, not magic number

	val, err := kv.engine.KV().Get(varName)
	if err != nil {
		luaState.Push(lua.LString(""))
		luaState.Push(lua.LString(err.Error()))
		return 2 //nolint:mnd // its arguments count, not magic number
	}

	luaState.Push(lua.LString(val))
	luaState.Push(lua.LNil)

	return 2 //nolint:mnd // its arguments count, not magic number
}

func (kv *KV) put(luaState *lua.LState) int {
	varName := luaState.Get(1).String() //nolint:mnd // its argument index, not magic number
	varVal := luaState.Get(2).String()  //nolint:mnd // its argument index, not magic number

	err := kv.engine.KV().Put(varName, varVal)
	if err != nil {
		luaState.Push(lua.LString(err.Error()))
		return 1 //nolint:mnd // its arguments count, not magic number
	}

	return 0 //nolint:mnd // its arguments count, not magic number
}

func (kv *KV) upsert(luaState *lua.LState) int {
	varName := luaState.Get(1).String() //nolint:mnd // its argument index, not magic number
	varVal := luaState.Get(2).String()  //nolint:mnd // its argument index, not magic number

	err := kv.engine.KV().Upsert(varName, varVal)
	if err != nil {
		luaState.Push(lua.LString(err.Error()))
		return 1 //nolint:mnd // its arguments count, not magic number
	}

	return 0 //nolint:mnd // its arguments count, not magic number
}

func (kv *KV) delete(luaState *lua.LState) int {
	varName := luaState.Get(1).String() //nolint:mnd // its argument index, not magic number

	err := kv.engine.KV().Delete(varName)
	if err != nil {
		luaState.Push(lua.LString(err.Error()))
		return 1 //nolint:mnd // its arguments count, not magic number
	}

	return 0 //nolint:mnd // its arguments count, not magic number
}
