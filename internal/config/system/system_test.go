package system

import "testing"

func TestSystem_Validate(t *testing.T) {
	type fields struct {
		JobWorkersCount int
		CronLocation    string
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errValue string
	}{
		{
			name:     "empty cron location",
			fields:   fields{},
			wantErr:  false,
			errValue: "",
		},
		{
			name: "correct cron location",
			fields: fields{
				JobWorkersCount: 0,
				CronLocation:    "UTC",
			},
			wantErr:  false,
			errValue: "",
		},
		{
			name: "incorrect cron location",
			fields: fields{
				JobWorkersCount: 0,
				CronLocation:    "foo",
			},
			wantErr:  true,
			errValue: "error parse cronLocation, unknown time zone foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &System{
				JobWorkersCount: tt.fields.JobWorkersCount,
				CronLocation:    tt.fields.CronLocation,
			}
			err := s.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errValue {
				t.Errorf("unexpectd error value = %s, wantErr %s", err.Error(), tt.errValue)
			}
		})
	}
}
