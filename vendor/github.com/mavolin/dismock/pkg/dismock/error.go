package dismock

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/utils/httputil"
	"github.com/stretchr/testify/require"
)

// Error simulates an error response for the given path using the given method.
func (m *Mocker) Error(method, path string, e httputil.HTTPError) {
	m.MockAPI("Error", method, path, func(w http.ResponseWriter, r *http.Request, t *testing.T) {
		w.WriteHeader(e.Status)
		err := json.NewEncoder(w).Encode(e)
		require.NoError(t, err)
	})
}
