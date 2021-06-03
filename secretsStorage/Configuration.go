package secretsStorage

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
	Storages map[string]Storage `yaml:"storages"`
	Profiles map[string]Profile `yaml:"profiles"`
}

type Storage struct {
	StorageType string            `yaml:"type"`
	Config      map[string]string `yaml:"config"`
}

type Profile struct {
	Storage   string            `yaml:"storage"`
	Path      string            `yaml:"path"`
	ConstEnv  map[string]string `yaml:"constEnv"`
	Env       map[string]string `yaml:"env"`
	DependsOn []string          `yaml:"dependsOn"`
}

func LoadConfigurationFromFile(path string) (*Configuration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Configuration

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

//Validate checks the validity of the profile. The storage and all profiles this
//profile depends on must be known to the registry.
func (p Profile) Validate() []string {
	var out []string
	registry := GetRegistry()
	if !registry.HasStorage(p.Storage) {
		out = append(out, fmt.Sprintf("references storage %s which is not defined", p.Storage))
	}
	for i := 0; i < len(p.DependsOn); i++ {
		if !registry.HasProfile(p.DependsOn[i]) {
			out = append(out, fmt.Sprintf("depends on %s which is not defined", p.DependsOn[i]))
		}
	}
	return out
}
