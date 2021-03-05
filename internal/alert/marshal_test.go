package alert

import (
	"testing"
	"time"
)

func TestAlert_Marshal(t *testing.T) {
	type fields struct {
		Name       string
		Level      Level
		LastChange time.Time
		Start      time.Time
		Count      int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "marshal",
			fields: fields{
				Name:       "1",
				Level:      2,
				LastChange: time.Date(2020, 01, 02, 03, 04, 05, 00, time.UTC),
				Start:      time.Date(2021, 01, 02, 03, 04, 05, 00, time.UTC),
				Count:      3,
			},
			want: `{"name":"1","level":"warning","level_num":2,"count":3,"last_change":"2020-01-02T03:04:05Z","start":"2021-01-02T03:04:05Z"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Alert{
				Name:       tt.fields.Name,
				Level:      tt.fields.Level,
				LastChange: tt.fields.LastChange,
				Start:      tt.fields.Start,
				Count:      tt.fields.Count,
			}
			if got := a.Marshal(); string(got) != tt.want {
				t.Errorf("Marshal() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestAlerts_Marshal(t *testing.T) {
	tests := []struct {
		name string
		a    Alerts
		want string
	}{
		{
			name: "marshal",
			a: Alerts{
				{
					Name:       "1",
					Level:      2,
					LastChange: time.Date(2020, 01, 02, 03, 04, 05, 00, time.UTC),
					Start:      time.Date(2021, 01, 02, 03, 04, 05, 00, time.UTC),
					Count:      3,
				},
				{
					Name:       "10",
					Level:      2,
					LastChange: time.Date(2020, 01, 02, 03, 04, 05, 00, time.UTC),
					Start:      time.Date(2021, 01, 02, 03, 04, 05, 00, time.UTC),
					Count:      30,
				},
			},
			want: `[{"name":"1","level":"warning","level_num":2,"count":3,"last_change":"2020-01-02T03:04:05Z","start":"2021-01-02T03:04:05Z"},{"name":"10","level":"warning","level_num":2,"count":30,"last_change":"2020-01-02T03:04:05Z","start":"2021-01-02T03:04:05Z"}]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Marshal(); string(got) != tt.want {
				t.Errorf("Marshal() = %s, want %s", got, tt.want)
			}
		})
	}
}
