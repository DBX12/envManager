package secretsStorage

import (
	"reflect"
	"testing"
)

func TestGetRegistry(t *testing.T) {
	actual := GetRegistry()
	if actual == nil {
		t.Fatal("Did not create instance")
	}
}

func TestRegistry_AddProfile(t *testing.T) {
	type fields struct {
		storages map[string]StorageAdapter
		profiles map[string]Profile
	}
	type args struct {
		name    string
		profile Profile
	}
	emptyRegistry := fields{
		storages: map[string]StorageAdapter{},
		profiles: map[string]Profile{},
	}
	validProfile := Profile{
		Storage:   "test01",
		Path:      "group1/entry1",
		ConstEnv:  map[string]string{"const1": "cval1"},
		Env:       map[string]string{"dynamic1": "dval1"},
		DependsOn: []string{},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Empty name",
			fields: emptyRegistry,
			args: args{
				name:    "",
				profile: validProfile,
			},
			wantErr: true,
		},
		{
			name:   "Non-empty name",
			fields: emptyRegistry,
			args: args{
				name:    "profile1",
				profile: validProfile,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Registry{
				storages: tt.fields.storages,
				profiles: tt.fields.profiles,
			}
			err := r.AddProfile(tt.args.name, tt.args.profile)
			if err != nil {
				if !tt.wantErr {
					t.Error("AddProfile() didn't want error but got one")
				}
			} else {
				if tt.wantErr {
					t.Error("AddProfile() wanted error but got none")
				}
				profile, _ := r.GetProfile(tt.args.name)
				if !reflect.DeepEqual(profile.name, tt.args.name) {
					t.Error("Name was not injected into profile")
				}
			}
		})
	}
}

func TestRegistry_AddStorage(t *testing.T) {
	type fields struct {
		storages map[string]StorageAdapter
		profiles map[string]Profile
	}
	type args struct {
		name    string
		storage StorageAdapter
	}
	emptyRegistry := fields{
		storages: map[string]StorageAdapter{},
		profiles: map[string]Profile{},
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantStorage bool
	}{
		{
			name:   "Add with empty name",
			fields: emptyRegistry,
			args: args{
				name: "",
				storage: &Keepass{
					Name:     "keepass01",
					FilePath: "/tmp/keepass.kdbx",
				},
			},
			wantErr:     true,
			wantStorage: false,
		},
		{
			name:   "Add nil",
			fields: emptyRegistry,
			args: args{
				name:    "nil-storage",
				storage: nil,
			},
			wantErr:     true,
			wantStorage: false,
		},
		{
			name:   "Add valid storage",
			fields: emptyRegistry,
			args: args{
				name: "keepass01",
				storage: &Keepass{
					Name:     "keepass01",
					FilePath: "/tmp/keepass.kdbx",
				},
			},
			wantErr:     false,
			wantStorage: true,
		},
		{
			name: "Add valid storage when one is existing",
			fields: fields{
				storages: map[string]StorageAdapter{
					"keepass01": &Keepass{
						Name:     "keepass01",
						FilePath: "/tmp/initialFile.kdbx",
					},
				},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "keepass01",
				storage: &Keepass{
					Name:     "keepass01",
					FilePath: "/tmp/newFile.kdbx",
				},
			},
			wantErr:     false,
			wantStorage: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Registry{
				storages: tt.fields.storages,
				profiles: tt.fields.profiles,
			}
			if err := r.AddStorage(tt.args.name, tt.args.storage); (err != nil) != tt.wantErr {
				t.Errorf("AddStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if r.HasStorage(tt.args.name) {
				if !tt.wantStorage {
					t.Error("The storage exists in the registry while it should not")
				}
				actualStorage, _ := r.GetStorage(tt.args.name)
				if !reflect.DeepEqual(*actualStorage, tt.args.storage) {
					t.Errorf("Storage in registry = %v, want %v", actualStorage, tt.args.storage)
				}
			} else if tt.wantStorage {
				t.Error("The storage does not exist in the registry while it should")
			}
		})
	}
}

