package global

import (
	"github.com/balerter/balerter/internal/config/global/api"
	"github.com/balerter/balerter/internal/config/global/service"
	"github.com/balerter/balerter/internal/config/global/storages"
)

type Global struct {
	SendStartNotification []string          `json:"sendStartNotification" yaml:"sendStartNotification"`
	SendStopNotification  []string          `json:"sendStopNotification" yaml:"sendStopNotification"`
	API                   api.API           `json:"api" yaml:"api"`
	Storages              storages.Storages `json:"storages" yaml:"storages"`
	LuaModulesPath        string            `json:"luaModulesPath" yaml:"luaModulesPath"`
	Service               service.Service   `json:"service" yaml:"service"`
}

func (cfg *Global) Validate() error {
	if err := cfg.API.Validate(); err != nil {
		return err
	}
	if err := cfg.Storages.Validate(); err != nil {
		return err
	}
	if err := cfg.Service.Validate(); err != nil {
		return err
	}

	return nil
}
