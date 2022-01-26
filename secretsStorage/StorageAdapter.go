package secretsStorage

import (
	"gopkg.in/errgo.v2/fmt/errors"
)

//StorageAdapter provides methods to interact with a secrets secretsStorage (e.g. keepass)
type StorageAdapter interface {
	//GetEntry retrieves an entry from the secretsStorage. The entry is addressed by the key parameter, it depends on the
	//implementation of the StorageAdapter how the key is interpreted.
	GetEntry(key string) (*Entry, error)
	//IsCaseSensitive indicates if the storage provider is case-sensitive (path and attributes)
	IsCaseSensitive() bool
	//Validate verifies that all provided data is valid (e.g. checks for existing files).
	//It returns an error value indicating if the validation was successful and a slice of strings containing information
	//for the user.
	Validate() (error, []string)
	//GetDefaultConfig returns the default config of the storage adapter. The return value of this function will be used
	//to initialize a storage adapter section in the config file.
	GetDefaultConfig() map[string]string
}

//CreateStorageAdapter is a factory method which creates a specific storage adapter determined by data["type"] and calls
//StorageAdapter.Validate on the created instance. Should StorageAdapter.Validate return an error, it is handed through
//to the caller of CreateStorageAdapter
func CreateStorageAdapter(name string, config Storage) (StorageAdapter, error) {
	var storage StorageAdapter
	switch config.StorageType {
	case KeepassTypeIdentifier:
		storage = &Keepass{
			Name:     name,
			FilePath: config.Config["path"],
		}
	case PassTypeIdentifier:
		storage = &Pass{
			Name:   name,
			Prefix: config.Config["prefix"],
		}
	default:
		return nil, errors.Newf("Unknown storage type %s", config.StorageType)
	}
	err, _ := storage.Validate()
	if err != nil {
		return nil, err
	}
	return storage, err
}

//GetStorageAdapterTypes returns a list of type identifiers for storage adapters
func GetStorageAdapterTypes() []string {
	return []string{
		KeepassTypeIdentifier,
		PassTypeIdentifier,
	}
}

//GetStorageAdapterDefaultConfig returns the default config for a given storage type
func GetStorageAdapterDefaultConfig(storageType string) (map[string]string, error) {
	var storage StorageAdapter
	switch storageType {
	case KeepassTypeIdentifier:
		storage = &Keepass{}
	case PassTypeIdentifier:
		storage = &Pass{}
	default:
		return nil, errors.Newf("Unknown storage type %s", storageType)
	}

	return storage.GetDefaultConfig(), nil
}
