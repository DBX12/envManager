package environment

import (
	"os"
	"reflect"
	"testing"
)

func TestEnvironment_Load(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
	}{
		{
			name:    "Empty environment",
			envVars: map[string]string{},
		},
		{
			name: "One environment variable",
			envVars: map[string]string{
				"FOO_PATH": "/tmp/foo",
			},
		},
		{
			name: "Multiple environment variables",
			envVars: map[string]string{
				"FOO_PATH":  "/tmp/foo",
				"BAR_COUNT": "5",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := prepareEnv(true, tt.envVars)
			if err != nil {
				t.Fatalf("Failed to prepare environment for test %s with error: %s", t.Name(), err.Error())
			}
			e := NewEnvironment()
			e.Load()
			if !reflect.DeepEqual(e.current, tt.envVars) {
				t.Errorf("Wanted %#v but got %#v", tt.envVars, e.current)
			}
		})
	}
}

func TestEnvironment_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "Set string pair",
			args: args{
				key:   "EDITOR",
				value: "vi",
			},
			want: map[string]string{
				"EDITOR": "vi",
			},
			wantErr: false,
		},
		{
			name: "Set empty value",
			args: args{
				key:   "EDITOR",
				value: "",
			},
			want: map[string]string{
				"EDITOR": "",
			},
			wantErr: false,
		},
		{
			name: "Set empty key",
			args: args{
				key:   "",
				value: "vi",
			},
			want:    map[string]string{},
			wantErr: true,
		},
		{
			name: "Set empty key and value",
			args: args{
				key:   "",
				value: "",
			},
			want:    map[string]string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEnvironment()
			err := e.Set(tt.args.key, tt.args.value)
			if err != nil && !tt.wantErr {
				t.Error("Got an error but wanted none")
			} else if err == nil && tt.wantErr {
				t.Error("Got no error but wanted one")
			}
			if !reflect.DeepEqual(e.addVars, tt.want) {
				t.Errorf("Wanted %#v but got %#v", tt.want, e.addVars)
			}
		})
	}
}

func TestEnvironment_Unset(t *testing.T) {
	type fields struct {
		current map[string]string
		addVars map[string]string
		delVars map[string]bool
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    fields
		wantErr bool
	}{
		{
			name: "Remove not added, not removed var",
			fields: fields{
				current: map[string]string{},
				addVars: map[string]string{},
				delVars: map[string]bool{},
			},
			args: args{
				key: "EDITOR",
			},
			want: fields{
				current: map[string]string{},
				addVars: map[string]string{},
				delVars: map[string]bool{"EDITOR": true},
			},
			wantErr: false,
		},
		{
			name: "Remove added, not removed var",
			fields: fields{
				current: map[string]string{},
				addVars: map[string]string{
					"EDITOR": "vim",
				},
				delVars: map[string]bool{},
			},
			args: args{
				key: "EDITOR",
			},
			want: fields{
				current: map[string]string{},
				addVars: map[string]string{},
				delVars: map[string]bool{"EDITOR": true},
			},
			wantErr: false,
		},
		{
			name: "Remove not added, removed var",
			fields: fields{
				current: map[string]string{},
				addVars: map[string]string{},
				delVars: map[string]bool{"EDITOR": true},
			},
			args: args{
				key: "EDITOR",
			},
			want: fields{
				current: map[string]string{},
				addVars: map[string]string{},
				delVars: map[string]bool{"EDITOR": true},
			},
			wantErr: false,
		},
		{
			name: "Remove empty string",
			fields: fields{
				current: map[string]string{},
				addVars: map[string]string{},
				delVars: map[string]bool{},
			},
			args: args{
				key: "",
			},
			want: fields{
				current: map[string]string{},
				addVars: map[string]string{},
				delVars: map[string]bool{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Environment{
				current: tt.fields.current,
				addVars: tt.fields.addVars,
				delVars: tt.fields.delVars,
			}
			err := e.Unset(tt.args.key)
			if err != nil && !tt.wantErr {
				t.Error("Got an error but wanted none")
			} else if err == nil && tt.wantErr {
				t.Error("Got no error but wanted one")
			}
			if !reflect.DeepEqual(tt.want.current, e.current) {
				t.Errorf("Unset() current = %v, want %v", e.current, tt.want.current)
			}
			if !reflect.DeepEqual(tt.want.addVars, e.addVars) {
				t.Errorf("Unset() addVars = %v, want %v", e.addVars, tt.want.addVars)
			}
			if !reflect.DeepEqual(tt.want.delVars, e.delVars) {
				t.Errorf("Unset() delVars = %v, want %v", e.delVars, tt.want.delVars)
			}
		})
	}
}

func TestEnvironment_WriteStatements(t *testing.T) {
	e := NewEnvironment()
	_ = e.Set("FOO_PATH", "/tmp/foo")
	_ = e.Set("BAR_PATH", "/tmp/bar")
	_ = e.Unset("EDITOR")
	want := "export FOO_PATH=\"/tmp/foo\";export BAR_PATH=\"/tmp/bar\";unset EDITOR"
	got := e.WriteStatements()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("WriteStatements() = %v, want %v", got, want)
	}
}

func TestNewEnvironment(t *testing.T) {
	e := NewEnvironment()
	if e.addVars == nil {
		t.Errorf("addVars was not initialized")
	}
	if e.delVars == nil {
		t.Errorf("delVars was not initialized")
	}
	if e.current == nil {
		t.Errorf("current was not initialized")
	}
}

func prepareEnv(clearCurrent bool, newVars map[string]string) error {
	if clearCurrent {
		os.Clearenv()
	}
	for name, value := range newVars {
		err := os.Setenv(name, value)
		if err != nil {
			return err
		}
	}
	return nil
}
