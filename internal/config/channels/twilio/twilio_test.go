package twilio

import "testing"

func TestTwilio_Validate(t *testing.T) {
	type fields struct {
		Name   string
		SID    string
		Token  string
		From   string
		To     string
		TwiML  string
		Ignore bool
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errValue string
	}{
		{
			name:     "empty name",
			fields:   fields{},
			wantErr:  true,
			errValue: "name must be not empty",
		},
		{
			name: "empty sid",
			fields: fields{
				Name:   "1",
				SID:    "",
				Token:  "",
				From:   "",
				To:     "",
				TwiML:  "",
				Ignore: false,
			},
			wantErr:  true,
			errValue: "sid must be not empty",
		},
		{
			name: "empty token",
			fields: fields{
				Name:   "1",
				SID:    "1",
				Token:  "",
				From:   "",
				To:     "",
				TwiML:  "",
				Ignore: false,
			},
			wantErr:  true,
			errValue: "token must be not empty",
		},
		{
			name: "empty from",
			fields: fields{
				Name:   "1",
				SID:    "1",
				Token:  "1",
				From:   "",
				To:     "",
				TwiML:  "",
				Ignore: false,
			},
			wantErr:  true,
			errValue: "from must be not empty",
		},
		{
			name: "empty to",
			fields: fields{
				Name:   "1",
				SID:    "1",
				Token:  "1",
				From:   "1",
				To:     "",
				TwiML:  "",
				Ignore: false,
			},
			wantErr:  true,
			errValue: "to must be not empty",
		},
		{
			name: "ok",
			fields: fields{
				Name:   "1",
				SID:    "1",
				Token:  "1",
				From:   "1",
				To:     "1",
				TwiML:  "",
				Ignore: false,
			},
			wantErr:  false,
			errValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := Twilio{
				Name:   tt.fields.Name,
				SID:    tt.fields.SID,
				Token:  tt.fields.Token,
				From:   tt.fields.From,
				To:     tt.fields.To,
				TwiML:  tt.fields.TwiML,
				Ignore: tt.fields.Ignore,
			}
			err := tw.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr && err.Error() != tt.errValue {
				t.Errorf("unexpected error value = %s, want %s", err.Error(), tt.errValue)
			}
		})
	}
}
