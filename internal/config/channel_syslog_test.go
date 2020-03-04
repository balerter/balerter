package config

import "testing"

func Test_validatePriority(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty string",
			args: args{
				p: "",
			},
			wantErr: false,
		},
		{
			name: "many parts",
			args: args{
				p: "EMERG|EMERG|EMERG",
			},
			wantErr: true,
		},
		{
			name: "correct severity",
			args: args{
				p: "EMERG",
			},
			wantErr: false,
		},
		{
			name: "incorrect severity",
			args: args{
				p: "EMERG-BAD",
			},
			wantErr: true,
		},
		{
			name: "correct facility",
			args: args{
				p: "EMERG|FTP",
			},
			wantErr: false,
		},
		{
			name: "incorrect facility",
			args: args{
				p: "EMERG|FTP-BAD",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePriority(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("validatePriority() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
