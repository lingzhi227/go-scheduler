package memory

// Backend is the storage interface for memory entries.
type Backend interface {
	Set(scope Scope, scopeID, key string, value any) error
	Get(scope Scope, scopeID, key string) (any, bool, error)
	Delete(scope Scope, scopeID, key string) error
	List(scope Scope, scopeID string) ([]string, error)
	SearchVector(scope Scope, scopeID string, query []float64, limit int) ([]VectorResult, error)
}
