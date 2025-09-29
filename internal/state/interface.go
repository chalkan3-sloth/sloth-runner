package state

// StateManagerInterface defines the common interface for state management
type StateManagerInterface interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Delete(key string) error
	List(prefix string) ([]string, error)
	SetWithTTL(key string, value interface{}, ttlSeconds int) error
	GetMetadata(key string) (*StateMetadata, error)
	AcquireLock(key string, holder string, ttlSeconds int) error
	ReleaseLock(key string, holder string) error
	WithLock(key string, holder string, ttlSeconds int, fn func() error) error
	Close() error
}