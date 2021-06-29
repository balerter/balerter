package kv

import (
	"github.com/balerter/balerter/internal/corestorage"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// KV represents KV API module
type KV struct {
	storage corestorage.KV
	logger  *zap.Logger
}

// New creates new KV API module
func New(storage corestorage.KV, logger *zap.Logger) *KV {
	kv := &KV{
		storage: storage,
		logger:  logger,
	}

	return kv
}

// Handler creates API handlers for KV API module
func (kv *KV) Handler(r chi.Router) {
	r.Get("/", kv.handlerIndex)
}
