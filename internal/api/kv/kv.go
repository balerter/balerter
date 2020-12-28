package kv

import (
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type KV struct {
	storage corestorage.KV
	logger  *zap.Logger
}

func New(storage corestorage.KV, logger *zap.Logger) *KV {
	kv := &KV{
		storage: storage,
		logger:  logger,
	}

	return kv
}

func (kv *KV) Handler(r chi.Router) {
	r.Get("/", kv.handlerIndex)
}
