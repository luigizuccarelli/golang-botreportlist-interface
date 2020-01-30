package connectors

import (
	gocb "github.com/couchbase/gocb/v2"
)

// Client Interface - used as a receiver and can be overriden for testing
type Clients interface {
	Upsert(collection string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error)
	Close() error
}
