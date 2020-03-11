package alerts

import (
	"github.com/balerter/balerter/internal/alert/alert"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func Test_filter_no(t *testing.T) {
	data := []*alertManager.AlertInfo{
		{Level: alert.LevelError, Name: "foo4"},
		{Level: alert.LevelError, Name: "foo1"},
		{Level: alert.LevelError, Name: "foo2"},
		{Level: alert.LevelWarn, Name: "foo3"},
		{Level: alert.LevelSuccess, Name: "foo4"},
		{Level: alert.LevelSuccess, Name: "foo5"},
	}

	req := &http.Request{URL: &url.URL{RawQuery: "level=&name="}}

	output, err := filter(req, data)
	require.NoError(t, err)
	require.Equal(t, 6, len(output))
}

func Test_filter(t *testing.T) {
	data := []*alertManager.AlertInfo{
		{Level: alert.LevelError, Name: "foo4"},
		{Level: alert.LevelError, Name: "foo1"},
		{Level: alert.LevelError, Name: "foo2"},
		{Level: alert.LevelWarn, Name: "foo3"},
		{Level: alert.LevelSuccess, Name: "foo4"},
		{Level: alert.LevelSuccess, Name: "foo5"},
	}

	req := &http.Request{URL: &url.URL{RawQuery: "level=warn,success&name=foo3,foo4"}}

	output, err := filter(req, data)
	require.NoError(t, err)
	require.Equal(t, 2, len(output))
	require.Contains(t, output, &alertManager.AlertInfo{Level: alert.LevelWarn, Name: "foo3"})
	require.Contains(t, output, &alertManager.AlertInfo{Level: alert.LevelSuccess, Name: "foo4"})
}

func Test_filter_bad_level(t *testing.T) {
	data := []*alertManager.AlertInfo{
		{Level: alert.LevelError, Name: "foo1"},
	}

	req := &http.Request{URL: &url.URL{RawQuery: "level=warn2"}}

	_, err := filter(req, data)
	require.Error(t, err)
	assert.Equal(t, "bad level value", err.Error())
}
