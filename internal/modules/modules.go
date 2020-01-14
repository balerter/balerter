package modules

import (
	lua "github.com/yuin/gopher-lua"
)

type Module interface {
	Name() string
	GetLoader() lua.LGFunction
	Stop() error
}
