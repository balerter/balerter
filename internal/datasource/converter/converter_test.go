package converter

import (
	lua "github.com/yuin/gopher-lua"
	"reflect"
	"testing"
)

func TestFromDateBytes(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want lua.LValue
	}{
		{
			name: "empty",
			args: args{
				v: &[]byte{},
			},
			want: lua.LString(""),
		},
		{
			name: "string",
			args: args{
				v: &[]byte{0x23, 0x24, 0x25},
			},
			want: lua.LString("#$%"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromDateBytes(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromDateBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
