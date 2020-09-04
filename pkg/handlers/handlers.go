package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/schema"
	"github.com/aws/aws-sdk-go/service/s3"
	gocb "github.com/couchbase/gocb/v2"
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

// ListHandler - handler that returns servisBOT accuracy
func ListHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
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
		msg := "ListHandler body data error : access forbidden %v"
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace("ListHandler request body : %s", string(body))

	errs := json.Unmarshal(body, &servisbotRequest)
	if errs != nil {
		msg := "ListHandler could not unmarshal input data from servisBOT to schema %v"
		con.Error(msg, errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}

	// check the jwt token
	_, err = verifyJwtToken(servisbotRequest.JwtToken)
	if err != nil {
		msg := "ListHandler verifyToken  %v"
		con.Error(msg, err)
		b := responseErrorFormat(http.StatusForbidden, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	// update the database
	res, err := con.GetList(vars["from"], vars["to"])
	if err != nil {
		msg := "ListHandler (get) couchbase  %v"
		con.Error(msg, err)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	response := &schema.Response{Code: http.StatusOK, Status: "OK", Message: "ListHandler retrieved data successfully ", Reports: res}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

// ReportUpdateHandler - handler that returns servisBOT accuracy
func ReportUpdateHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
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
		msg := "ReportUpdateHandler body data error : %v"
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace(" ReportUpdateHandler request body : %s", string(body))

	errs := json.Unmarshal(body, &servisbotRequest)
	if errs != nil {
		msg := "ReportUpdateHandler could not unmarshal input data from servisBOT to schema %v"
		con.Error(msg, errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}

	// check the jwt token
	_, err = verifyJwtToken(servisbotRequest.JwtToken)
	if err != nil {
		msg := "ReportUpdateHandler verifyToken  %v"
		con.Error(msg, err)
		b := responseErrorFormat(http.StatusForbidden, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	// update the database
	res, err := con.Upsert(servisbotRequest.Data.Id, servisbotRequest.Data.ServisbotStats, &gocb.UpsertOptions{})
	if err != nil {
		msg := "ReportUpdateHandler (post) couchbase  %v"
		con.Error(msg, err)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	response := &schema.Response{Code: http.StatusOK, Status: "OK", Message: fmt.Sprintf("ReportUpdateHandler posted data successfully %v", res)}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

// StatsHandler - handler that returns servisBOT accuracy
func StatsHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
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
		msg := "StatsHandler body data error : %v"
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace("StatsHandler request body : %s", string(body))

	errs := json.Unmarshal(body, &servisbotRequest)
	if errs != nil {
		msg := "StatsHandler could not unmarshal input data from servisBOT to schema %v"
		con.Error(msg, errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}

	// check the jwt token
	_, err = verifyJwtToken(servisbotRequest.JwtToken)
	if err != nil {
		msg := "StatsHandler verifyToken  %v"
		con.Error(msg, err)
		b := responseErrorFormat(http.StatusForbidden, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	// get the stats data from couchbase
	res, err := con.GetAllStats()
	if err != nil {
		msg := "StatsHandler (post) couchbase  %v"
		con.Error(msg, err)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	response := &schema.StatsResponse{Code: http.StatusOK, Status: "OK", Message: "StatsHandler retrieved data successfully", Stats: res}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

// ReportObjectHandler - handler that interfaces with s3 bucket
func ReportObjectHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var servisbotRequest *schema.ServisBOTRequest

	bucket := os.Getenv(AWSBUCKET)
	addHeaders(w, r)

	// read the jwt token data in the body
	// we don't use authorization header as the token can get quite large due to form data
	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ReportObjectHandler body data %v"
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

	key := servisbotRequest.Data.Id
	opts := &s3.GetObjectInput{Bucket: &bucket, Key: &key}
	data, e := con.GetObject(opts)
	if e != nil {
		msg := "ReportObjectHandler %v"
		con.Error(msg, e)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}

	con.Trace("ReportObjectHandler data %v", data)
	response := &schema.ReportResponse{Code: http.StatusOK, Status: "OK", Message: "ReportObjectHandler s3 object call successful", Report: *data}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

// IsAlive - readiness & liveliness probe
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
