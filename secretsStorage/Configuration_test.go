package secretsStorage

import (
	"bytes"
	"envManager/environment"
	"envManager/helper"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestConfiguration_LoadFromFile(t *testing.T) {
	//copyFixtureFile(t, "envManager.yml")
	type fields struct {
		Storages map[string]Storage
		Profiles map[string]Profile
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
			args:    args{path: helper.GetTestDataFile(t, "not-existing.yml")},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Invalid yaml",
			fields: fields{
				Storages: map[string]Storage{},
				Profiles: map[string]Profile{},
			},
			args:    args{path: helper.GetTestDataFile(t, "invalid.yml")},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Existing file",
			fields: fields{
				Storages: map[string]Storage{},
				Profiles: map[string]Profile{},
			},
			args:    args{path: helper.GetTestDataFile(t, "envManager.yml")},
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
		Storages map[string]Storage
		Profiles map[string]Profile
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
			wantedFile:      helper.GetTestDataFile(t, "expectedConfigFile1.yaml"),
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
			wantedFile:      helper.GetTestDataFile(t, "expectedConfigFile1.yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configuration{
				Storages: tt.fields.Storages,
				Profiles: tt.fields.Profiles,
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

func TestProfile_AddToEnvironment(t *testing.T) {
	type fields struct {
		name      string
		Storage   string
		Path      string
		ConstEnv  map[string]string
		Env       map[string]string
		DependsOn []string
	}
	type args struct {
		env *environment.Environment
	}

	// setup of dummy storage
	const storageName = "keepass"
	_ = GetRegistry().AddStorage(storageName, &Keepass{
		FilePath: helper.GetTestDataFile(t, "keepass.kdbx"),
	})

	// setup of dummy environments
	emptyEnv := environment.NewEnvironment()
	filledEnv := environment.NewEnvironment()
	_ = filledEnv.Set("preset_1", "preset_1_value")

	// pushing fixture passwords for dummy environments into input helper
	helper.GetInput().Inputs = []string{"1234"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Add to empty environment",
			fields: fields{
				name:    "the-profile",
				Storage: storageName,
				Path:    "entry1",
				ConstEnv: map[string]string{
					"CONST_ENV1": "CONST_VAL1",
				},
				Env: map[string]string{
					"user": "UserName",
					"pass": "Password",
				},
				DependsOn: nil,
			},
			args: args{
				env: &emptyEnv,
			},
			wantErr: false,
		},
		{
			name: "Add to filled environment",
			fields: fields{
				name:    "the-profile",
				Storage: storageName,
				Path:    "entry1",
				ConstEnv: map[string]string{
					"CONST_ENV1": "CONST_VAL1",
				},
				Env: map[string]string{
					"user": "UserName",
					"pass": "Password",
				},
				DependsOn: nil,
			},
			args: args{
				env: &filledEnv,
			},
			wantErr: false,
		},
		{
			name: "Set empty constEnv key",
			fields: fields{
				name:    "the-profile",
				Storage: storageName,
				Path:    "entry1",
				ConstEnv: map[string]string{
					"": "CONST_VAL1",
				},
				Env:       map[string]string{},
				DependsOn: nil,
			},
			args: args{
				env: &emptyEnv,
			},
			wantErr: true,
		},
		{
			name: "Set empty env key",
			fields: fields{
				name:     "the-profile",
				Storage:  storageName,
				Path:     "entry1",
				ConstEnv: map[string]string{},
				Env: map[string]string{
					"": "UserName",
				},
				DependsOn: nil,
			},
			args: args{
				env: &emptyEnv,
			},
			wantErr: true,
		},
		{
			name: "Use nonexistent storage",
			fields: fields{
				name:    "the-profile",
				Storage: "null",
				Path:    "entry1",
				ConstEnv: map[string]string{
					"CONST_ENV1": "CONST_VAL1",
				},
				Env: map[string]string{
					"user": "UserName",
					"pass": "Password",
				},
				DependsOn: nil,
			},
			args: args{
				env: &emptyEnv,
			},
			wantErr: true,
		},
		{
			name: "Use nonexistent entry",
			fields: fields{
				name:     "the-profile",
				Storage:  storageName,
				Path:     "null",
				ConstEnv: map[string]string{},
				Env: map[string]string{
					"user": "UserName",
					"pass": "Password",
				},
				DependsOn: nil,
			},
			args: args{
				env: &filledEnv,
			},
			wantErr: true,
		},
		{
			name: "Use nonexistent attribute",
			fields: fields{
				name:     "the-profile",
				Storage:  storageName,
				Path:     "entry1",
				ConstEnv: map[string]string{},
				Env: map[string]string{
					"user": "UserName",
					"pass": "null",
				},
				DependsOn: nil,
			},
			args: args{
				env: &filledEnv,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Profile{
				name:      tt.fields.name,
				Storage:   tt.fields.Storage,
				Path:      tt.fields.Path,
				ConstEnv:  tt.fields.ConstEnv,
				Env:       tt.fields.Env,
				DependsOn: tt.fields.DependsOn,
			}
			err := p.AddToEnvironment(tt.args.env)
			if err != nil {
				if tt.wantErr == false {
					t.Fatalf("AddToEnvironment() got error but wanted none. error = %v", err)
				}
				// got an error and wanted one, nothing more to test in this case
				return
			} else if tt.wantErr {
				t.Fatal("AddToEnvironment() got no error but wanted one.")
			}
		})
	}
}

func TestProfile_GetDependencies(t *testing.T) {
	type fields struct {
		name      string
		Storage   string
		Path      string
		ConstEnv  map[string]string
		Env       map[string]string
		DependsOn []string
	}
	type args struct {
		alreadyVisited []string
	}

	// setup of dummy registry
	const storageName = "keepass"
	_ = GetRegistry().AddStorage(storageName, &Keepass{
		FilePath: helper.GetTestDataFile(t, "keepass.kdbx"),
	})
	_ = GetRegistry().AddProfile("dependency1", getEmptyProfile())
	_ = GetRegistry().AddProfile("dependency2", Profile{
		name:      "dependency2",
		Storage:   storageName,
		Path:      "entry1",
		ConstEnv:  map[string]string{},
		Env:       map[string]string{},
		DependsOn: []string{"dependency1"},
	})
	_ = GetRegistry().AddProfile("circular", Profile{
		name:     "dependency2",
		Storage:  storageName,
		Path:     "entry1",
		ConstEnv: map[string]string{},
		Env:      map[string]string{},
		DependsOn: []string{
			"profile1",
		},
	})
	_ = GetRegistry().AddProfile("bad_dependency", Profile{
		name:     "dependency2",
		Storage:  storageName,
		Path:     "entry1",
		ConstEnv: map[string]string{},
		Env:      map[string]string{},
		DependsOn: []string{
			"null",
		},
	})

	/*
		Registry contents:
		- dependency1 (no deps)
		- dependency2 (depends on dependency1)
		- circular (depends on profile1, the profile for which dependencies are searched)
		- bad_dependency (depends on null, which does not exist)
	*/

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Profile without dependencies",
			fields: fields{
				name:      "profile1",
				Storage:   storageName,
				Path:      "group1/g1e1",
				ConstEnv:  map[string]string{},
				Env:       map[string]string{},
				DependsOn: []string{},
			},
			args:    args{alreadyVisited: []string{}},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Profile with dependency",
			fields: fields{
				name:     "profile1",
				Storage:  storageName,
				Path:     "group1/g1e1",
				ConstEnv: map[string]string{},
				Env:      map[string]string{},
				DependsOn: []string{
					"dependency1",
				},
			},
			args:    args{alreadyVisited: []string{}},
			want:    []string{"dependency1"},
			wantErr: false,
		},
		{
			name: "Profile with deep dependency",
			fields: fields{
				name:     "profile1",
				Storage:  storageName,
				Path:     "group1/g1e1",
				ConstEnv: map[string]string{},
				Env:      map[string]string{},
				DependsOn: []string{
					"dependency2",
				},
			},
			args:    args{alreadyVisited: []string{}},
			want:    []string{"dependency2", "dependency1"},
			wantErr: false,
		},
		{
			name: "Profile with circular dependency",
			fields: fields{
				name:     "profile1",
				Storage:  storageName,
				Path:     "group1/g1e1",
				ConstEnv: map[string]string{},
				Env:      map[string]string{},
				DependsOn: []string{
					"circular",
				},
			},
			args:    args{alreadyVisited: []string{}},
			want:    []string{"circular"},
			wantErr: false,
		},
		{
			name: "Profile with unknown dependency",
			fields: fields{
				name:     "profile1",
				Storage:  storageName,
				Path:     "group1/g1e1",
				ConstEnv: map[string]string{},
				Env:      map[string]string{},
				DependsOn: []string{
					"null",
				},
			},
			args:    args{alreadyVisited: []string{}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Profile with unknown dependency in dependency",
			fields: fields{
				name:     "profile1",
				Storage:  storageName,
				Path:     "group1/g1e1",
				ConstEnv: map[string]string{},
				Env:      map[string]string{},
				DependsOn: []string{
					"bad_dependency",
				},
			},
			args:    args{alreadyVisited: []string{}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				name:      tt.fields.name,
				Storage:   tt.fields.Storage,
				Path:      tt.fields.Path,
				ConstEnv:  tt.fields.ConstEnv,
				Env:       tt.fields.Env,
				DependsOn: tt.fields.DependsOn,
			}
			got, err := p.GetDependencies(tt.args.alreadyVisited)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDependencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDependencies() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfile_RemoveFromEnvironment(t *testing.T) {
	type fields struct {
		name      string
		Storage   string
		Path      string
		ConstEnv  map[string]string
		Env       map[string]string
		DependsOn []string
	}
	type args struct {
		env *environment.Environment
	}

	// setup of dummy storage
	const storageName = "keepass"
	_ = GetRegistry().AddStorage(storageName, &Keepass{
		FilePath: helper.GetTestDataFile(t, "keepass.kdbx"),
	})

	// setup of dummy environments
	emptyEnv := environment.NewEnvironment()
	filledEnv := environment.NewEnvironment()
	_ = filledEnv.Set("preset_1", "preset_1_value")

	// pushing fixture passwords for dummy environments into input helper
	helper.GetInput().Inputs = []string{"1234"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Remove from empty environment",
			fields: fields{
				name:    "the-profile",
				Storage: storageName,
				Path:    "entry1",
				ConstEnv: map[string]string{
					"CONST_ENV1": "CONST_VAL1",
				},
				Env: map[string]string{
					"user": "UserName",
					"pass": "Password",
				},
				DependsOn: nil,
			},
			args: args{
				env: &emptyEnv,
			},
			wantErr: false,
		},
		{
			name: "Add to filled environment",
			fields: fields{
				name:    "the-profile",
				Storage: storageName,
				Path:    "entry1",
				ConstEnv: map[string]string{
					"CONST_ENV1": "CONST_VAL1",
				},
				Env: map[string]string{
					"user": "UserName",
					"pass": "Password",
				},
				DependsOn: nil,
			},
			args: args{
				env: &filledEnv,
			},
			wantErr: false,
		},
		{
			name: "Set empty constEnv key",
			fields: fields{
				name:    "the-profile",
				Storage: storageName,
				Path:    "entry1",
				ConstEnv: map[string]string{
					"": "CONST_VAL1",
				},
				Env:       map[string]string{},
				DependsOn: nil,
			},
			args: args{
				env: &emptyEnv,
			},
			wantErr: true,
		},
		{
			name: "Set empty env key",
			fields: fields{
				name:     "the-profile",
				Storage:  storageName,
				Path:     "entry1",
				ConstEnv: map[string]string{},
				Env: map[string]string{
					"": "UserName",
				},
				DependsOn: nil,
			},
			args: args{
				env: &emptyEnv,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Profile{
				name:      tt.fields.name,
				Storage:   tt.fields.Storage,
				Path:      tt.fields.Path,
				ConstEnv:  tt.fields.ConstEnv,
				Env:       tt.fields.Env,
				DependsOn: tt.fields.DependsOn,
			}
			if err := p.RemoveFromEnvironment(tt.args.env); (err != nil) != tt.wantErr {
				t.Errorf("RemoveFromEnvironment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfile_SetName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		profile Profile
		args    args
		want    string
	}{
		{
			name:    "Set empty name",
			profile: getEmptyProfile(),
			args:    args{name: ""},
			want:    "",
		},
		{
			name:    "Set name",
			profile: getEmptyProfile(),
			args:    args{name: "the-name"},
			want:    "the-name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.profile.SetName(tt.args.name)
			if !reflect.DeepEqual(tt.want, tt.profile.name) {
				t.Errorf("SetName() = %v, want %v", tt.profile.name, tt.want)
			}
		})
	}
}

func TestProfile_Validate(t *testing.T) {
	type fields struct {
		name      string
		Storage   string
		Path      string
		ConstEnv  map[string]string
		Env       map[string]string
		DependsOn []string
	}

	// setup of dummy storage
	const storageName = "keepass"
	_ = GetRegistry().AddStorage(storageName, &Keepass{
		FilePath: helper.GetTestDataFile(t, "keepass.kdbx"),
	})
	_ = GetRegistry().AddProfile("dependency1", getEmptyProfile())

	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Valid profile with dependencies",
			fields: fields{
				name:     "the-profile",
				Storage:  storageName,
				Path:     "entry1",
				ConstEnv: nil,
				Env:      nil,
				DependsOn: []string{
					"dependency1",
				},
			},
			want: nil,
		},
		{
			name: "Valid profile without dependencies",
			fields: fields{
				name:      "the-profile",
				Storage:   storageName,
				Path:      "entry1",
				ConstEnv:  nil,
				Env:       nil,
				DependsOn: nil,
			},
			want: nil,
		},
		{
			name: "Invalid storage",
			fields: fields{
				name:      "the-profile",
				Storage:   "null",
				Path:      "entry1",
				ConstEnv:  nil,
				Env:       nil,
				DependsOn: nil,
			},
			want: []string{"references storage null which is not defined"},
		},
		{
			name: "Invalid dependency",
			fields: fields{
				name:      "the-profile",
				Storage:   storageName,
				Path:      "entry1",
				ConstEnv:  nil,
				Env:       nil,
				DependsOn: []string{"null"},
			},
			want: []string{"depends on null which is not defined"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				name:      tt.fields.name,
				Storage:   tt.fields.Storage,
				Path:      tt.fields.Path,
				ConstEnv:  tt.fields.ConstEnv,
				Env:       tt.fields.Env,
				DependsOn: tt.fields.DependsOn,
			}
			if got := p.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getEmptyProfile() Profile {
	return Profile{
		name:      "dummy profile",
		Storage:   "keepass",
		Path:      "entry1",
		ConstEnv:  map[string]string{},
		Env:       map[string]string{},
		DependsOn: []string{},
	}
}
