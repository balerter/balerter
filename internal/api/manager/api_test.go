package manager

import (
	"github.com/stretchr/testify/assert"
	httpTestify "github.com/stretchr/testify/http"
	"go.uber.org/zap"
	"testing"
)

func TestAPI_Liveness(t *testing.T) {
	api := &API{
		logger: zap.NewNop(),
	}

	rw := &httpTestify.TestResponseWriter{}

	api.handlerLiveness(rw, nil)

	assert.Equal(t, "ok", rw.Output)
}
