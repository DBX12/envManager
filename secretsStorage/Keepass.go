package secretsStorage

import (
	"fmt"
	"github.com/tobischo/gokeepasslib/v3"
	"gopkg.in/errgo.v2/fmt/errors"
	"os"
	"strings"
)

const KeepassTypeIdentifier = "keepass"

type Keepass struct {
	Name     string
	FilePath string
	database *gokeepasslib.Database
}

func (k Keepass) IsCaseSensitive() bool {
	return true
}

func (k Keepass) Validate() (error, []string) {
	var out []string
	validationFailed := false

	fileExists := true
	_, err := os.Stat(k.FilePath)
	if err != nil {
		fileExists = false
		validationFailed = true
	}
	out = append(out, fmt.Sprintf("Configured file is %s\nFile exists: %t", k.FilePath, fileExists))

	if validationFailed {
		return errors.Newf("Validation of %s failed. Run debug storage %s to check it in detail", k.Name, k.Name), out
	}
	return nil, out
}

func (k *Keepass) promptCredentials() string {
	fmt.Printf("Enter password for %s\n> ", k.Name)
	var password string
	_, _ = fmt.Scanln(&password)
	return password
}

func (k *Keepass) GetEntry(key string) (*Entry, error) {
	if k.database == nil {
		_ = k.openDatabase()
	}
	currentGroup := &k.database.Content.Root.Groups[0]
	parts := strings.Split(key, "/")
	lastIndex := len(parts) - 1
	entryName := parts[lastIndex]
	parts = parts[:lastIndex] // all but the last part
	for _, part := range parts {
		var err error
		currentGroup, err = findGroup(currentGroup, part)
		if err != nil {
			return nil, err
		}
	}
	kpEntry, err := findEntry(currentGroup, entryName)
	if err != nil {
		return nil, err
	}
	entry, err := toEntry(kpEntry)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func toEntry(kpEntry *gokeepasslib.Entry) (*Entry, error) {
	entry := NewEntry()
	for _, valueData := range kpEntry.Values {
		err := entry.SetAttribute(valueData.Key, valueData.Value.Content)
		if err != nil {
			return nil, err
		}
	}
	return &entry, nil
}

func findEntry(group *gokeepasslib.Group, name string) (*gokeepasslib.Entry, error) {
	for _, entry := range group.Entries {
		if entry.GetTitle() == name {
			return &entry, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Could not find entry with name %s in group %s", name, group.Name))
}

func findGroup(group *gokeepasslib.Group, name string) (*gokeepasslib.Group, error) {
	for _, subGroup := range group.Groups {
		if subGroup.Name == name {
			return &subGroup, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Could not find subgroup with name %s in group %s", name, group.Name))
}

func (k *Keepass) openDatabase() error {
	fileHandle, err := os.Open(k.FilePath)
	if err != nil {
		return err
	}
	k.database = gokeepasslib.NewDatabase()
	k.database.Credentials = gokeepasslib.NewPasswordCredentials(k.promptCredentials())
	_ = gokeepasslib.NewDecoder(fileHandle).Decode(k.database)
	_ = k.database.UnlockProtectedEntries()
	return nil
}
