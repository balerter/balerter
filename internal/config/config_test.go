package config

import (
	"github.com/balerter/balerter/internal/util"
	"testing"
)

func Test_checkUnique(t *testing.T) {
	type args struct {
		data []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty slice",
			args: args{
				data: []string{},
			},
			want: "",
		},
		{
			name: "one element",
			args: args{
				data: []string{"foo"},
			},
			want: "",
		},
		{
			name: "without duplicates",
			args: args{
				data: []string{"foo", "bar"},
			},
			want: "",
		},
		{
			name: "with duplicates",
			args: args{
				data: []string{"foo", "bar", "baz", "foo"},
			},
			want: "foo",
		},
		{
			name: "with duplicates case insensitive",
			args: args{
				data: []string{"foo", "bar", "baz", "FOO"},
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.CheckUnique(tt.args.data); got != tt.want {
				t.Errorf("checkUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}
