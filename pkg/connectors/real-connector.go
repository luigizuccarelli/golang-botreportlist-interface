// +build real

package connectors

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

// Real "runtime" client connections

// Connections struct - all backend connections in a common object
type Connectors struct {
	S3Service *s3.S3
	Bucket    *gocb.Bucket
	Cluster   *gocb.Cluster
	Logger    *simple.Logger
	Mode      string
}

// NewClientConnections - fucntion that creates all client connections and returns the interface
func NewClientConnections(logger *simple.Logger) Clients {
	// setup aws session
	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		logger.Error(fmt.Sprintf("NewClientConnections : %v", err))
		panic(err)
	}

	svc := s3.New(sess)

	opts := gocb.ClusterOptions{
		Username: os.Getenv("COUCHBASE_USER"),
		Password: os.Getenv("COUCHBASE_PASSWORD"),
	}
	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_HOST"), opts)
	if err != nil {
		logger.Error(fmt.Sprintf("Couchbase connection: %v", err))
		panic(err)
	}

	// get a bucket reference
	// bucket := cluster.Bucket(os.Getenv("COUCHBASE_BUCKET"), &gocb.BucketOptions{}) v.2.0.0-beta-1
	bucket := cluster.Bucket(os.Getenv("COUCHBASE_BUCKET"))
	logger.Info(fmt.Sprintf("Couchbase connection: %v", bucket))

	return &Connectors{Bucket: bucket, Cluster: cluster, S3Service: svc, Logger: logger}
}
