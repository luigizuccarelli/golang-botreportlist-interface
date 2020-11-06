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

// GetConfusiorMatrix - get confusion matrix stats for bot accuracy
func (c *Connectors) GetConfusionMatrix() (*schema.ConfusionMatrix, error) {
	var cm = &schema.ConfusionMatrix{}

	// we first get all the counts for the bot
	query := "select distinct ProcessOutcome,'UserClassification' as UserCalssififcation,count(ProcessOutcome) as count from servisbotstats where  UserClassification != \"\" group by ProcessOutcome"
	statsTotal, err := getStatsData(query, c)
	if err != nil {
		return cm, err
	}

	query = "select count(ProcessOutcome) as count,ProcessOutcome,UserClassification from servisbotstats where ProcessOutcome == \"No Action\" and UserClassification != \"No Action\" group by ProcessOutcome,UserClassification"
	statsA, err := getStatsData(query, c)
	if err != nil {
		return cm, err
	}

	// now get all no actions
	query = "select count(ProcessOutcome) as count,ProcessOutcome,UserClassification from servisbotstats where ProcessOutcome == \"Cancel Subscription\" and UserClassification != \"Cancel Subscription\" group by ProcessOutcome,UserClassification"
	statsB, err := getStatsData(query, c)
	if err != nil {
		return cm, err
	}

	// now get all no actions
	query = "select count(ProcessOutcome) as count,ProcessOutcome,UserClassification from servisbotstats where ProcessOutcome == \"Cancel Autorenewal\" and UserClassification != \"Cancel Autorenewal\" group by ProcessOutcome,UserClassification"
	statsC, err := getStatsData(query, c)
	if err != nil {
		return cm, err
	}

	// we now have all te relevant values for the 3X3 matrix
	// let build our return matrix
	// we know that the statsA array is for NoAction
	cm = updateStatsStruct(statsA, statsB, statsC, statsTotal)

	return cm, nil
}

func getStatsData(query string, c *Connectors) ([]schema.Stat, error) {
	var stats []schema.Stat
	var stat *schema.Stat
	c.Info("Function getStatsData %s", query)
	res, err := c.Cluster.Query(query, &gocb.QueryOptions{})
	defer res.Close()
	if err != nil {
		c.Error("Function getStatsData (query) %v", err)
		return stats, err
	}

	// iterate through each object
	// struct with int64,string,string
	for res.Next() {
		err := res.Row(&stat)
		c.Trace("Function getStatsData data %v", stat)
		if err != nil {
			c.Error("Function getStatsData (next loop) %v", err)
			break
		}
		stats = append(stats, *stat)
	}

	// always check for errors after iterating
	err = res.Err()
	if err != nil {
		return stats, err
	}
	c.Trace("Function getStatsData alldata %v", stats)
	return stats, nil
}

func updateStatsStruct(in ...[]schema.Stat) *schema.ConfusionMatrix {

	var cm = &schema.ConfusionMatrix{}

	for _, item := range in[0] {
		switch item.UserClassification {
		case "Cancel Subscription":
			cm.NoAction.Cancel = item.Count
			break
		case "Cancel Autorenewal":
			cm.NoAction.CancelAR = item.Count
			break
		}
	}

	for _, item := range in[1] {
		switch item.UserClassification {
		case "No Action":
			cm.Cancel.NoAction = item.Count
			break
		case "Cancel Autorenewal":
			cm.Cancel.CancelAR = item.Count
			break
		}
	}

	for _, item := range in[2] {
		switch item.UserClassification {
		case "No Action":
			cm.CancelAR.NoAction = item.Count
			break
		case "Cancel Subscription":
			cm.CancelAR.Cancel = item.Count
			break
		}
	}

	for _, item := range in[3] {
		switch item.ProcessOutcome {
		case "No Action":
			cm.NoAction.NoAction = item.Count - (cm.NoAction.Cancel + cm.NoAction.CancelAR)
			break
		case "Cancel Subscription":
			cm.Cancel.Cancel = item.Count - (cm.Cancel.NoAction + cm.Cancel.CancelAR)
			break
		case "Cancel Autorenewal":
			cm.CancelAR.CancelAR = item.Count - (cm.CancelAR.NoAction + cm.CancelAR.Cancel)
			break
		}
	}
	return cm
}
