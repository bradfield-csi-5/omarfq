package common

type Storer interface {
	// Returns the value by a given key or an error
	// if the key does not exist
	Write(key string) (string, error)
	// Sets or updates a new key-value pair in the
	// data store
	Read(key, value string) error
}
