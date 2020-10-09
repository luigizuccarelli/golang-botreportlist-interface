package connectors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/schema"
	"github.com/aws/aws-sdk-go/service/s3"
	gocb "github.com/couchbase/gocb/v2"
)

// Using the +build directive we can plugin (via the receiver) fake or real connectors

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

// Meta - used for testing ignored in real implementation
func (c *Connectors) Meta(force string) string {
	return force
}

// GetObject - S3 Object download wrapper
func (c *Connectors) GetObject(opts *s3.GetObjectInput) (*schema.ReportContent, error) {
	var rc *schema.ReportContent
	result, err := c.S3Service.GetObject(opts)
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

	query := "select meta().id as id,* from servisbotstats order by `servisbotstats`.`Timestamp` desc offset " + offset + " limit " + limit
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
func (c *Connectors) GetListCount() (*int64, error) {
	var count map[string]int64
	query := "select count(meta().id) as count from servisbotstats"
	c.Trace("Function GetListCount %s", query)
	res, err := c.Cluster.Query(query, &gocb.QueryOptions{})
	if err != nil {
		v := int64(0)
		return &v, err
	}
	err = res.One(&count)
	c.Trace("Function GetListCount count %v", count)
	if err != nil {
		c.Trace("Function GetListCount %v", err)
		v := int64(0)
		return &v, err
	}

	// always check for errors after iterating
	err = res.Err()
	if err != nil {
		v := int64(0)
		return &v, err
	}
	v := count["count"]
	return &v, nil
}

// GetAllStats - get stats for bot accuracy
func (c *Connectors) GetAllStats() ([]schema.Stat, error) {
	var stats []schema.Stat
	var stat *schema.Stat

	query := "select count(meta().id) as count,`Success` from servisbotstats group by `Success`"
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
