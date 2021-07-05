package loki

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryOptions_Validate_error(t *testing.T) {
	q := &queryOptions{Direction: "bad"}
	err := q.Validate()
	require.Error(t, err)
	assert.Equal(t, "option Direction support only values: 'forward' and 'backward'", err.Error())
}

func TestQueryOptions_Validate(t *testing.T) {
	q := &queryOptions{Direction: ""}
	err := q.Validate()
	require.NoError(t, err)
}

func TestRangeOptions_Validate_error(t *testing.T) {
	q := &rangeOptions{Direction: "bad"}
	err := q.Validate()
	require.Error(t, err)
	assert.Equal(t, "option Direction support only values: 'forward' and 'backward'", err.Error())
}

func TestRangeOptions_Validate(t *testing.T) {
	q := &rangeOptions{Direction: ""}
	err := q.Validate()
	require.NoError(t, err)
}
