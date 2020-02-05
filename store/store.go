package store

// Universal interface that btree must implement
// Only support string as key and map[string]string as value type
// Store type is used by the api package - btree and lsm is never accessed direct

type Value map[string]string

type Store interface {
	Insert(string, Value) error
	Update(string, Value) error
	Search(string) Value
	Remove(string) error
}
