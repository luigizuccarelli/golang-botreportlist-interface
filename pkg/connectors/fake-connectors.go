// +build fake

package connectors

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/schema"
	"github.com/aws/aws-sdk-go/service/s3"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

// Fake connectors for testing in this package
// Used by connections_test.go

var count int = 0

// Connectors - overrides the real implemntation (using gocb.* dependencies)
// The file directive +build mock ensures its use (see the first line of this file)_
type Connectors struct {
	Bucket    *FakeBucket
	Cluster   *FakeCluster
	S3Service *FakeS3
	Logger    *simple.Logger
	Flag      string
}

// FakeCluster
type FakeCluster struct {
	Force string
}

// FakeBucket
type FakeBucket struct {
	Force string
}

// FakeQuery
type FakeQuery struct {
}

// FakeCollection
type FakeCollection struct {
	Force string
}

// FakeResult
type FakeResult struct {
	Force string
}

// Fake AWS S3 session
type FakeS3 struct {
	Force string
}

func (fs3 *FakeS3) GetObject(opts *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	var obj *s3.GetObjectOutput
	var data []byte

	if fs3.Force == "true" {
		return obj, errors.New("Function GetObject forced error")
	}
	if fs3.Force == "data" {
		data = []byte("{ test")
	} else {
		data, _ = ioutil.ReadFile("../../tests/report-payload.json")
	}
	obj = &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewBuffer(data))}
	return obj, nil
}

// Couchbase fake

// Query - inject our implementation for testing
func (fc *FakeCluster) Query(query string, opts *gocb.QueryOptions) (*FakeResult, error) {
	if fc.Force == "true" {
		return &FakeResult{}, errors.New("Function Query forced error")
	}
	return &FakeResult{Force: fc.Force}, nil
}

// DefaultCollection - override the original gocb implementation
func (fb *FakeBucket) DefaultCollection() *FakeCollection {
	return &FakeCollection{Force: fb.Force}
}

// Next - override the original golang implementation
func (fr *FakeResult) Next() bool {
	if count < 2 {
		count++
		return true
	}
	count = 0
	return false
}

// Row - override the original golang implementation
func (fr *FakeResult) Row(ptr interface{}) error {
	if fr.Force == "true" {
		return errors.New("Function Row forced error")
	}
	if reflect.TypeOf(ptr).String() == "**schema.Stat" {
		var stat *schema.Stat
		data := `{
			"count": 60,
			"success": true
		}`
		json.Unmarshal([]byte(data), &stat)
		*ptr.(**schema.Stat) = stat
	} else {
		var rl *schema.ReportList
		data := `{ "id": "096esbpfrk8b3nhdlfhditsmk10gj03g06i3c201.json",
    "servisbotstats": {
      "EmailClassification": "Cancel",
      "ProcessOutcome": "No Action",
      "UserClassification": "",
      "success": false
    }}`
		json.Unmarshal([]byte(data), &rl)
		*ptr.(**schema.ReportList) = rl
	}
	return nil
}

// One - override the original golang implementation
func (fr *FakeResult) One(ptr interface{}) error {
	var count = make(map[string]int64)
	if fr.Force == "true" {
		return errors.New("Function One forced error")
	}
	count["count"] = int64(1234)
	*ptr.(*map[string]int64) = count
	return nil
}

// Err - fake the error handler
func (fr *FakeResult) Err() error {
	return nil
}

// Close - do we need to explain ?
func (fr *FakeResult) Close() {
}

// Upsert : wrapper function for couchbase collection upsert
func (fc *FakeCollection) Upsert(col string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error) {
	if fc.Force == "error" {
		return &gocb.MutationResult{}, errors.New("Forced collection upsert error")
	}
	return &gocb.MutationResult{}, nil
}
