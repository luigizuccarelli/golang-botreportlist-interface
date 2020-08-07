package connectors

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/microlib/simple"
)

var count int = 0

type FakeS3 struct {
}

// Mock all connections
type MockConnectors struct {
	S3Service *FakeS3
	Http      *http.Client
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

// SetMode - simple push pull flag setting
func (c *MockConnectors) SetMode(mode string) {
	c.Mode = mode
}

// GetMode - simple flag check routine
func (c *MockConnectors) GetMode() string {
	return c.Mode
}

// Do - log
func (c *MockConnectors) Do(req *http.Request) (*http.Response, error) {
	if c.Flag == "true" {
		return nil, errors.New("forced http error")
	}
	return c.Http.Do(req)
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
	if count >= 2 {
		truncated = false
		count = 0
	}
	count++
	s := &s3.ListObjectsV2Output{Contents: objs, IsTruncated: &truncated}
	return s, nil
}

// GetObject - S3 Object download wrapper
func (c *MockConnectors) GetObject(opts *s3.GetObjectInput) ([]byte, error) {
	var b []byte
	if c.Flag == "true" {
		return b, errors.New("forced s3 ListObjectsV2 error")
	}
	b = []byte("GetObject -> this is working!!!")
	return b, nil
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

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// NewTestConnector - creates all test connectors
func NewTestConnectors(file string, code int, logger *simple.Logger) Clients {

	// we first load the json payload to simulate a call to middleware
	// for now just ignore failures.
	data, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Error(fmt.Sprintf("file data %v\n", err))
		panic(err)
	}
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString(string(data))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	conns := &MockConnectors{S3Service: &FakeS3{}, Http: httpclient, Logger: logger, Flag: "false"}
	return conns
}
