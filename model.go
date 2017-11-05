package csnotes

// Model implements a function for saving to the database, and a function for
// loading from it. They may use the helper methods Select, Sync, or Delete on
// a resource, as needed.
type Model interface {
	Save() error
	Load() error
}
