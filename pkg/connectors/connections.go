package connectors

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/microlib/simple"
)

// Connections struct - all backend connections in a common object
type Connectors struct {
	S3Session *session.Session
	Http      *http.Client
	Logger    *simple.Logger
	Mode      string
}

// NewClientConnections - fucntion that creates all client connections and returns the interface
func NewClientConnections(logger *simple.Logger) Clients {
	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		logger.Error(fmt.Sprintf("NewClientConnections : %v", err))
		panic(err)
	}
	// svc := s3.New(sess)

	return &Connectors{S3Session: sess, Http: httpClient, Logger: logger}
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

// SetMode - simple push pull flag setting
func (c *Connectors) SetMode(mode string) {
	c.Mode = mode
}

// GetMode - simple flag check routine
func (c *Connectors) GetMode() string {
	return c.Mode
}

// ListObjectsV2 - wrapper for the s3 list service
func (c *Connectors) ListObjectsV2(in *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	svc := s3.New(c.S3Session)
	return svc.ListObjectsV2(in)
}

// GetObject - S3 Object download wrapper
func (c *Connectors) GetObject(opts *s3.GetObjectInput) ([]byte, error) {
	var b []byte
	svc := s3.New(c.S3Session)
	result, err := svc.GetObject(opts)
	if err != nil {
		// Message from an error.
		c.Error("Function GetObject %v", err)
		return b, err
	}

	b, err = ioutil.ReadAll(result.Body)
	if err != nil {
		c.Error("Function GetObject %v", err)
		return b, err
	}
	return b, nil
}

// PutObject - S3 Object uploader wrapper
func (c *Connectors) PutObject(opts *s3.PutObjectInput) (*string, error) {
	svc := s3.New(c.S3Session)
	// Body:   aws.ReadSeekCloser(strings.NewReader("HappyFace.jpg"))
	result, err := svc.PutObject(opts)
	if err != nil {
		// Message from an error.
		c.Error("Function PutObject %v", err)
		s := "error"
		return &s, err
	}
	return result.ETag, nil
}
