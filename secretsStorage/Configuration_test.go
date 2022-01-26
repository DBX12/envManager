package secretsStorage

import (
	"bytes"
	"envManager/internal"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestConfiguration_LoadFromFile(t *testing.T) {
	//copyFixtureFile(t, "envManager.yml")
	type fields struct {
		Storages         map[string]Storage
		Profiles         map[string]Profile
		DirectoryMapping map[string][]string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *fields
	}{
		{
			name: "Non-existent file",
			fields: fields{
				Storages: map[string]Storage{},
				Profiles: map[string]Profile{},
			},
			args:    args{path: internal.GetTestDataFile(t, "not-existing.yml")},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Invalid yaml",
			fields: fields{
				Storages: map[string]Storage{},
				Profiles: map[string]Profile{},
			},
			args:    args{path: internal.GetTestDataFile(t, "invalid.yml")},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Existing file",
			fields: fields{
				Storages: map[string]Storage{},
				Profiles: map[string]Profile{},
			},
			args:    args{path: internal.GetTestDataFile(t, "envManager.yml")},
			wantErr: false,
			want: &fields{
				Storages: map[string]Storage{
					"keepass01": {
						StorageType: "keepass",
						Config: map[string]string{
							"path": "/tmp/keepass.kdbx",
						},
					},
				},
				Profiles: map[string]Profile{
					"root": {
						name:      "",
						Storage:   "keepass01",
						Path:      "entry1",
						ConstEnv:  map[string]string{"ROOT_PROF": "root_entry"},
						Env:       nil,
						DependsOn: nil,
					},
					"prof1": {
						name:    "",
						Storage: "keepass01",
						Path:    "group1/g1e1",
						ConstEnv: map[string]string{
							"PROF1_CONST": "foobar",
						},
						Env: map[string]string{
							"PROF1_USER": "UserName",
							"PROF1_PASS": "Password",
						},
						DependsOn: []string{"root"},
					},
				},
				DirectoryMapping: map[string][]string{
					"/tmp/projectA": {"prof1", "root"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Configuration{
				Storages: tt.fields.Storages,
				Profiles: tt.fields.Profiles,
			}
			if err := c.LoadFromFile(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				// do not test further if we wanted an error
				return
			}
			if !reflect.DeepEqual(tt.want.Storages, c.Storages) {
				t.Errorf("LoadFromFile() got = %v, want %v", c.Storages, tt.want.Storages)
			}
			if !reflect.DeepEqual(tt.want.Profiles, c.Profiles) {
				t.Errorf("LoadFromFile() got = %v, want %v", c.Profiles, tt.want.Profiles)
			}
		})
	}
}

func TestConfiguration_WriteToFile(t *testing.T) {
	type fields struct {
		Storages         map[string]Storage
		Profiles         map[string]Profile
		DirectoryMapping map[string][]string
	}
	type args struct {
		path    string
		replace bool
	}

	// prepare output directory
	outputDir, err := os.MkdirTemp(os.TempDir(), "Configuration_test")
	if err != nil {
		t.Fatal("Failed to prepare output directory")
	}
	// create a file for tests with existing files
	_, err = os.Create(path.Join(outputDir, "existing.yaml"))
	if err != nil {
		t.Fatal("Failed to create existing.yaml")
	}
	//goland:noinspection GoUnhandledErrorResult
	defer os.RemoveAll(outputDir)

	defaultFields := fields{
		Storages: map[string]Storage{
			"storage1": {
				StorageType: "dummy",
				Config: map[string]string{
					"key1": "val1",
				},
			},
		},
		Profiles: map[string]Profile{
			"profile1": {
				name:    "profile1",
				Storage: "storage1",
				Path:    "profile1",
				ConstEnv: map[string]string{
					"const1key": "const1value",
				},
				Env: map[string]string{
					"dyn1key": "dyn1value",
				},
				DependsOn: []string{
					"profile2", "profile3",
				},
			},
		},
		DirectoryMapping: map[string][]string{
			"/tmp/projectA": {"profile1"},
		},
	}

	tests := []struct {
		name            string
		fields          fields
		args            args
		fileShouldExist bool
		wantErr         bool
		wantedFile      string
	}{
		{
			name:   "Valid path",
			fields: defaultFields,
			args: args{
				path:    path.Join(outputDir, "test1.yaml"),
				replace: false,
			},
			fileShouldExist: false,
			wantErr:         false,
			wantedFile:      internal.GetTestDataFile(t, "expectedConfigFile1.yaml"),
		},
		{
			name:   "File exists",
			fields: defaultFields,
			args: args{
				path:    path.Join(outputDir, "existing.yaml"),
				replace: false,
			},
			fileShouldExist: true,
			wantErr:         true,
			wantedFile:      "",
		},
		{
			name:   "File exists and replace is set",
			fields: defaultFields,
			args: args{
				path:    path.Join(outputDir, "existing.yaml"),
				replace: true,
			},
			fileShouldExist: true,
			wantErr:         false,
			wantedFile:      internal.GetTestDataFile(t, "expectedConfigFile1.yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configuration{
				Storages:         tt.fields.Storages,
				Profiles:         tt.fields.Profiles,
				DirectoryMapping: tt.fields.DirectoryMapping,
			}

			// check that the test setup is correct
			fileInfo, _ := os.Stat(tt.args.path)
			if fileInfo == nil && tt.fileShouldExist {
				t.Fatalf("Expected file %s but it does not exist", tt.args.path)
			} else if fileInfo != nil && !tt.fileShouldExist {
				t.Fatalf("Expected no file at %s but there is one", tt.args.path)
			}
			err := c.WriteToFile(tt.args.path, tt.args.replace)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Got an error but wantend none. err = %v", err.Error())
				}
				// there was an error and we wanted one, nothing more to test
				return
			} else if tt.wantErr {
				t.Fatalf("Got no error but wantend one")
			}

			// check that file contents are equal
			expectedContent, _ := ioutil.ReadFile(tt.wantedFile)
			actualContent, _ := ioutil.ReadFile(tt.args.path)

			if !bytes.Equal(expectedContent, actualContent) {
				t.Errorf("File contents differ")
			}
		})
	}
}

func TestNewConfiguration(t *testing.T) {
	actual := NewConfiguration()
	if actual.Storages == nil {
		t.Error("storages property was not initialized")
	}
	if actual.Profiles == nil {
		t.Error("profiles property was not initialized")
	}
}
