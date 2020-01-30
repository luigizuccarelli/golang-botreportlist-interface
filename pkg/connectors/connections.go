package connectors

import (
	"os"

	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

// Connections struct - all backend connections in a common object
type Connections struct {
	Bucket  *gocb.Bucket
	Cluster *gocb.Cluster
}

// Upsert call implementation
func (r *Connections) Upsert(col string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error) {
	//collection := bucket.Scope("default").Collection(col, &gocb.CollectionOptions{})
	collection := r.Bucket.DefaultCollection()
	return collection.Upsert(col, value, opts)
}

func (r *Connections) Close() error {
	return r.Cluster.Close(&gocb.ClusterCloseOptions{})
}

// NewClientConnectors returns Connectors struct
func NewClientConnections(logger *simple.Logger) Clients {

	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: os.Getenv("COUCHBASE_USER"),
			Password: os.Getenv("COUCHBASE_PASSWORD"),
		},
	}
	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_HOST"), opts)
	if err != nil {
		panic(err)
	}

	// get a bucket reference
	bucket := cluster.Bucket(os.Getenv("COUCHBASE_BUCKET"), &gocb.BucketOptions{})

	conns := &Connections{Bucket: bucket, Cluster: cluster}
	return conns
}
