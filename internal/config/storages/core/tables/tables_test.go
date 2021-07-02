package tables

import "testing"

func TestTableAlerts_Validate(t1 *testing.T) {
	type fields struct {
		Table  string
		Fields AlertFields
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errValue string
	}{
		{
			name:     "no table",
			fields:   fields{},
			wantErr:  true,
			errValue: "table must be not empty",
		},
		{
			name: "field name",
			fields: fields{
				Table: "1",
				Fields: AlertFields{
					Name:      "",
					Level:     "",
					Count:     "",
					UpdatedAt: "",
					CreatedAt: "",
				},
			},
			wantErr:  true,
			errValue: "field name must be not empty",
		},
		{
			name: "field level",
			fields: fields{
				Table: "1",
				Fields: AlertFields{
					Name:      "1",
					Level:     "",
					Count:     "",
					UpdatedAt: "",
					CreatedAt: "",
				},
			},
			wantErr:  true,
			errValue: "field level must be not empty",
		},
		{
			name: "field count",
			fields: fields{
				Table: "1",
				Fields: AlertFields{
					Name:      "1",
					Level:     "1",
					Count:     "",
					UpdatedAt: "",
					CreatedAt: "",
				},
			},
			wantErr:  true,
			errValue: "field count must be not empty",
		},
		{
			name: "field updated_at",
			fields: fields{
				Table: "1",
				Fields: AlertFields{
					Name:      "1",
					Level:     "1",
					Count:     "1",
					UpdatedAt: "",
					CreatedAt: "1",
				},
			},
			wantErr:  true,
			errValue: "field updatedAt must be not empty",
		},
		{
			name: "field created_at",
			fields: fields{
				Table: "1",
				Fields: AlertFields{
					Name:      "1",
					Level:     "1",
					Count:     "1",
					UpdatedAt: "1",
					CreatedAt: "",
				},
			},
			wantErr:  true,
			errValue: "field createdAt must be not empty",
		},
		{
			name: "ok",
			fields: fields{
				Table: "1",
				Fields: AlertFields{
					Name:      "1",
					Level:     "1",
					Count:     "1",
					UpdatedAt: "1",
					CreatedAt: "1",
				},
			},
			wantErr:  false,
			errValue: "",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TableAlerts{
				Table:  tt.fields.Table,
				Fields: tt.fields.Fields,
			}
			err := t.Validate()
			if (err != nil) != tt.wantErr {
				t1.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr && err.Error() != tt.errValue {
				t1.Errorf("unexpected error value '%s', expect '%s'", err.Error(), tt.errValue)
			}
		})
	}
}

func TestTableKV_Validate(t1 *testing.T) {
	type fields struct {
		Table       string
		Fields      KVFields
		CreateTable bool
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		errValue string
	}{
		{
			name: "no table",
			fields: fields{
				Table:       "",
				Fields:      KVFields{},
				CreateTable: false,
			},
			wantErr:  true,
			errValue: "table must be not empty",
		},
		{
			name: "field key",
			fields: fields{
				Table: "1",
				Fields: KVFields{
					Key:   "",
					Value: "",
				},
				CreateTable: false,
			},
			wantErr:  true,
			errValue: "field key must be not empty",
		},
		{
			name: "field value",
			fields: fields{
				Table: "1",
				Fields: KVFields{
					Key:   "1",
					Value: "",
				},
				CreateTable: false,
			},
			wantErr:  true,
			errValue: "field value must be not empty",
		},
		{
			name: "ok",
			fields: fields{
				Table: "1",
				Fields: KVFields{
					Key:   "1",
					Value: "1",
				},
				CreateTable: false,
			},
			wantErr:  false,
			errValue: "",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TableKV{
				Table:       tt.fields.Table,
				Fields:      tt.fields.Fields,
				CreateTable: tt.fields.CreateTable,
			}
			err := t.Validate()
			if (err != nil) != tt.wantErr {
				t1.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr && err.Error() != tt.errValue {
				t1.Errorf("unexpected error value '%s', expect '%s'", err.Error(), tt.errValue)
			}
		})
	}
}
