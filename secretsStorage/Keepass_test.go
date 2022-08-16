package secretsStorage

import (
	"envManager/internal"
	"github.com/stretchr/testify/assert"
	"github.com/tobischo/gokeepasslib/v3"
	"testing"
)

func TestKeepass_IsCaseSensitive(t *testing.T) {
	k := Keepass{}
	assert.Truef(t, k.IsCaseSensitive(), "IsCaseSensitive()")
}

func TestKeepass_Validate(t *testing.T) {
	existingFile := internal.GetTestDataFile(t, "keepass.kdbx")
	notExistingFile := internal.GetTestDataFile(t, "missing.kdbx")
	type fields struct {
		Name     string
		FilePath string
		database *gokeepasslib.Database
	}
	tests := []struct {
		name             string
		fields           fields
		wantError        bool
		wantErrorMessage string
		wantMessage      []string
	}{
		{
			name: "Existing file",
			fields: fields{
				Name:     "test_keepass",
				FilePath: existingFile,
				database: nil,
			},
			wantError:        false,
			wantErrorMessage: "",
			wantMessage:      []string{"Configured file is " + existingFile + "\nFile exists: true"},
		},
		{
			name: "Missing file",
			fields: fields{
				Name:     "test_keepass",
				FilePath: notExistingFile,
				database: nil,
			},
			wantError:        true,
			wantErrorMessage: "Validation of test_keepass failed. Run debug storage test_keepass to check it in detail",
			wantMessage:      []string{"Configured file is " + notExistingFile + "\nFile exists: false"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := Keepass{
				Name:     tt.fields.Name,
				FilePath: tt.fields.FilePath,
				database: tt.fields.database,
			}
			gotErr, gotMessages := k.Validate()

			if tt.wantError {
				assert.Error(t, gotErr, "Validate()")
			} else {
				assert.NoError(t, gotErr, "Validate()")
			}

			assert.Equalf(t, tt.wantMessage, gotMessages, "Validate()")
		})
	}
}
