package secretsStorage

import (
	"errors"
	"fmt"
)

// Entry is a storage independent representation of an entry
type Entry struct {
	attributes map[string]string
}

// NewEntry instantiates an Entry object
func NewEntry() Entry {
	return Entry{
		attributes: map[string]string{},
	}
}

// SetAttribute sets an attribute in the entry. It will return an error if the
// key is an empty string.
func (e *Entry) SetAttribute(key string, value string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	e.attributes[key] = value
	return nil
}

// GetAttribute retrieves an attribute from this entry. It will return an error
// if the key is an empty string or does not exist.
func (e *Entry) GetAttribute(key string) (*string, error) {
	if key == "" {
		return nil, errors.New("key must not be empty")
	}
	value, exists := e.attributes[key]
	if exists == false {
		return nil, errors.New(fmt.Sprintf("unknown attribute %s", key))
	}
	return &value, nil
}

// GetAttributeNames returns a slice containing keys of the attributes of this entry.
func (e Entry) GetAttributeNames() []string {
	var out []string
	for key, _ := range e.attributes {
		out = append(out, key)
	}
	return out
}
