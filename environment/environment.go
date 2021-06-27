package environment

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

//Environment contains the current environment, variables which will be added to it and variables which will be removed
//from it the next time WriteStatements is called.
type Environment struct {
	//current holds the shell's environment after calling Load
	current map[string]string
	//addVars holds the variables to add with WriteStatements (will generate export key=value statements)
	addVars map[string]string
	//delVars holds the variable names to remove with WriteStatements (will generate unset value statements)
	delVars map[string]bool
}

//NewEnvironment creates a new Environment object and initializes the fields with empty maps / slices
func NewEnvironment() Environment {
	return Environment{
		current: map[string]string{},
		addVars: map[string]string{},
		delVars: map[string]bool{},
	}
}

func (e *Environment) Load() {
	mapped := os.Environ()
	for _, element := range mapped {
		parts := strings.Split(element, "=")
		e.current[parts[0]] = strings.Join(parts[1:], "=")
	}
}

//Retrieves a currently set environment variable by the given key. If it does not
//exist, the defaultValue is returned.
func (e *Environment) GetCurrent(key string, defaultValue string) string {
	value, exists := e.current[key]
	if !exists {
		return defaultValue
	}
	return value
}

//Set adds an environment variable with given key and value to the list of variables to set. Call WriteStatements to
//create export statements consumable by a shell.
func (e *Environment) Set(key string, value string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	e.addVars[key] = value
	return nil
}

//Unset removes an already set environment variable. Also removes it from the list of variables which will be set.
func (e *Environment) Unset(key string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	delete(e.addVars, key)
	e.delVars[key] = true
	return nil
}

//WriteStatements writes a list of export and unset statements to update the environment of a shell.
func (e *Environment) WriteStatements() string {
	var output []string
	for key, value := range e.addVars {
		output = append(output, fmt.Sprintf("export %s=\"%s\"", key, value))
	}
	for key, _ := range e.delVars {
		output = append(output, fmt.Sprintf("unset %s", key))
	}
	return strings.Join(output, ";")
}
