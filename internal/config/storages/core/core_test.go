package core

import (
	"github.com/balerter/balerter/internal/config/storages/core/sqlite"
	"github.com/balerter/balerter/internal/config/storages/core/tables"
	"testing"
)

func createSqliteItem(name string) sqlite.Sqlite {
	item := sqlite.Sqlite{
		Name:    name,
		Path:    "1",
		Timeout: 0,
		TableAlerts: tables.TableAlerts{
			Table: "1",
			Fields: tables.AlertFields{
				Name:      "1",
				Level:     "2",
				Count:     "3",
				UpdatedAt: "4",
				CreatedAt: "5",
			},
		},
		TableKV: tables.TableKV{
			Table: "1",
			Fields: tables.KVFields{
				Key:   "1",
				Value: "2",
			},
		},
	}

	return item
}

func TestStoragesCore_Validate(t *testing.T) {
	type fields struct {
		File []sqlite.Sqlite
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name: "duplicated storage name",
			fields: fields{
				File: []sqlite.Sqlite{
					createSqliteItem("a"),
					createSqliteItem("a"),
				},
			},
			wantErr: true,
			errText: "found duplicated name for core storages 'sqlite': a",
		},
		{
			name: "ok",
			fields: fields{
				File: []sqlite.Sqlite{
					createSqliteItem("a1"),
					createSqliteItem("a2"),
				},
			},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Core{
				Sqlite: tt.fields.File,
			}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errText {
				t.Errorf("Validate() error = '%s', wantErrText '%s'", err.Error(), tt.errText)
			}
		})
	}
}
