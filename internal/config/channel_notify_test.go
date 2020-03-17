package config

import "testing"

func TestChannelNotify_Validate(t *testing.T) {
	type fields struct {
		Name  string
		Icons ChannelNotifyIcons
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty name",
			fields: fields{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "with name",
			fields: fields{
				Name: "foo",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ChannelNotify{
				Name:  tt.fields.Name,
				Icons: tt.fields.Icons,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
