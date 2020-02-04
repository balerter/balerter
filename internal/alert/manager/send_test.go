package manager

import "testing"

func Test_contains(t *testing.T) {
	type args struct {
		v   string
		arr []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty array",
			args: args{
				v:   "foo",
				arr: nil,
			},
			want: false,
		},
		{
			name: "not exists",
			args: args{
				v:   "foo",
				arr: []string{"bar"},
			},
			want: false,
		},
		{
			name: "exists",
			args: args{
				v:   "foo",
				arr: []string{"bar", "foo", "baz"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.v, tt.args.arr); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
