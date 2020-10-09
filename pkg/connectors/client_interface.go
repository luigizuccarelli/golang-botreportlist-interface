package connectors

import (
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/schema"
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
	GetAllStats() ([]schema.Stat, error)
	GetList(string, string) ([]schema.ReportList, error)
	GetListCount() (*int64, error)
	Upsert(uuid string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error)
	GetObject(in *s3.GetObjectInput) (*schema.ReportContent, error)
}
