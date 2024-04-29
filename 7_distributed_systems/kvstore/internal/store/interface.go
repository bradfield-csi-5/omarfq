package store

type KVStore interface {
	// Returns the value by a given key or an error
	// if the key does not exist
	Get(key string) (string, error)
	// Sets or updates a new key-value pair in the
	// data store
	Set(key, value string) error
}
