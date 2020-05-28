package connectors

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/microlib/simple"
)

// Mock all connections
type MockConnectors struct {
	Http   *http.Client
	Logger *simple.Logger
	Flag   string
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

// Do - log wrapper
func (c *MockConnectors) Do(req *http.Request) (*http.Response, error) {
	if c.Flag == "true" {
		return nil, errors.New("forced http error")
	}
	return c.Http.Do(req)
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

	conns := &MockConnectors{Http: httpclient, Logger: logger, Flag: "false"}
	return conns
}
