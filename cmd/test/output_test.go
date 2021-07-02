package main

import (
	"bytes"
	"github.com/balerter/balerter/internal/modules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_outputJSON(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	err := outputJSON([]modules.TestResult{
		{
			ScriptName: "a",
			ModuleName: "b",
			Message:    "c",
			Ok:         true,
		},
	}, buf)
	require.NoError(t, err)

	assert.Equal(t, `[{"script":"a","module":"b","message":"c","ok":true}]`, buf.String())
}

func Test_outputPlainColored(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	err := outputPlainColored([]modules.TestResult{
		{
			ScriptName: "a",
			ModuleName: "b",
			Message:    "c",
			Ok:         true,
		},
	}, buf)
	require.NoError(t, err)

	assert.Equal(t, "[PASS]\t[a]\t[b]\tc\n", buf.String())
}

func Test_output(t *testing.T) {
	mods := []modules.TestResult{
		{
			ScriptName: "a",
			ModuleName: "b",
			Message:    "c",
			Ok:         true,
		},
	}

	coloredBuf := bytes.NewBuffer(nil)
	jsonBuf := bytes.NewBuffer(nil)

	err := outputPlainColored(mods, coloredBuf)
	require.NoError(t, err)
	err = outputJSON(mods, jsonBuf)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)

	err = output(mods, buf, true)
	require.NoError(t, err)
	assert.Equal(t, buf.String(), jsonBuf.String())
	buf.Reset()

	err = output(mods, buf, false)
	require.NoError(t, err)
	assert.Equal(t, buf.String(), coloredBuf.String())
}
