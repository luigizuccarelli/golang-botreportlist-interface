package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-s3bucket-manager/pkg/connectors"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Inject (force) readAll test error")
}

// TestAllHandlers - main test entry point
func TestAllHandlers(t *testing.T) {

	logger := &simple.Logger{Level: "trace"}

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v2/sys/info/isalive", nil)
		connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(IsAlive)
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

	t.Run("ListBucketHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/s3objects/list/12345", nil)
		conn := connectors.NewTestConnectors("../../tests/payload-smt-power.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ListBucketHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ListBucketHandler", rr.Code, STATUS))
		}
	})

	t.Run("GetObjectHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/s3objects/4324324324", nil)
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			GetObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "GetObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("GetObjectHandler : should fail (force error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/s3objects/4324324324", nil)
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		conn.Meta("true")
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			GetObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "GetObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("PutObjectHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		payload := "{\"test\":\"data\"}"
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/s3objects/4324324324", bytes.NewBuffer([]byte(payload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PutObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "PutObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("PutObjectHandler : should fail (force error readall)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/s3objects/4324324324", errReader(0))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PutObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "PutObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("PutObjectHandler : should fail (force error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		payload := "{\"test\":\"data\"}"
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/s3objects/4324324324", bytes.NewBuffer([]byte(payload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		conn.Meta("true")
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PutObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "PutObjectHandler", rr.Code, STATUS))
		}
	})

}
