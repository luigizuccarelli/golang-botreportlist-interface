package connectors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/schema"
	"github.com/aws/aws-sdk-go/service/s3"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

var count int = 0

type FakeS3 struct {
}

// Mock all connections
type MockConnectors struct {
	S3Service *FakeS3
	Logger    *simple.Logger
	Flag      string
	Mode      string
}

// Error - log wrapper
func (c *MockConnectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

// Info - log wrapper
func (c *MockConnectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

// Debug - log wrapper
func (c *MockConnectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

// Trace - log wrapper
func (c *MockConnectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

// Meta - log wrapper
func (c *MockConnectors) Meta(flag string) string {
	c.Flag = flag
	return flag
}

// Upsert : wrapper function for couchbase update
func (c *MockConnectors) Upsert(uuid string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error) {
	if c.Flag == "true" {
		return &gocb.MutationResult{}, errors.New("Upsert (forced error)")
	}
	return &gocb.MutationResult{}, nil
}

// ListObjectsV2 - S3 wrapper
func (c *MockConnectors) ListObjectsV2(in *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	var objs []*s3.Object
	var truncated bool = true

	if c.Flag == "true" {
		return nil, errors.New("forced s3 ListObjectsV2 error")
	}
	name := "test"
	last := time.Now()
	size := int64(3232)
	sc := "TEST"
	objs = append(objs, &s3.Object{Key: &name, LastModified: &last, Size: &size, StorageClass: &sc})

	newname := "nextest"
	newlast := time.Now()
	newsize := int64(6464)
	newsc := "NEXTTEST"
	objs = append(objs, &s3.Object{Key: &newname, LastModified: &newlast, Size: &newsize, StorageClass: &newsc})
	if count >= 1 {
		truncated = false
		count = 0
	}
	count++
	s := &s3.ListObjectsV2Output{Contents: objs, IsTruncated: &truncated}
	return s, nil
}

// GetList - Couchbase list wrapper
func (c *MockConnectors) GetList(offset string, limit string) ([]schema.ReportList, error) {
	var list []schema.ReportList
	if c.Flag == "true" {
		return list, errors.New("forced GetList (DB) error")
	}
	b, _ := ioutil.ReadFile("../../tests/payload-reportlist-01.json")
	c.Trace("GetList mock response %s", string(b))
	json.Unmarshal(b, &list)
	return list, nil
}

// GetAllStats - Couchbase stats wrapper
func (c *MockConnectors) GetAllStats() ([]schema.Stat, error) {
	var stats []schema.Stat
	if c.Flag == "true" {
		return stats, errors.New("forced GetAllStats (DB) error")
	}
	b, _ := ioutil.ReadFile("../../tests/payload-stats.json")
	c.Trace("GetAllStats mock response %s", string(b))
	json.Unmarshal(b, &stats)
	return stats, nil
}

// GetObject - S3 Object download wrapper
func (c *MockConnectors) GetObject(opts *s3.GetObjectInput) (*schema.ReportContent, error) {
	var rc *schema.ReportContent
	if c.Flag == "true" {
		return rc, errors.New("forced s3 ListObjectsV2 error")
	}
	b, _ := ioutil.ReadFile("../../tests/report-payload.json")
	json.Unmarshal(b, &rc)
	return rc, nil
}

// PutObject - S3 Object uploader wrapper
func (c *MockConnectors) PutObject(opts *s3.PutObjectInput) (*string, error) {
	if c.Flag == "true" {
		s := "error"
		return &s, errors.New("forced s3 ListObjectsV2 error")
	}
	s := "PutObject This is working !!!"
	return &s, nil
}

// NewTestConnector - creates all test connectors
func NewTestConnectors(code int, logger *simple.Logger) Clients {
	conns := &MockConnectors{Logger: logger, Flag: "false"}
	return conns
}
