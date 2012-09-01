package storage

type ConnectFn func(db, host string) DB

type DB interface {
	Get(set string) Set
	Close() error
	String() string
}

type Set interface {
	// Fills in the blank spaces of the supplied struct.
	Fill(data interface{}) error
	
	// Stores the supplied structs in the Set.
	Insert(data ...interface{}) error
	
	// Gets all results that satisfies the conditions.
	Get(res interface{}, where ...string) (int, error)
	
	// Gets the first result that satisfies the conditions.
	GetOne(res interface{}, where ...string) (err error)
	String() string
}

type StorageError struct {
	typ string
}

func NewStorageError(message string) *StorageError {
	return &StorageError{message}
}

func (e StorageError) Error() string {
	return e.typ
}
