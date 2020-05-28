package connectors

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/microlib/simple"
)

// Connections struct - all backend connections in a common object
type Connectors struct {
	Http   *http.Client
	Logger *simple.Logger
}

// NewClientConnections - fucntion that creates all client connections and returns the interface
func NewClientConnections(logger *simple.Logger) Clients {
	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	return &Connectors{Http: httpClient, Logger: logger}
}

// Error - log wrapper
func (c *Connectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

// Info - log wrapper
func (c *Connectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

// Debug - log wrapper
func (c *Connectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

// Trace - log wrapper
func (c *Connectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

// Meta - test wrapper
func (c *Connectors) Meta(info string) string {
	return info
}

// Do - http wrapper
func (c *Connectors) Do(req *http.Request) (*http.Response, error) {
	return c.Http.Do(req)
}
