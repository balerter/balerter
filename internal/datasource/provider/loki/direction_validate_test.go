package loki

import "testing"

func Test_directionValidate(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				v: "",
			},
			wantErr: false,
		},
		{
			name: "forward",
			args: args{
				v: "forward",
			},
			wantErr: false,
		},
		{
			name: "backward",
			args: args{
				v: "backward",
			},
			wantErr: false,
		},
		{
			name: "wrong",
			args: args{
				v: "wrong",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := directionValidate(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("directionValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
