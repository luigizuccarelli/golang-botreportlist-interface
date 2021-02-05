package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/connectors"
	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/schema"
	"github.com/aws/aws-sdk-go/service/s3"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

// Fake connectors used in this package for testing only
// Used by handlers_test.go

var count int = 0

type FakeS3 struct {
}

// Fake all connections
type FakeConnectors struct {
	S3Service *FakeS3
	Logger    *simple.Logger
	Flag      string
	Mode      string
}

// Error - log wrapper
func (c *FakeConnectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

// Info - log wrapper
func (c *FakeConnectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

// Debug - log wrapper
func (c *FakeConnectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

// Trace - log wrapper
func (c *FakeConnectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

// Meta - log wrapper
func (c *FakeConnectors) Meta(flag string) string {
	c.Flag = flag
	return flag
}

// Upsert : wrapper function for couchbase update
func (c *FakeConnectors) Upsert(uuid string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error) {
	if c.Flag == "true" {
		return &gocb.MutationResult{}, errors.New("Upsert (forced error)")
	}
	return &gocb.MutationResult{}, nil
}

// GetList - Couchbase list wrapper
func (c *FakeConnectors) GetList(offset string, limit string) ([]schema.ReportList, error) {
	var list []schema.ReportList
	if c.Flag == "true" {
		return list, errors.New("forced GetList (DB) error")
	}
	b, _ := ioutil.ReadFile("../../tests/payload-reportlist-01.json")
	c.Trace("GetList mock response %s", string(b))
	json.Unmarshal(b, &list)
	return list, nil
}

// GetConfusionMatrix - Couchbase stats wrapper
func (c *FakeConnectors) GetConfusionMatrix() (*schema.ConfusionMatrix, error) {
	var stats *schema.ConfusionMatrix
	if c.Flag == "true" {
		return stats, errors.New("forced GetAllStats (DB) error")
	}
	b, _ := ioutil.ReadFile("../../tests/confusion-matrix.json")
	c.Trace("GetAllStats mock response %s", string(b))
	json.Unmarshal(b, &stats)
	return stats, nil
}

// GetListCount - Couchbase list count wrapper
func (c *FakeConnectors) GetListCount() (*int64, error) {
	if c.Flag == "true" {
		val := int64(0)
		return &val, errors.New("forced GetListCount (DB) error")
	}
	c.Trace("GetListCount 1234")
	val := int64(1234)
	return &val, nil
}

// GetObject - S3 Object download wrapper
func (c *FakeConnectors) GetObject(opts *s3.GetObjectInput) (*schema.ReportContent, error) {
	var rc *schema.ReportContent
	if c.Flag == "true" {
		return rc, errors.New("forced s3 ListObjectsV2 error")
	}
	b, _ := ioutil.ReadFile("../../tests/report-payload.json")
	json.Unmarshal(b, &rc)
	return rc, nil
}

// NewTestConnector - creates all test connectors
func NewTestConnectors(code int, logger *simple.Logger) connectors.Clients {
	conns := &FakeConnectors{Logger: logger, Flag: "false"}
	return conns
}
