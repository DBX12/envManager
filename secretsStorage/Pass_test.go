package secretsStorage

import (
	"context"
	"envManager/internal"
	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/gopasspw/gopass/pkg/gopass/secrets"
	"github.com/stretchr/testify/assert"
	"gopkg.in/errgo.v2/fmt/errors"
	"reflect"
	"testing"
)

func TestPass_GetDefaultConfig(t *testing.T) {
	p := &Pass{}
	want := map[string]string{
		"prefix": "",
	}
	got := p.GetDefaultConfig()
	if reflect.DeepEqual(want, got) == false {
		t.Errorf("GetDefaultConfig() got = %v, want %v", got, want)
	}
}

func TestPass_GetEntry_successful(t *testing.T) {
	const key = "key1"
	goPassMock := new(internal.MockGoPass)
	// return a gopass.Secret with the password "pass" and the attribute
	// "username" set to "john.doe"
	goPassMock.On(
		"Get",
		context.Background(),
		key,
		"",
	).Return(
		secrets.NewKVWithData(
			"pass",
			map[string][]string{
				"username": {"john.doe"},
			},
			"",
			false,
		),
		nil,
	)
	p := &Pass{
		store: goPassMock,
	}
	want := &Entry{map[string]string{
		"password": "pass",
		"username": "john.doe",
	}}

	got, err := p.GetEntry(key)

	assert.NoError(t, err, "No error occurred")
	assert.Equal(t, want, got, "Entry is loaded correctly")
}

// tests the behavior when p.store.Get() returns an error
func TestPass_GetEntry_getFailure(t *testing.T) {
	const key = "key1"
	goPassMock := new(internal.MockGoPass)
	// return a gopass.Secret with the password "pass" and the attribute
	// "username" set to "john.doe"
	goPassMock.On(
		"Get",
		context.Background(),
		key,
		"",
	).Return(
		nil,
		errors.New("Something went wrong"),
	)
	p := &Pass{
		store: goPassMock,
	}
	got, gotErr := p.GetEntry(key)
	assert.Error(t, gotErr, "Received an error")
	assert.Nil(t, got, "No Entry on error")
}

// tests the behavior when secret.Get(sKey) returns success = false
func TestPass_GetEntry_getKeyFailure(t *testing.T) {
	const entryName = "key1"
	goPassMock := new(internal.MockGoPass)
	goSecretMock := new(internal.MockGoSecret)
	goSecretMock.On("Keys").Return([]string{"username"})
	goSecretMock.On("Password").Return("password")
	goSecretMock.On("Get", "username").Return("", false)
	// return a gopass.Secret with the password "pass" and the attribute
	// "username" set to "john.doe"
	goPassMock.On(
		"Get",
		context.Background(),
		entryName,
		"",
	).Return(goSecretMock, nil)
	p := &Pass{
		store: goPassMock,
	}
	got, gotErr := p.GetEntry(entryName)
	assert.Error(t, gotErr, "Received an error")
	assert.Nil(t, got, "No Entry on error")
}

// tests behavior when Pass.Prefix is set
func TestPass_GetEntry_successfulWithPrefix(t *testing.T) {
	const key = "key1"
	goPassMock := new(internal.MockGoPass)
	// return a gopass.Secret with the password "pass" and the attribute
	// "username" set to "john.doe"
	goPassMock.On(
		"Get",
		context.Background(),
		"personal/"+key,
		"",
	).Return(
		secrets.NewKVWithData(
			"pass",
			map[string][]string{
				"username": {"john.doe"},
			},
			"",
			false,
		),
		nil,
	)
	p := &Pass{
		Prefix: "personal",
		store:  goPassMock,
	}
	want := &Entry{map[string]string{
		"password": "pass",
		"username": "john.doe",
	}}

	got, err := p.GetEntry(key)

	assert.NoError(t, err, "No error occurred")
	assert.Equal(t, want, got, "Entry is loaded correctly")
}

func TestPass_IsCaseSensitive(t *testing.T) {
	p := &Pass{}
	if p.IsCaseSensitive() != true {
		t.Error("IsCaseSensitive() should be true")
	}
}

func TestPass_Validate(t *testing.T) {
	p := &Pass{}
	gotErr, gotSlice := p.Validate()

	if gotErr != nil {
		t.Error("Validate() should not return an error")
	}

	if len(gotSlice) != 0 {
		t.Error("Validate() should not return any error message")
	}
}

func TestPass_initStore(t *testing.T) {
	type fields struct {
		store gopass.Store
	}
	store := &api.Gopass{}
	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		wantStore gopass.Store
	}{
		{
			name:      "Store is nil",
			fields:    fields{store: nil},
			wantErr:   false,
			wantStore: nil,
		},
		{
			// An existing store must not be replaced
			name:      "Store already initialized",
			fields:    fields{store: store},
			wantErr:   false,
			wantStore: store,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pass{
				store: tt.fields.store,
			}
			err := p.initStore()
			if tt.wantErr {
				assert.Error(t, err, "Want error")
			}
			assert.NotNil(t, p.store, "Store is initialized")
			if tt.wantStore != nil && p.store != tt.wantStore {
				t.Fatalf("initStore() store was overwritten, got = %p, want = %p", p.store, tt.wantStore)
			}
		})
	}
}