func TestRegistry_GetAllProfiles(t *testing.T) {
	type fields struct {
		storages map[string]StorageAdapter
		profiles map[string]Profile
	}
	emptyRegistry := fields{
		storages: map[string]StorageAdapter{},
		profiles: map[string]Profile{},
	}
	validProfile := Profile{
		name:      "profile1",
		Storage:   "test01",
		Path:      "group1/entry1",
		ConstEnv:  map[string]string{"const1": "cval1"},
		Env:       map[string]string{"dynamic1": "dval1"},
		DependsOn: []string{},
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]Profile
	}{
		{
			name:   "Empty registry",
			fields: emptyRegistry,
			want:   map[string]Profile{},
		},
		{
			name: "With profile",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{
					"profile1": validProfile,
				},
			},
			want: map[string]Profile{
				"profile1": validProfile,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Registry{
				storages: tt.fields.storages,
				profiles: tt.fields.profiles,
			}
			if got := r.GetAllProfiles(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllProfiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegistry_GetAllStorages(t *testing.T) {
	type fields struct {
		storages map[string]StorageAdapter
		profiles map[string]Profile
	}
	dummyStorages := map[string]StorageAdapter{
		"keepass0": &Keepass{
			Name:     "keepass0",
			FilePath: "/tmp/keepass0.kdbx",
		},
		"keepass1": &Keepass{
			Name:     "keepass1",
			FilePath: "/tmp/keepass1.kdbx",
		},
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]StorageAdapter
	}{
		{
			name: "Empty registry",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{},
			},
			want: map[string]StorageAdapter{},
		},
		{
			name: "With dummyStorages",
			fields: fields{
				storages: dummyStorages,
				profiles: map[string]Profile{},
			},
			want: dummyStorages,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Registry{
				storages: tt.fields.storages,
				profiles: tt.fields.profiles,
			}
			if got := r.GetAllStorages(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllStorages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegistry_GetProfile(t *testing.T) {
	type fields struct {
		storages map[string]StorageAdapter
		profiles map[string]Profile
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Profile
		wantErr bool
	}{
		{
			name: "Empty registry",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "awsMain",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get with empty name",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{
					"awsMain": {
						name:      "awsMain",
						Storage:   "keepass0",
						Path:      "aws/awsMain",
						ConstEnv:  map[string]string{},
						Env:       map[string]string{},
						DependsOn: []string{},
					},
				},
			},
			args: args{
				name: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get existing",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{
					"awsMain": {
						name:      "awsMain",
						Storage:   "keepass0",
						Path:      "aws/awsMain",
						ConstEnv:  map[string]string{},
						Env:       map[string]string{},
						DependsOn: []string{},
					},
				},
			},
			args: args{
				name: "awsMain",
			},
			want: &Profile{
				name:      "awsMain",
				Storage:   "keepass0",
				Path:      "aws/awsMain",
				ConstEnv:  map[string]string{},
				Env:       map[string]string{},
				DependsOn: []string{},
			},
			wantErr: false,
		},
		{
			name: "Get not existing",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{
					"awsMain": {
						name:      "awsMain",
						Storage:   "keepass0",
						Path:      "aws/awsMain",
						ConstEnv:  map[string]string{},
						Env:       map[string]string{},
						DependsOn: []string{},
					},
				},
			},
			args: args{
				name: "awsProd",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Registry{
				storages: tt.fields.storages,
				profiles: tt.fields.profiles,
			}
			got, err := r.GetProfile(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProfile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegistry_GetProfileNames(t *testing.T) {
	type fields struct {
		storages map[string]StorageAdapter
		profiles map[string]Profile
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Empty registry",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{},
			},
			want: nil,
		},
		{
			name: "Filled registry",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{
					"awsMain": {
						name:      "awsMain",
						Storage:   "keepass0",
						Path:      "aws/main",
						ConstEnv:  map[string]string{},
						Env:       map[string]string{},
						DependsOn: []string{},
					},
					"awsProd": {
						name:      "awsProd",
						Storage:   "keepass0",
						Path:      "aws/main",
						ConstEnv:  map[string]string{},
						Env:       map[string]string{},
						DependsOn: []string{},
					},
				},
			},
			want: []string{"awsMain", "awsProd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Registry{
				storages: tt.fields.storages,
				profiles: tt.fields.profiles,
			}
			if got := r.GetProfileNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProfileNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegistry_GetStorage(t *testing.T) {
	type fields struct {
		storages map[string]StorageAdapter
		profiles map[string]Profile
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    StorageAdapter
		wantErr bool
	}{
		{
			name: "Empty registry",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "keepass0",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get with empty name",
			fields: fields{
				storages: map[string]StorageAdapter{
					"keepass0": &Keepass{
						Name:     "keepass0",
						FilePath: "/tmp/keepass0.kdbx",
					},
				},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get existing",
			fields: fields{
				storages: map[string]StorageAdapter{
					"keepass0": &Keepass{
						Name:     "keepass0",
						FilePath: "/tmp/keepass0.kdbx",
					},
				},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "keepass0",
			},
			want: &Keepass{
				Name:     "keepass0",
				FilePath: "/tmp/keepass0.kdbx",
			},
			wantErr: false,
		},
		{
			name: "Get not existing",
			fields: fields{
				storages: map[string]StorageAdapter{
					"keepass0": &Keepass{
						Name:     "keepass0",
						FilePath: "/tmp/keepass0.kdbx",
					},
				},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "keepass1",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Registry{
				storages: tt.fields.storages,
				profiles: tt.fields.profiles,
			}
			got, err := r.GetStorage(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				if tt.want != nil {
					t.Fatal("Wanted non-nil value but got nil")
				}
			} else {
				// only dereference got if it is not nil
				if !reflect.DeepEqual(*got, tt.want) {
					t.Errorf("GetStorage() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRegistry_GetStorageNames(t *testing.T) {
	type fields struct {
		storages map[string]StorageAdapter
		profiles map[string]Profile
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Empty registry",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{},
			},
			want: nil,
		},
		{
			name: "Filled registry",
			fields: fields{
				storages: map[string]StorageAdapter{
					"keepass0": &Keepass{
						Name:     "keepass0",
						FilePath: "/tmp/keepass0.kdbx",
					},
					"keepass1": &Keepass{
						Name:     "keepass1",
						FilePath: "/tmp/keepass1.kdbx",
					},
				},
				profiles: map[string]Profile{},
			},
			want: []string{"keepass0", "keepass1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Registry{
				storages: tt.fields.storages,
				profiles: tt.fields.profiles,
			}
			if got := r.GetStorageNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStorageNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegistry_HasProfile(t *testing.T) {
	type fields struct {
		storages map[string]StorageAdapter
		profiles map[string]Profile
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Empty registry",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "awsMain",
			},
			want: false,
		},
		{
			name: "Empty profile name",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{
					"awsProd": {
						name:      "awsProd",
						Storage:   "keepass0",
						Path:      "aws/prod",
						ConstEnv:  map[string]string{},
						Env:       map[string]string{},
						DependsOn: []string{},
					},
				},
			},
			args: args{
				name: "",
			},
			want: false,
		},
		{
			name: "Unknown profile",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{
					"awsProd": {
						name:      "awsProd",
						Storage:   "keepass0",
						Path:      "aws/prod",
						ConstEnv:  map[string]string{},
						Env:       map[string]string{},
						DependsOn: []string{},
					},
				},
			},
			args: args{
				name: "awsMain",
			},
			want: false,
		},
		{
			name: "Known profile",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{
					"awsProd": {
						name:      "awsProd",
						Storage:   "keepass0",
						Path:      "aws/prod",
						ConstEnv:  map[string]string{},
						Env:       map[string]string{},
						DependsOn: []string{},
					},
				},
			},
			args: args{
				name: "awsProd",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Registry{
				storages: tt.fields.storages,
				profiles: tt.fields.profiles,
			}
			if got := r.HasProfile(tt.args.name); got != tt.want {
				t.Errorf("HasProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegistry_HasStorage(t *testing.T) {
	type fields struct {
		storages map[string]StorageAdapter
		profiles map[string]Profile
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Empty registry",
			fields: fields{
				storages: map[string]StorageAdapter{},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "keepass0",
			},
			want: false,
		},
		{
			name: "Empty storage name",
			fields: fields{
				storages: map[string]StorageAdapter{
					"keepass0": &Keepass{
						Name:     "keepass0",
						FilePath: "/tmp/keepass0.kdbx",
					},
				},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "",
			},
			want: false,
		},
		{
			name: "Unknown storage",
			fields: fields{
				storages: map[string]StorageAdapter{
					"keepass0": &Keepass{
						Name:     "keepass0",
						FilePath: "/tmp/keepass0.kdbx",
					},
				},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "keepass1",
			},
			want: false,
		},
		{
			name: "Known storage",
			fields: fields{
				storages: map[string]StorageAdapter{
					"keepass0": &Keepass{
						Name:     "keepass0",
						FilePath: "/tmp/keepass0.kdbx",
					},
				},
				profiles: map[string]Profile{},
			},
			args: args{
				name: "keepass0",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Registry{
				storages: tt.fields.storages,
				profiles: tt.fields.profiles,
			}
			if got := r.HasStorage(tt.args.name); got != tt.want {
				t.Errorf("HasStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newRegistry(t *testing.T) {
	registry := newRegistry()
	if registry == nil {
		t.Fatal("Did not create instance")
	}
	if registry.profiles == nil {
		t.Error("Did not initialize profiles map")
	}
	if registry.storages == nil {
		t.Error("Did not initialize storages map")
	}
}
