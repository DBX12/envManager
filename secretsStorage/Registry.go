package secretsStorage

import (
	"gopkg.in/errgo.v2/fmt/errors"
	"sync"
)

type Registry struct {
	storages map[string]StorageAdapter
	profiles map[string]Profile
}

var instance *Registry
var lock = &sync.Mutex{}

func GetRegistry() *Registry {
	// skip expensive locking if it already exists
	if instance == nil {
		// acquire a lock to ensure only one goroutine can make an instance
		lock.Lock()
		defer lock.Unlock()
		// after acquiring the lock, check again if another goroutine got here
		// first and made an instance
		if instance == nil {
			instance = newRegistry()
		}
	}
	return instance
}

func newRegistry() *Registry {
	return &Registry{
		storages: map[string]StorageAdapter{},
		profiles: map[string]Profile{},
	}
}

func (r *Registry) AddStorage(name string, storage StorageAdapter) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	r.storages[name] = storage
	return nil
}

func (r *Registry) AddProfile(name string, profile Profile) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	r.profiles[name] = profile
	return nil
}

//GetProfile retrieves a profile with given name. Will return an error if given
//name is empty or unknown to the registry
func (r Registry) GetProfile(name string) (*Profile, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	profile, exists := r.profiles[name]
	if !exists {
		return nil, errors.Newf("profile with name %s does not exist", name)
	}
	return &profile, nil
}

//GetStorage retrieves a storage instance with given name. Will return an error
//if given name is empty or unknown to the registry
func (r Registry) GetStorage(name string) (*StorageAdapter, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	storage, exists := r.storages[name]
	if !exists {
		return nil, errors.Newf("storage with name %s does not exist", name)
	}
	return &storage, nil
}

//HasStorage checks if the registry knows about a storage with this name
func (r Registry) HasStorage(name string) bool {
	_, exists := r.storages[name]
	return exists
}

//HasProfile checks if the registry knows about a profile with this name
func (r Registry) HasProfile(name string) bool {
	_, exists := r.profiles[name]
	return exists
}

//GetAllStorages returns all storages known to the registry
func (r Registry) GetAllStorages() map[string]StorageAdapter {
	return r.storages
}
