package kv

import (
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
)

type KV struct {
}

func New() *KV {
	kv := &KV{}

	return kv
}

func (kv *KV) Name() string {
	return "kv"
}

func (kv *KV) GetLoader(script *script.Script) lua.LGFunction {
	return nil
}

func (kv *KV) Stop() error {
	return nil
}
