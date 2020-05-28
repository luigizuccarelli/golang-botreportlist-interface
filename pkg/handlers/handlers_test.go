package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-lytics-interface/pkg/connectors"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
)

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

	t.Run("ProfileHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/profile/", nil)
		conn := connectors.NewTestConnectors("../../tests/payload-smt-power.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"email": "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/profile/", nil)
		conn := connectors.NewTestConnectors("../../tests/payload-smt-active.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"email": "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/profile/", nil)
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"email": "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/profile/", nil)
		conn := connectors.NewTestConnectors("../../tests/payload-not-found.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"email": "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should fail (no email)", func(t *testing.T) {
		var STATUS int = 500
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/profile/", nil)
		conn := connectors.NewTestConnectors("../../tests/payload-not-found.json", STATUS, logger)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should fail (unmarshal json)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/profile/", nil)
		conn := connectors.NewTestConnectors("../../tests/fail.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"email": "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should fail (json status 500 return error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/profile/", nil)
		conn := connectors.NewTestConnectors("../../tests/payload-status-code.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"email": "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should fail (force do request error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/profile/", nil)
		conn := connectors.NewTestConnectors("../../tests/payload-status-code.json", STATUS, logger)
		// force request Do error
		conn.Meta("true")

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"email": "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

}
