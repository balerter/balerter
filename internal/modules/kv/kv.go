package kv

import (
	"fmt"
	"github.com/balerter/balerter/internal/modules/kv/provider"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
)

type KVer interface {
	Put(string, string) error
	Get(string) (string, error)
	Upsert(string, string) error
	Delete(string) error
}

type ProviderType int

const (
	ProviderTypeMemory ProviderType = 0
)

type KV struct {
	engine KVer
}

func New(providerType ProviderType) (*KV, error) {
	kv := &KV{}

	switch providerType {
	case ProviderTypeMemory:
		kv.engine = provider.New()
	default:
		return nil, fmt.Errorf("unexpected kv provider type")
	}

	return kv, nil
}

func (kv *KV) Name() string {
	return "kv"
}

func (kv *KV) GetLoader(script *script.Script) lua.LGFunction {
	return func() lua.LGFunction {
		return func(L *lua.LState) int {
			var exports = map[string]lua.LGFunction{
				"get":    kv.get,
				"put":    kv.put,
				"delete": kv.delete,
				"upsert": kv.upsert,
			}

			mod := L.SetFuncs(L.NewTable(), exports)

			L.Push(mod)
			return 1
		}
	}()
}

func (kv *KV) Stop() error {
	return nil
}

func (kv *KV) get(L *lua.LState) int {
	varName := L.Get(1).String()

	val, err := kv.engine.Get(varName)
	if err != nil {
		L.Push(lua.LString(""))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(val))
	L.Push(lua.LNil)

	return 2
}

func (kv *KV) put(L *lua.LState) int {
	varName := L.Get(1).String()
	varVal := L.Get(2).String()

	err := kv.engine.Put(varName, varVal)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	return 0
}

func (kv *KV) upsert(L *lua.LState) int {
	varName := L.Get(1).String()
	varVal := L.Get(2).String()

	err := kv.engine.Upsert(varName, varVal)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	return 0
}

func (kv *KV) delete(L *lua.LState) int {
	varName := L.Get(1).String()

	err := kv.engine.Delete(varName)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	return 0
}
