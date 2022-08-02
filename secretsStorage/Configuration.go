package secretsStorage

import (
	"gopkg.in/errgo.v2/fmt/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Options          Options             `yaml:"options,omitempty"`
	Storages         map[string]Storage  `yaml:"storages"`
	Profiles         map[string]Profile  `yaml:"profiles"`
	DirectoryMapping map[string][]string `yaml:"directoryMapping"`
}

type Storage struct {
	StorageType string            `yaml:"type"`
	Config      map[string]string `yaml:"config"`
}

//Options controls the general behavior of envManager. It is only read from the config file in the home directory and
//cannot be overridden by other config files.
type Options struct {
	//DisableCollisionDetection allows overwriting profiles, storages and mappings instead of returning an error
	DisableCollisionDetection bool `yaml:"disableCollisionDetection"`
}

//NewConfiguration creates a new, empty configuration object
func NewConfiguration() Configuration {
	return Configuration{
		Options:          Options{},
		Storages:         map[string]Storage{},
		Profiles:         map[string]Profile{},
		DirectoryMapping: map[string][]string{},
	}
}

//LoadFromFile loads the config file at the given path. It will merge the current config with the config from the file.
func (c *Configuration) LoadFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	return nil
}

//WriteToFile writes the current config to given path. It will not overwrite an existing file
//except when replace is set to true.
func (c Configuration) WriteToFile(path string, replace bool) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	fileInfo, err := os.Stat(path)
	if fileInfo != nil && !replace {
		return errors.Newf("Will not overwrite %s without being explicitly told to do so.", path)
	}
	err = ioutil.WriteFile(path, data, 0600)

	return err
}
