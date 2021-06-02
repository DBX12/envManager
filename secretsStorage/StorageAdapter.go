package secretsStorage

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
}
