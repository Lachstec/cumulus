package db

// Predicate represents a search constraint
// in order to filter rows returned by the database.
// The Predicate should return true if a value T should
// be included, else false.
type Predicate[T any] func(T) bool

// Store represents types that can be used to persist
// an associated entity to an underlying data store.
type Store[T any] interface {
	// GetById searches the data store for a record with the given id.
	// If a record was not found, the return value will be nil.
	GetById(id int64) (*T, error)

	// Find returns all records in the store where predicate is true.
	Find(predicate Predicate[*T]) ([]*T, error)

	// Add persists a given type T to the underlying data store.
	Add(*T) (int64, error)

	// Update changes the stored record if it exists, else it gets inserted.
	Update(*T) (*T, error)

	// Delete removes the given record from the data store.
	Delete(*T) error
}
