package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-s3bucket-manager/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-s3bucket-manager/pkg/schema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
	HANDLERESPONSE  string = "Function handleResponse "
)

// ListBucketHandler - handler that interfaces with s3 bucket
func ListBucketHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var in *s3.ListObjectsV2Input
	var files []schema.S3FileMetaData
	vars := mux.Vars(r)

	//sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	bucket := os.Getenv("AWS_BUCKET")

	// Get the list of items
	if vars["lastobject"] != "" {
		in = &s3.ListObjectsV2Input{Bucket: aws.String(bucket), StartAfter: aws.String(vars["lastobject"])}
	} else {
		in = &s3.ListObjectsV2Input{Bucket: aws.String(bucket)}
	}

	resp, err := con.ListObjectsV2(in)
	if err != nil {
		con.Error("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		file := &schema.S3FileMetaData{Name: *item.Key, LastModified: *item.LastModified, Size: *item.Size, StorageClass: *item.StorageClass}
		files = append(files, *file)
	}
	con.Trace("ListBucketHandler found %d items in bucket %s", len(resp.Contents), bucket)
	response := &schema.Response{Code: http.StatusOK, Status: "OK", Message: "ListBucketHandler s3 object call successful", Payload: files}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

// GetObjectHandler - handler that interfaces with s3 bucket
func GetObjectHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	bucket := os.Getenv("AWS_BUCKET")
	vars := mux.Vars(r)

	filename := vars["key"]
	opts := &s3.GetObjectInput{Bucket: &bucket, Key: &filename}
	data, e := con.GetObject(opts)
	if e != nil {
		msg := "GetObjectHandler %v"
		con.Error(msg, e)
		b := responseError(w, msg, e)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace("GetObjectHandler data %s", string(data))
	response := &schema.Response{Code: http.StatusOK, Status: "OK", Message: "ListBucketHandler s3 object call successful", EmailContent: string(data)}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

// PutObjectHandler - handler that interfaces with s3 bucket
func PutObjectHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	bucket := os.Getenv("AWS_BUCKET")
	vars := mux.Vars(r)

	filename := vars["key"]
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "PutObjectHandler %v"
		con.Error(msg, err)
		b := responseError(w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	opts := &s3.PutObjectInput{Bucket: &bucket, Key: &filename, Body: aws.ReadSeekCloser(strings.NewReader(string(data)))}
	res, e := con.PutObject(opts)
	if e != nil {
		msg := "PutObjectHandler %v"
		con.Error(msg, e)
		b := responseError(w, msg, e)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace("PutObjectHandler data %s", res)

	response := &schema.Response{Code: http.StatusOK, Status: "OK", Message: "PutObjectHandler s3 object call successful"}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	// add header
	addHeaders(w, r)
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \""+os.Getenv("NAME")+"\" }")
	return
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("API-KEY") != "" {
		w.Header().Set("API_KEY_PT", r.Header.Get("API_KEY"))
	}
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// responsError - utility function
func responseError(w http.ResponseWriter, msg string, val ...interface{}) []byte {
	var b []byte
	response := &schema.Response{Code: http.StatusInternalServerError, StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf(msg, val...)}
	w.WriteHeader(http.StatusInternalServerError)
	b, _ = json.MarshalIndent(response, "", "	")
	return b
}
