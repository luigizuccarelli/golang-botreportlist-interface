package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-lytics-interface/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-lytics-interface/pkg/schema"
	"github.com/gorilla/mux"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

// ProfileHandler - handler that calls lytics audience endpoint
func ProfileHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var audience *schema.AudienceSchema
	var vars = mux.Vars(r)

	addHeaders(w, r)

	email := vars["email"]
	token := os.Getenv("TOKEN")
	if token == "" || email == "" {
		err := errors.New("input params token/email are empty")
		con.Error("ProfileHandler %v", err)
		b := responseError(w, "ProfileHandler %v", err)
		fmt.Fprintf(w, string(b))
		return
	}

	// we first check the smt_power audience
	res := strings.NewReplacer("{email}", email, "{token}", token, "{audience}", "smt_power")
	// Replace all pairs.
	url := res.Replace(os.Getenv("URL"))
	body, errs := makeRequest(url, con)
	if errs != nil {
		con.Error("ProfileHandler %v", errs)
		b := responseError(w, "ProfileHandler %v", errs)
		fmt.Fprintf(w, string(b))
		return
	}

	errs = json.Unmarshal(body, &audience)
	if errs != nil {
		msg := "ProfileHandler could not unmarshal profile message data to schema (smt_power) %v"
		con.Error(msg, errs)
		b := responseError(w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}

	// we check if status is 200 - this should always work
	// if not it means lytics is offline
	if audience.Status == 200 {
		// we should usually get success
		// check in  "highly engaged audience" list
		data, bFlag := handleResponse(audience, "smt_power", con)
		if bFlag {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, string(data))
			return
		} else {
			// check in "currently engaged audience" list
			data, bFlag = handleResponse(audience, "smt_active", con)
			res = strings.NewReplacer("{email}", email, "{token}", token, "{audience}", "smt_active")
			// Replace all pairs.
			url = res.Replace(os.Getenv("URL"))
			body, errs = makeRequest(url, con)
			if errs != nil {
				con.Error("ProfileHandler %v", errs)
				b := responseError(w, "ProfileHandler %v", errs)
				fmt.Fprintf(w, string(b))
				return
			}

			errs = json.Unmarshal(body, &audience)
			if errs != nil {
				msg := "ProfileHandler could not unmarshal profile message data to schema (smt_active) %v"
				con.Error(msg, errs)
				b := responseError(w, msg, errs)
				fmt.Fprintf(w, string(b))
				return
			}

			data, _ := handleResponse(audience, "smt_active", con)
			fmt.Fprintf(w, string(data))
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	msg := "ProfileHandler %v"
	err := errors.New("lytics request error")
	con.Error(msg, err)
	b := responseError(w, msg, err)
	fmt.Fprintf(w, string(b))
	return
}

// IsAlive - endpoint call used for openshift readiness and liveliness probes
func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \""+os.Getenv("NAME")+"\" }")
	return
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// responsError - utility function
func responseError(w http.ResponseWriter, msg string, val ...interface{}) []byte {
	var b []byte
	response := &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf(msg, val...)}
	w.WriteHeader(http.StatusInternalServerError)
	b, _ = json.MarshalIndent(response, "", "	")
	return b
}

// handleResponse - utility function
func handleResponse(a *schema.AudienceSchema, list string, con connectors.Clients) ([]byte, bool) {
	var b []byte
	// should not be empty
	if a.Message == "success" {
		if contains(a.Data.Segments, list) {
			msg := "engaged : " + list
			con.Debug("Function handleResponse "+msg, "")
			response := &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: msg}
			b, _ = json.MarshalIndent(response, "", "	")
			return b, true
		} else {
			msg := "engaged : none"
			con.Debug("Function handleResponse "+msg, "")
			response := &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: msg}
			b, _ = json.MarshalIndent(response, "", "	")
			return b, false
		}
	}
	msg := "engaged : not found"
	con.Debug("Function handleResponse "+msg, "")
	response := &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: msg}
	b, _ = json.MarshalIndent(response, "", "	")
	return b, false
}

// makeRequest - private utility function
func makeRequest(url string, con connectors.Clients) ([]byte, error) {
	var b []byte
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := con.Do(req)
	if err != nil {
		con.Error("Function makeRequest http request %v", err)
		return b, err
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		con.Error("Function makeRequest %v", e)
		return b, err
	}
	con.Debug("Function makeRequest response from middleware %s", string(body))
	return body, nil
}

// contains - private function that iterates through the []string
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
