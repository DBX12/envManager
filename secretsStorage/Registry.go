package secretsStorage

import (
	"gopkg.in/errgo.v2/fmt/errors"
	"sync"
)

type Registry struct {
	storages         map[string]StorageAdapter
	profiles         map[string]Profile
	directoryMapping map[string][]string
}

var instance *Registry
var once sync.Once

func GetRegistry() *Registry {
	once.Do(func() {
		instance = newRegistry()
	})
	return instance
}

func newRegistry() *Registry {
	return &Registry{
		storages:         map[string]StorageAdapter{},
		profiles:         map[string]Profile{},
		directoryMapping: map[string][]string{},
	}
}

// AddStorage adds a storage adapter to the registry. If the given name already exists, the old storage adapter will be
// replaced. Will return an error if the storage name is empty or the storage adapter is nil.
func (r *Registry) AddStorage(name string, storage StorageAdapter) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if storage == nil {
		return errors.New("storage cannot be nil")
	}
	r.storages[name] = storage
	return nil
}

// AddProfile adds a profile to the registry. If the given name already exists, the old profile instance will be
// replaced. Will return an error if the profile name is empty.
func (r *Registry) AddProfile(name string, profile Profile) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	profile.SetName(name)
	r.profiles[name] = profile
	return nil
}

// AddDirectoryMapping adds a directory mapping to the registry. If the given path
// already exists, the mapping will be replaced. Will return an error if the path
// or profiles are empty
func (r *Registry) AddDirectoryMapping(path string, profiles []string) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}
	if len(profiles) == 0 {
		return errors.New("profiles cannot be empty")
	}
	r.directoryMapping[path] = profiles
	return nil
}

// GetProfile retrieves a profile with given name. Will return an error if given
// name is empty or unknown to the registry
func (r *Registry) GetProfile(name string) (*Profile, error) {
	if name == "" {
		return nil, errors.New("profile name cannot be empty")
	}
	profile, exists := r.profiles[name]
	if !exists {
		return nil, errors.Newf("profile with name %s does not exist", name)
	}
	return &profile, nil
}

// GetStorage retrieves a storage instance with given name. Will return an error
// if given name is empty or unknown to the registry
func (r *Registry) GetStorage(name string) (*StorageAdapter, error) {
	if name == "" {
		return nil, errors.New("storage name cannot be empty")
	}
	storage, exists := r.storages[name]
	if !exists {
		return nil, errors.Newf("storage with name %s does not exist", name)
	}
	return &storage, nil
}

// GetDirectoryMapping retrieves the profile names mapped to the given path.
// Will return an error if given path is empty or unknown to the registry
func (r *Registry) GetDirectoryMapping(path string) ([]string, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}
	profiles, exists := r.directoryMapping[path]
	if !exists {
		return nil, errors.Newf("directory mapping for path %s does not exist", path)
	}
	return profiles, nil
}

// HasStorage checks if the registry knows about a storage with this name
func (r *Registry) HasStorage(name string) bool {
	_, exists := r.storages[name]
	return exists
}

// HasProfile checks if the registry knows about a profile with this name
func (r *Registry) HasProfile(name string) bool {
	_, exists := r.profiles[name]
	return exists
}

// HasDirectoryMapping checks if the registry knows about a directory mapping for this path
func (r *Registry) HasDirectoryMapping(path string) bool {
	_, exists := r.directoryMapping[path]
	return exists
}

// GetAllStorages returns all storages known to the registry
func (r *Registry) GetAllStorages() map[string]StorageAdapter {
	return r.storages
}

// GetAllProfiles returns all profiles known to the registry
func (r *Registry) GetAllProfiles() map[string]Profile {
	return r.profiles
}

// GetAllDirectoryMappings returns all directory mappings known to the registry
func (r *Registry) GetAllDirectoryMappings() map[string][]string {
	return r.directoryMapping
}

// GetStorageNames returns the names of all storages known to the registry
func (r *Registry) GetStorageNames() []string {
	var out []string
	for name := range r.storages {
		out = append(out, name)
	}
	return out
}

// GetProfileNames returns all profiles names known to the registry
func (r *Registry) GetProfileNames() []string {
	var out []string
	for name := range r.profiles {
		out = append(out, name)
	}
	return out
}

// GetDirectoryMappedPaths returns all paths which have a mapping in the registry
func (r *Registry) GetDirectoryMappedPaths() []string {
	var out []string
	for name := range r.directoryMapping {
		out = append(out, name)
	}
	return out
}
