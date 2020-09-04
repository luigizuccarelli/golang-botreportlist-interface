package connectors

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/schema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

// Connections struct - all backend connections in a common object
type Connectors struct {
	S3Session *session.Session
	Bucket    *gocb.Bucket
	Cluster   *gocb.Cluster
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
	opts := gocb.ClusterOptions{
		Username: os.Getenv("COUCHBASE_USER"),
		Password: os.Getenv("COUCHBASE_PASSWORD"),
	}
	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_HOST"), opts)
	if err != nil {
		logger.Error(fmt.Sprintf("Couchbase connection: %v", err))
		panic(err)
	}

	// get a bucket reference
	// bucket := cluster.Bucket(os.Getenv("COUCHBASE_BUCKET"), &gocb.BucketOptions{}) v.2.0.0-beta-1
	bucket := cluster.Bucket(os.Getenv("COUCHBASE_BUCKET"))
	logger.Info(fmt.Sprintf("Couchbase connection: %v", bucket))

	return &Connectors{Bucket: bucket, Cluster: cluster, S3Session: sess, Http: httpClient, Logger: logger}
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

// GetObject - S3 Object download wrapper
func (c *Connectors) GetObject(opts *s3.GetObjectInput) (*schema.ReportContent, error) {
	var rc *schema.ReportContent
	svc := s3.New(c.S3Session)
	result, err := svc.GetObject(opts)
	if err != nil {
		// Message from an error.
		c.Error("Function GetObject %v", err)
		return rc, err
	}

	b, err := ioutil.ReadAll(result.Body)
	if err != nil {
		c.Error("Function GetObject %v", err)
		return rc, err
	}
	err = json.Unmarshal(b, &rc)
	if err != nil {
		c.Error("Function GetObject %v", err)
		return rc, err
	}

	return rc, nil
}

// Upsert : wrapper function for couchbase update
func (c *Connectors) Upsert(uuid string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error) {
	collection := c.Bucket.DefaultCollection()
	return collection.Upsert(uuid, value, opts)
}

// GetList - get all reports list
func (c *Connectors) GetList(offset string, limit string) ([]schema.ReportList, error) {
	var stats []schema.ReportList
	var stat *schema.ReportList

	query := "select meta().id as id,* from servisbotstats offset " + offset + " limit " + limit
	c.Trace("Function GetList %s", query)
	res, err := c.Cluster.Query(query, &gocb.QueryOptions{})
	if err != nil {
		return stats, err
	}

	// iterate through each object
	for res.Next() {
		err := res.Row(&stat)
		if err != nil {
			break
		}
		stats = append(stats, *stat)
	}

	// always check for errors after iterating
	err = res.Err()
	if err != nil {
		return stats, err
	}
	return stats, nil
}

// GetListCount - get total of reports in related to s3 report bucket from couchbase
func (c *Connectors) GetListCount() (int64, error) {
	var count int64

	query := "select count(meta().id) as count from servisbotstats"
	c.Trace("Function GetListCount %s", query)
	res, err := c.Cluster.Query(query, &gocb.QueryOptions{})
	if err != nil {
		return count, err
	}

	err = res.One(&count)
	if err != nil {
		return count, err
	}

	// always check for errors after iterating
	err = res.Err()
	if err != nil {
		return count, err
	}
	return count, nil
}

// GetAllStats - get stats for bot accuracy
func (c *Connectors) GetAllStats() ([]schema.Stat, error) {
	var stats []schema.Stat
	var stat *schema.Stat

	query := "select count(meta().id) as count,`success` from servisbotstats group by `success`"
	c.Trace("Function GetAllStats %s", query)
	res, err := c.Cluster.Query(query, &gocb.QueryOptions{})
	if err != nil {
		return stats, err
	}

	// iterate through each object
	for res.Next() {
		err := res.Row(&stat)
		if err != nil {
			break
		}
		stats = append(stats, *stat)
	}

	// always check for errors after iterating
	err = res.Err()
	if err != nil {
		return stats, err
	}
	return stats, nil
}
