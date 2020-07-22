package test

import (
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_splitScripts(t *testing.T) {
	src := []*script.Script{
		{Name: "s1", IsTest: false},
		{Name: "s1_test", IsTest: true, TestTarget: "s1"},
		{Name: "s2", IsTest: false},
		{Name: "s2_test", IsTest: true, TestTarget: "s2"},
	}

	pairs, err := splitScripts(src)
	require.NoError(t, err)

	assert.Equal(t, 2, len(pairs))

	p, ok := pairs["s1_test"]
	require.True(t, ok)
	assert.Equal(t, "s1_test", p.test.Name)
	assert.Equal(t, "s1", p.main.Name)

	p, ok = pairs["s2_test"]
	require.True(t, ok)
	assert.Equal(t, "s2_test", p.test.Name)
	assert.Equal(t, "s2", p.main.Name)
}

func Test_splitScripts_no_main(t *testing.T) {
	src := []*script.Script{
		{Name: "s1_test", IsTest: true, TestTarget: "s1"},
	}

	_, err := splitScripts(src)
	require.Error(t, err)

	assert.Equal(t, "main script for test 's1_test' not found", err.Error())
}
