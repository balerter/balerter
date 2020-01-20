package api

import (
	"github.com/balerter/balerter/internal/config"
)

type API struct {
	address string
}

func New(cfg config.API) *API {
	api := &API{
		address: cfg.Address,
	}

	return api
}
