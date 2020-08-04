package connectors

import (
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
)

// Client Interface - used as a receiver and can be overriden for testing
type Clients interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	Meta(string) string
	SetMode(string)
	GetMode() string
	Do(req *http.Request) (*http.Response, error)
	ListObjectsV2(in *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
	GetObject(in *s3.GetObjectInput) ([]byte, error)
	PutObject(in *s3.PutObjectInput) (*string, error)
}
