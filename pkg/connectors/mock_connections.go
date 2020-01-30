package connectors

import (
	"errors"
	"testing"

	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

type FakeCouchbase struct {
}

// Mock all connections
type MockConnections struct {
	Bucket *FakeCouchbase
}

func (r *MockConnections) Close() error {
	return nil
}

func (r *MockConnections) Upsert(collection string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error) {
	if collection == "" {
		return &gocb.MutationResult{}, errors.New("Empty collection value")
	}
	//r.Bucket
	return &gocb.MutationResult{}, nil
}

// NewTestConnections - create all mock connections
func NewTestConnections(file string, code int, logger *simple.Logger) Clients {

	bc := &FakeCouchbase{}
	conns := &MockConnections{Bucket: bc}
	return conns
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
