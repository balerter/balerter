package scripts

import (
	"github.com/balerter/balerter/internal/config/scripts/file"
	"github.com/balerter/balerter/internal/config/scripts/folder"
	"github.com/balerter/balerter/internal/config/scripts/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestScripts_Validate(t *testing.T) {
	type fields struct {
		UpdateInterval int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "small updateInterval",
			fields:  fields{UpdateInterval: -1},
			wantErr: true,
			errText: "updateInterval must be not less than 0",
		},
		{
			name:    "ok",
			fields:  fields{UpdateInterval: 1},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Scripts{
				UpdateInterval: tt.fields.UpdateInterval,
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

func TestScripts_Validate_check_name_folder(t *testing.T) {
	f := folder.Folder{Name: "1", Path: os.TempDir()}
	s := Scripts{
		Folder: []folder.Folder{f, f},
	}

	err := s.Validate()
	require.Error(t, err)
	assert.Equal(t, "found duplicated name for scritsource 'folder': 1", err.Error())
}

func TestScripts_Validate_check_name_file(t *testing.T) {
	ff, err := os.CreateTemp("", "")
	require.NoError(t, err)
	f := file.File{Name: "1", Filename: ff.Name()}
	s := Scripts{
		File: []file.File{f, f},
	}

	err = s.Validate()
	require.Error(t, err)
	assert.Equal(t, "found duplicated name for scritsource 'file': 1", err.Error())
}

func TestScripts_Validate_check_name_postgres(t *testing.T) {
	f := postgres.Postgres{
		Name:    "1",
		Host:    "2",
		Port:    10,
		Timeout: 10,
		Query:   "3",
	}
	s := Scripts{
		Postgres: []postgres.Postgres{f, f},
	}

	err := s.Validate()
	require.Error(t, err)
	assert.Equal(t, "found duplicated name for scritsource 'postgres': 1", err.Error())
}
