package secretsStorage

import (
	"envManager/internal"
	"reflect"
	"testing"
)

func TestEntry_GetAttribute(t *testing.T) {
	type args struct {
		key string
	}
	johnDoe := "john.doe"
	tests := []struct {
		name              string
		initialAttributes map[string]string
		args              args
		want              *string
		wantErr           bool
	}{
		{
			name: "Get existing attribute",
			initialAttributes: map[string]string{
				"username": "john.doe",
				"password": "secret123",
			},
			args:    args{key: "username"},
			want:    &johnDoe,
			wantErr: false,
		},
		{
			name: "Get not existing attribute",
			initialAttributes: map[string]string{
				"username": "john.doe",
				"password": "secret123",
			},
			args:    args{key: "otherfield"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get existing attribute with space",
			initialAttributes: map[string]string{
				"user name": "john.doe",
				"pass word": "secret123",
			},
			args:    args{key: "user name"},
			want:    &johnDoe,
			wantErr: false,
		},
		{
			name: "Get not existing attribute with space",
			initialAttributes: map[string]string{
				"user name": "john.doe",
				"pass word": "secret123",
			},
			args:    args{key: "other field"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get attribute with empty key",
			initialAttributes: map[string]string{
				"username": "john.doe",
				"password": "secret123",
			},
			args:    args{key: ""},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				attributes: tt.initialAttributes,
			}
			got, err := e.GetAttribute(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAttribute() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_SetAttribute(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name              string
		initialAttributes map[string]string
		args              args
		want              map[string]string
		wantErr           bool
	}{
		{
			name:              "Set with empty key",
			initialAttributes: map[string]string{},
			args: args{
				key:   "",
				value: "vim",
			},
			want:    map[string]string{},
			wantErr: true,
		},
		{
			name: "Set with non-existing key",
			initialAttributes: map[string]string{
				"username": "john.doe",
			},
			args: args{
				key:   "password",
				value: "secret123",
			},
			want: map[string]string{
				"username": "john.doe",
				"password": "secret123",
			},
			wantErr: false,
		},
		{
			name: "Set with existing key",
			initialAttributes: map[string]string{
				"username": "john.doe",
			},
			args: args{
				key:   "username",
				value: "jack.daniels",
			},
			want: map[string]string{
				"username": "jack.daniels",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{
				attributes: tt.initialAttributes,
			}
			if err := e.SetAttribute(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SetAttribute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.want, e.attributes) {
				t.Errorf("e.attributes = %v, want %v", e.attributes, tt.want)
			}
		})
	}
}

func TestNewEntry(t *testing.T) {
	entry := NewEntry()
	if entry.attributes == nil {
		t.Error("attributes was not initialized")
	}
}

func TestEntry_GetAttributeNames(t *testing.T) {
	type fields struct {
		attributes map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Test empty",
			fields: fields{
				attributes: map[string]string{},
			},
			want: nil,
		},
		{
			name: "Test with attributes",
			fields: fields{
				attributes: map[string]string{
					"UserName": "john.doe",
					"Password": "secret123",
				},
			},
			want: []string{"UserName", "Password"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Entry{
				attributes: tt.fields.attributes,
			}
			got := e.GetAttributeNames()
			if !internal.AssertStringSliceEqual(t, tt.want, got) {
				t.Errorf("GetAttributeNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
