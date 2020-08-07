package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-s3bucket-manager/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-s3bucket-manager/pkg/schema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
	HANDLERESPONSE  string = "Function handleResponse "
	AWSBUCKET       string = "AWS_BUCKET"
	AWSREPORTBUCKET string = "AWS_REPORT_BUCKET"
	CHANNEL         string = "Email/"
	EMAIL           string = "Email"
)

// ListBucketHandler - handler that interfaces with s3 bucket
func ListBucketHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var in *s3.ListObjectsV2Input
	var files []schema.S3FileMetaData
	var servisbotRequest *schema.ServisBOTRequest
	vars := mux.Vars(r)

	addHeaders(w, r)
	// read the jwt token data in the body
	// we don't use authorization header as the token can get quite large due to form data
	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "Body data (JWT) %v"
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace("Request body : %s", string(body))

	// unmarshal result from mw backend
	errs := json.Unmarshal(body, &servisbotRequest)
	if errs != nil {
		msg := "GenericHandler could not unmarshal input data from servisBOT to schema %v"
		con.Error(msg, errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}

	// check the jwt token
	_, err = verifyJwtToken(servisbotRequest.JwtToken)
	if err != nil {
		msg := "ListBucketHandler verifyToken  %v"
		con.Error(msg, err)
		b := responseErrorFormat(http.StatusForbidden, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	bucket := os.Getenv(AWSREPORTBUCKET)

	// Get the list of items
	if vars["lastobject"] != "" && vars["lastobject"] != "false" {
		in = &s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String(EMAIL), StartAfter: aws.String(vars["lastobject"])}
	} else {
		in = &s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String(EMAIL)}
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

// EmailObjectHandler - handler that interfaces with s3 bucket
func EmailObjectHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var servisbotRequest *schema.ServisBOTRequest

	bucket := os.Getenv(AWSBUCKET)
	vars := mux.Vars(r)

	addHeaders(w, r)
	// read the jwt token data in the body
	// we don't use authorization header as the token can get quite large due to form data
	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "EmailObjectHandler body data %v"
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace("EmailObjectHandler request body : %s", string(body))

	// unmarshal result from mw backend
	errs := json.Unmarshal(body, &servisbotRequest)
	if errs != nil {
		msg := "EmailObjectHandler could not unmarshal input data from servisBOT to schema %v"
		con.Error(msg, errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}

	// check the jwt token
	_, err = verifyJwtToken(servisbotRequest.JwtToken)
	if err != nil {
		msg := "EmailObjectHandler verifyToken  %v"
		con.Error(msg, err)
		b := responseErrorFormat(http.StatusForbidden, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	filename := vars["key"]
	opts := &s3.GetObjectInput{Bucket: &bucket, Key: &filename}
	data, e := con.GetObject(opts)
	if e != nil {
		msg := "EmailObjectHandler %v"
		con.Error(msg, e)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace("EmailObjectHandler data %s", string(data))
	response := &schema.Response{Code: http.StatusOK, Status: "OK", Message: "EmailObjectHandler s3 object call successful", Email: string(data)}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

// ReportObjectHandler - handler that interfaces with s3 bucket
func ReportObjectHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var servisbotRequest *schema.ServisBOTRequest
	var response *schema.Response
	var report *schema.ReportContent

	bucket := os.Getenv(AWSREPORTBUCKET)
	vars := mux.Vars(r)

	addHeaders(w, r)

	// read the jwt token data in the body
	// we don't use authorization header as the token can get quite large due to form data
	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ReportObjectHandler body data error : access forbidden %v"
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace("ReportObjectHandler request body : %s", string(body))

	// unmarshal result from mw backend
	errs := json.Unmarshal(body, &servisbotRequest)
	if errs != nil {
		msg := "ReportObjectHandler could not unmarshal input data from servisBOT to schema %v"
		con.Error(msg, errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}

	// check the jwt token
	_, err = verifyJwtToken(servisbotRequest.JwtToken)
	if err != nil {
		msg := "ReportObjectHandler verifyToken  %v"
		con.Error(msg, err)
		b := responseErrorFormat(http.StatusForbidden, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	// check for request method
	filename := CHANNEL + vars["key"]
	if con.GetMode() == "pull" {
		opts := &s3.GetObjectInput{Bucket: &bucket, Key: &filename}
		res, er := con.GetObject(opts)
		if er != nil {
			msg := "ReportObjectHandler %v"
			con.Error(msg, er)
			b := responseErrorFormat(http.StatusInternalServerError, w, msg, er)
			fmt.Fprintf(w, string(b))
			return
		}
		// unmarshal result from mw backend
		json.Unmarshal(res, &report)
		con.Trace("ReportObjectHandler (get) data %s", res)
		response = &schema.Response{Code: http.StatusOK, Status: "OK", Message: "ReportObjectHandler s3 object call (get) successful", Report: report}
	} else {
		// This is s POST
		opts := &s3.PutObjectInput{Bucket: &bucket, Key: &filename, Body: aws.ReadSeekCloser(strings.NewReader(servisbotRequest.Data))}
		res, e := con.PutObject(opts)
		if e != nil {
			msg := "ReportObjectHandler %v"
			con.Error(msg, e)
			b := responseErrorFormat(http.StatusInternalServerError, w, msg, e)
			fmt.Fprintf(w, string(b))
			return
		}
		con.Trace("ReportObjectHandler (post) data %s", res)
		response = &schema.Response{Code: http.StatusOK, Status: "OK", Message: "ReportObjectHandler s3 object call (post) successful"}
	}

	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

// GetStatsHandler - handler that returns servisBOT accuracy
func GetStatsHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var servisbotRequest *schema.ServisBOTRequest

	addHeaders(w, r)
	// read the jwt token data in the body
	// we don't use authorization header as the token can get quite large due to form data
	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "GetStatsHandler body data %v"
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace("GetStatsHandler request body : %s", string(body))

	// unmarshal result from mw backend
	errs := json.Unmarshal(body, &servisbotRequest)
	if errs != nil {
		msg := "GetStatsHandler could not unmarshal input data from servisBOT to schema %v"
		con.Error(msg, errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}

	// check the jwt token
	_, err = verifyJwtToken(servisbotRequest.JwtToken)
	if err != nil {
		msg := "GetStatsHandler verifyToken  %v"
		con.Error(msg, err)
		b := responseErrorFormat(http.StatusForbidden, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	bucket := os.Getenv(AWSREPORTBUCKET)
	filename := "Stats/stats.json"
	opts := &s3.GetObjectInput{Bucket: &bucket, Key: &filename}
	res, er := con.GetObject(opts)
	if er != nil {
		msg := "GetStatsHandler reading stats  %v"
		con.Error(msg, er)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, er)
		fmt.Fprintf(w, string(b))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(res))
	return
}

// StatsHandler - handler that interfaces with s3 bucket
func StatsHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var stats *schema.Stats
	var listOpts *s3.ListObjectsV2Input

	vars := mux.Vars(r)
	bucket := os.Getenv(AWSREPORTBUCKET)

	// we don't need to worry about jwt
	if vars["init"] != "" && vars["init"] == "true" {
		con.Trace("StatsHandler init set to true - starting re-run")
		// re-run from the first record
		listOpts = &s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String(EMAIL)}
	} else {
		// Get the list of items
		// Check if we have a lastobject item
		con.Trace("StatsHandler init not set")
		filename := "Stats/stats.json"
		opts := &s3.GetObjectInput{Bucket: &bucket, Key: &filename}
		res, _ := con.GetObject(opts)
		// update our schema
		err := json.Unmarshal([]byte(res), &stats)
		if err != nil {
			stats = &schema.Stats{}
			con.Error("StatsHandler converting json %v", err)
		}
		if stats.LastObject != "" {
			listOpts = &s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String(EMAIL), StartAfter: aws.String(stats.LastObject)}
		} else {
			listOpts = &s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String(EMAIL)}
		}
	}

	var wg sync.WaitGroup
	if os.Getenv("TESTING") != "" && os.Getenv("TESTING") == "true" {
		wg.Add(1)
	}

	go calculateStats(con, listOpts, &wg)

	if os.Getenv("TESTING") != "" && os.Getenv("TESTING") == "true" {
		wg.Wait()
	}

	response := &schema.Response{Code: http.StatusOK, Status: "OK", Message: fmt.Sprintf("StatsHandler process started in background - check timestamp %d", time.Now().Unix())}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

func calculateStats(con connectors.Clients, listOpts *s3.ListObjectsV2Input, wg *sync.WaitGroup) error {

	if os.Getenv("TESTING") != "" && os.Getenv("TESTING") == "true" {
		defer wg.Done()
	}
	con.Trace("Function getStats opts %v", listOpts)
	bucket := os.Getenv(AWSREPORTBUCKET)
	resp, err := con.ListObjectsV2(listOpts)
	if err != nil {
		msg := "Function calculateStats unable to list items in bucket %q, %v"
		con.Error(msg, bucket, err)
		return err
	}

	var name string = ""
	var count float64 = 0.0
	var accuracy float64 = 0.0

	con.Trace("Function calculateStats listObjects count %d", len(resp.Contents))
	for _, item := range resp.Contents {
		name = *item.Key
		count++
		accuracy = accuracy + getObject(con, bucket, name)
	}

	// check for more objects in the bucket
	//if *resp.IsTruncated {
	for {
		resp, err = con.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String(EMAIL), StartAfter: aws.String(name)})
		if err != nil {
			msg := "Function calculateStats unable to list items in bucket %q, %v"
			con.Error(msg, bucket, err)
			return err
		}
		con.Trace("Function calculateStats listObjects (second loop) count %d", len(resp.Contents))
		for _, item := range resp.Contents {
			name = *item.Key
			count++
			accuracy = accuracy + getObject(con, bucket, name)
		}
		if !*resp.IsTruncated {
			break
		}
	}
	//}

	con.Trace("Function calculateStats last object %s", name)
	con.Trace("Function calculateStats found %f items in bucket %s", count, bucket)
	con.Trace("Function calculateStats success items in bucket %f", accuracy)
	con.Trace("Function calculateStats bot accuracy %f", accuracy)
	s := &schema.Stats{RecordCount: count, SuccessCount: accuracy, Accuracy: (accuracy / count), LastObject: name, LastUpdated: time.Now().Unix()}
	con.Trace("Function calculateStats struct %v", s)
	b, _ := json.MarshalIndent(s, "", "	")
	filename := "Stats/stats.json"
	// store to s3
	opts := &s3.PutObjectInput{Bucket: &bucket, Key: &filename, Body: aws.ReadSeekCloser(strings.NewReader(string(b)))}
	_, e := con.PutObject(opts)
	if e != nil {
		msg := "Function calculateStats putobject %v"
		con.Error(msg, e)
		return e
	}
	con.Trace("Function calculateStats putobject succeeded %s", string(b))
	// all good
	return nil
}

func getObject(con connectors.Clients, bucket string, key string) float64 {
	var rc *schema.ReportContent
	const MSG string = "Function getObject %v"

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := con.GetObject(input)
	if err != nil {
		con.Error(MSG, err)
		return 0.0
	}

	// unmarshal result from mw backend
	errs := json.Unmarshal([]byte(result), &rc)
	if errs != nil {
		con.Error(MSG+" unmarshalling data to schema", errs)
		return 0.0
	}

	if rc.Success != "" {
		if rc.Success == "true" {
			return 1.0
		} else {
			return 0.0
		}
	}
	return 0.0
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	// add header (cors) override for vuejs FE
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
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept-Language")
}

// responsErrorFormat - utility function
func responseErrorFormat(code int, w http.ResponseWriter, msg string, val ...interface{}) []byte {
	var b []byte
	response := &schema.Response{Code: code, Status: "ERROR", Message: fmt.Sprintf(msg, val...)}
	w.WriteHeader(code)
	b, _ = json.MarshalIndent(response, "", "	")
	return b
}

// verifyJwtToken - private function
func verifyJwtToken(tokenStr string) (*schema.Credentials, error) {
	var creds *schema.Credentials

	if tokenStr == "" {
		return creds, errors.New("jwt token is invalid/empty")
	}
	// local function
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRETKEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["user"] == nil || claims["customerNumber"] == nil {
			return creds, errors.New("JWT invalid user/customerNumber empty")
		}
		user := claims["user"].(string)
		cn := claims["customerNumber"].(string)
		creds = &schema.Credentials{User: user, Password: "", CustomerNumber: cn}
		return creds, nil
	}
	return creds, errors.New("jwt token is invalid")
}
