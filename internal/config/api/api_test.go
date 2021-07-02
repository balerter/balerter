package api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAPI_Validate(t *testing.T) {
	a := API{}
	require.NoError(t, a.Validate())
}
