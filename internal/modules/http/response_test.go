package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func Test_response_toLuaTable(t *testing.T) {
	type fields struct {
		StatusCode int
		Body       []byte
		Headers    map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   func(t2 *testing.T, tbl *lua.LTable)
	}{
		{
			name: "test",
			fields: fields{
				StatusCode: 10,
				Body:       []byte("foo"),
				Headers:    map[string]string{"a": "b"},
			},
			want: func(t2 *testing.T, tbl *lua.LTable) {
				assert.Equal(t2, "10", tbl.RawGetString("status_code").String())
				assert.Equal(t2, "foo", tbl.RawGetString("body").String())
				h := tbl.RawGetString("headers")
				require.Equal(t2, lua.LTTable, h.Type())
				assert.Equal(t2, "b", h.(*lua.LTable).RawGetString("a").String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &response{
				StatusCode: tt.fields.StatusCode,
				Body:       tt.fields.Body,
				Headers:    tt.fields.Headers,
			}
			got := r.toLuaTable()
			tt.want(t, got)
		})
	}
}
