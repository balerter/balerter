package email

import (
	"github.com/balerter/balerter/internal/config/channels/email"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailName(t *testing.T) {
	s := &Email{name: "foo"}
	assert.Equal(t, "foo", s.Name())
}

func TestNew(t *testing.T) {
	e, err := New(email.Email{}, zap.NewNop())
	require.NoError(t, err)
	assert.IsType(t, &Email{}, e)
}

func TestEmail_Ignore(t *testing.T) {
	e := &Email{ignore: true}
	assert.True(t, e.Ignore())
}
