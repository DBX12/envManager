package secretsStorage

import (
	"envManager/environment"
	"envManager/helper"
	"reflect"
	"testing"
)

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
	_ = GetRegistry().AddProfile("dependency1", getEmptyProfile(t))
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
			profile: getEmptyProfile(t),
			args:    args{name: ""},
			want:    "",
		},
		{
			name:    "Set name",
			profile: getEmptyProfile(t),
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
	_ = GetRegistry().AddProfile("dependency1", getEmptyProfile(t))

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

func getEmptyProfile(t *testing.T) Profile {
	t.Helper()
	return Profile{
		name:      "dummy profile",
		Storage:   "keepass",
		Path:      "entry1",
		ConstEnv:  map[string]string{},
		Env:       map[string]string{},
		DependsOn: []string{},
	}
}
