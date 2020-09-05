// +build fake

package connectors

import (
	"fmt"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/schema"
	"github.com/aws/aws-sdk-go/service/s3"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

func TestConnections(t *testing.T) {

	var logger = &simple.Logger{Level: "trace"}

	t.Run("Logging : should pass", func(t *testing.T) {
		con := &Connectors{Bucket: &FakeBucket{}, Cluster: &FakeCluster{}, S3Service: &FakeS3{}, Logger: logger}
		con.Info("Log Info")
		con.Debug("Log Debug")
		con.Trace("Log Trace")
		con.Error("Log Error")
	})

	t.Run("GetObject : should pass", func(t *testing.T) {
		con := &Connectors{Bucket: &FakeBucket{}, Cluster: &FakeCluster{}, S3Service: &FakeS3{}, Logger: logger}
		bucket := "Email/"
		key := "12345"
		opts := &s3.GetObjectInput{Bucket: &bucket, Key: &key}
		data, err := con.GetObject(opts)
		if err != nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should be nil) -  got (%v) wanted (%v)", "GetObject", err, nil))
		}
		con.Info("Data result %v", data)
	})

	t.Run("GetObject : should fail (forced error)", func(t *testing.T) {
		con := &Connectors{Bucket: &FakeBucket{}, Cluster: &FakeCluster{}, S3Service: &FakeS3{Force: "true"}, Logger: logger}
		bucket := "Email/"
		key := "12345"
		opts := &s3.GetObjectInput{Bucket: &bucket, Key: &key}
		_, err := con.GetObject(opts)
		if err == nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should not be nil) -  got (%v) wanted (%v)", "GetObject", nil, "error"))
		}
		con.Info("Data result %v", err)
	})

	t.Run("GetObject : should fail (forced data error)", func(t *testing.T) {
		con := &Connectors{Bucket: &FakeBucket{}, Cluster: &FakeCluster{}, S3Service: &FakeS3{Force: "data"}, Logger: logger}
		bucket := "Email/"
		key := "12345"
		opts := &s3.GetObjectInput{Bucket: &bucket, Key: &key}
		_, err := con.GetObject(opts)
		if err == nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should not be nil) -  got (%v) wanted (%v)", "GetObject", nil, "error"))
		}
		con.Info("Data result %v", err)
	})

	t.Run("Upsert : should pass", func(t *testing.T) {
		con := &Connectors{Bucket: &FakeBucket{}, Cluster: &FakeCluster{}, S3Service: &FakeS3{}, Logger: logger}
		data, err := con.Upsert("123456", &schema.ReportContent{}, &gocb.UpsertOptions{})
		if err != nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should be nil) -  got (%v) wanted (%v)", "Upsert", err, nil))
		}
		con.Info("Data result %v", data)
	})

	t.Run("Upsert : should fail (forced error)", func(t *testing.T) {
		con := &Connectors{Bucket: &FakeBucket{Force: "error"}, Cluster: &FakeCluster{Force: "error"}, S3Service: &FakeS3{}, Logger: logger}
		_, err := con.Upsert("123456", &schema.ReportContent{}, &gocb.UpsertOptions{})
		if err == nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should not be nil) -  got (%v) wanted (%v)", "Upsert", nil, "error"))
		}
		con.Info("Data result %v", err)
	})

	t.Run("GetList : should pass", func(t *testing.T) {
		con := &Connectors{Bucket: &FakeBucket{}, Cluster: &FakeCluster{}, S3Service: &FakeS3{}, Logger: logger}
		data, err := con.GetList("0", "10")
		if err != nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should be nil) -  got (%v) wanted (%v)", "GetList", err, nil))
		}
		con.Info("Data result %v", data)
	})

	t.Run("GetList : should fail (forced error)", func(t *testing.T) {
		con := &Connectors{Bucket: &FakeBucket{}, Cluster: &FakeCluster{Force: "true"}, S3Service: &FakeS3{}, Logger: logger}
		data, err := con.GetList("0", "10")
		if err == nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should not be nil) -  got (%v) wanted (%v)", "GetList", err, nil))
		}
		con.Info("Data result %v", data)
	})

	t.Run("GetListCount : should pass", func(t *testing.T) {
		con := &Connectors{Bucket: &FakeBucket{}, Cluster: &FakeCluster{}, S3Service: &FakeS3{}, Logger: logger}
		data, err := con.GetListCount()
		if err != nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should be nil) -  got (%v) wanted (%v)", "GetListCount", err, nil))
		}
		con.Info("Data result %v", data)
	})

	t.Run("GetAllStats : should pass", func(t *testing.T) {
		con := &Connectors{Bucket: &FakeBucket{}, Cluster: &FakeCluster{}, S3Service: &FakeS3{}, Logger: logger}
		data, err := con.GetAllStats()
		if err != nil {
			t.Errorf(fmt.Sprintf("Function (%s) assert (error should be nil) -  got (%v) wanted (%v)", "GetAllStats", err, nil))
		}
		con.Info("Data result %v", data)
	})

}
