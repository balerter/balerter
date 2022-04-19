package prometheus_models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatrix(t *testing.T) {
	s := `{
 "status": "success",
 "data": {
   "resultType": "matrix",
   "result": [
     {
       "metric": {"foo":"1"},
       "values": [
         [
           1647162924,
           "0.633393218801267"
         ],
         [
           1647249324,
           "0.6034093331590693"
         ]
       ]
     },
     {
       "metric": {},
       "values": [
         [
           1747162924,
           "0.733393218801267"
         ],
         [
           1747249324,
           "0.7034093331590693"
         ]
       ]
     }
   ]
 }
}`

	resp := &APIResponse{}

	err := json.Unmarshal([]byte(s), resp)
	require.NoError(t, err)

	assert.Equal(t, ValMatrix, resp.Data.Type)
	require.IsType(t, Matrix{}, resp.Data.Value)

	m, ok := resp.Data.Value.(Matrix)
	require.True(t, ok)

	require.Equal(t, 2, len(m))

	expect := Matrix{
		{
			Metric: map[string]string{"foo": "1"},
			Values: []SamplePair{
				{Timestamp: 1647162924, Value: 0.633393218801267},
				{Timestamp: 1647249324, Value: 0.6034093331590693},
			},
		},
		{
			Metric: map[string]string{},
			Values: []SamplePair{
				{Timestamp: 1747162924, Value: 0.733393218801267},
				{Timestamp: 1747249324, Value: 0.7034093331590693},
			},
		},
	}

	require.Equal(t, expect, m)
}

func TestVector(t *testing.T) {
	s := `
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "datacenter": "eu-1",
          "version": "v2"
        },
        "value": [
          1647251030,
          "1"
        ]
      },
      {
        "metric": {
          "job": "bar"
        },
        "value": [
          1647251030,
          "1"
        ]
      },
      {
        "metric": {
          "job": "node_exporter",
          "node": "foo"
        },
        "value": [
          1647251030,
          "1"
        ]
      }
	]
  }
}
`

	resp := &APIResponse{}

	err := json.Unmarshal([]byte(s), resp)
	require.NoError(t, err)

	assert.Equal(t, ValVector, resp.Data.Type)
	require.IsType(t, Vector{}, resp.Data.Value)

	m, ok := resp.Data.Value.(Vector)
	require.True(t, ok)

	require.Equal(t, 3, len(m))

	expect := Vector{
		{
			Metric: map[string]string{"datacenter": "eu-1", "version": "v2"},
			Value:  SamplePair{Timestamp: 1647251030, Value: 1},
		},
		{
			Metric: map[string]string{"job": "bar"},
			Value:  SamplePair{Timestamp: 1647251030, Value: 1},
		},
		{
			Metric: map[string]string{"job": "node_exporter", "node": "foo"},
			Value:  SamplePair{Timestamp: 1647251030, Value: 1},
		},
	}

	require.Equal(t, expect, m)
}

func TestLoki(t *testing.T) {
	s := `
{
    "status": "success",
    "data": {
        "resultType": "streams",
        "result": [
            {
                "stream": {
                    "datacenter": "eu-1"
                },
                "values": [
                    [
                        "1647253163106308798",
                        "foobar"
                    ]
                ]
            },
            {
                "stream": {
                    "node": "node24",
                    "level": "error"
                },
                "values": [
                    [
                        "1647253143809759832",
                        "foobar"
                    ],
                    [
                        "1647253143809469604",
                        "barbar"
                    ],
                    [
                        "1647253143809464871",
                        "barbarbar"
                    ]
                ]
            }
        ]
    }
}
`

	resp := &APIResponse{}

	err := json.Unmarshal([]byte(s), resp)
	require.NoError(t, err)

	assert.Equal(t, ValStreams, resp.Data.Type)
	require.IsType(t, Streams{}, resp.Data.Value)

	m, ok := resp.Data.Value.(Streams)
	require.True(t, ok)

	require.Equal(t, 2, len(m))

	expect := Streams{
		{
			Metric: map[string]string{"datacenter": "eu-1"},
			Values: []StreamSamplePair{
				{1647253163106308798, "foobar"},
			},
		},
		{
			Metric: map[string]string{"node": "node24", "level": "error"},
			Values: []StreamSamplePair{
				{1647253143809759832, "foobar"},
				{1647253143809469604, "barbar"},
				{1647253143809464871, "barbarbar"},
			},
		},
	}

	require.Equal(t, expect, m)
}
