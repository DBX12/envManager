package secretsStorage

import (
	"context"
	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"gopkg.in/errgo.v2/fmt/errors"
)

const PassTypeIdentifier = "pass"

type Pass struct {
	Name   string
	Prefix string
	store  gopass.Store
}

func (p *Pass) GetEntry(key string) (*Entry, error) {
	if err := p.initStore(); err != nil {
		return nil, err
	}
	if p.Prefix != "" {
		// only add non-empty prefix, otherwise the key starts with /
		key = p.Prefix + "/" + key
	}
	secret, err := p.store.Get(context.Background(), key, "")
	if err != nil {
		return nil, err
	}
	entry := NewEntry()
	_ = entry.SetAttribute("password", secret.Password())
	for _, sKey := range secret.Keys() {
		value, success := secret.Get(sKey)
		if !success {
			return nil, errors.Newf("Got false when retrieving sKey %s on key %s", sKey, key)
		}
		err = entry.SetAttribute(sKey, value)
		if err != nil {
			return nil, err
		}
	}
	return &entry, nil
}
func (p *Pass) initStore() error {
	if p.store == nil {
		var err error
		p.store, err = api.New(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pass) IsCaseSensitive() bool {
	return true
}

func (p *Pass) Validate() (error, []string) {
	return nil, []string{}
}

func (p *Pass) GetDefaultConfig() map[string]string {
	return map[string]string{
		"prefix": "",
	}
}
