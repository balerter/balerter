package email

import (
	"testing"

	"github.com/balerter/balerter/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestEmailName(t *testing.T) {
	s := &Email{name: "foo"}
	assert.Equal(t, "foo", s.Name())
}

func TestEmailConfigNoAuthMechs(t *testing.T) {
	email := &Email{
		conf: &config.ChannelEmail{Name: "emailTest", AuthUsername: "test"}, logger: zap.NewExample(),
	}
	_, err := email.auth("")
	require.Error(t, err)
	require.Equal(t, err.Error(), "unknown auth mechanism: ")
}

func TestEmailConfigMissingAuthParam(t *testing.T) {
	conf := &config.ChannelEmail{AuthUsername: "test"}
	email := &Email{
		conf: conf, logger: zap.NewExample(),
	}
	_, err := email.auth("CRAM-MD5")
	require.Error(t, err)
	require.Equal(t, err.Error(), "missing secret for CRAM-MD5 auth mechanism")

	_, err = email.auth("PLAIN")
	require.Error(t, err)
	require.Equal(t, err.Error(), "missing password for PLAIN auth mechanism")

	_, err = email.auth("LOGIN")
	require.Error(t, err)
	require.Equal(t, err.Error(), "missing password for LOGIN auth mechanism")

	_, err = email.auth("PLAIN LOGIN")
	require.Error(t, err)
	require.Equal(t, err.Error(), "missing password for PLAIN auth mechanism; missing password for LOGIN auth mechanism")
}

func TestEmailNoUsernameStillOk(t *testing.T) {
	email := &Email{
		conf: &config.ChannelEmail{}, logger: zap.NewExample(),
	}
	a, err := email.auth("CRAM-MD5")
	require.NoError(t, err)
	require.Nil(t, a)
}
