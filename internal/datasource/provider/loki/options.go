package loki

import (
	"fmt"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type queryOptions struct {
	Time      string
	Limit     int
	Direction string
}

// Validate the options
func (q *queryOptions) Validate() error {
	if err := directionValidate(q.Direction); err != nil {
		return err
	}
	return nil
}

type rangeOptions struct {
	Limit     int
	Start     string
	End       string
	Step      string
	Direction string
}

// Validate the options
func (q *rangeOptions) Validate() error {
	if err := directionValidate(q.Direction); err != nil {
		return err
	}
	return nil
}

type options interface {
	Validate() error
}

func (m *Loki) parseOptions(luaState *lua.LState, opts options) error {
	options := luaState.Get(2) // nolint:gomnd // param position
	if options.Type() == lua.LTNil {
		return nil
	}

	if options.Type() != lua.LTTable {
		return fmt.Errorf("options must be a table")
	}

	err := gluamapper.Map(options.(*lua.LTable), &opts)
	if err != nil {
		return fmt.Errorf("error parse, %w", err)
	}

	if err := opts.Validate(); err != nil {
		return err
	}

	return nil
}
