package connectors

import (
	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/schema"
	"github.com/aws/aws-sdk-go/service/s3"
	gocb "github.com/couchbase/gocb/v2"
)

// Client Interface - used as a receiver and can be overriden for testing
type Clients interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	Meta(force string) string
	GetConfusionMatrix() (*schema.ConfusionMatrix, error)
	GetList(string, string) ([]schema.ReportList, error)
	GetListCount() (*int64, error)
	Upsert(uuid string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error)
	GetObject(in *s3.GetObjectInput) (*schema.ReportContent, error)
}
