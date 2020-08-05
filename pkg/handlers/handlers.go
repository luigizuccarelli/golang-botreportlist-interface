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
		in = &s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String("Email"), StartAfter: aws.String(vars["lastobject"])}
	} else {
		in = &s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String("Email")}
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
