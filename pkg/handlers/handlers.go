package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-middleware-interface/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-middleware-interface/pkg/schema"
	"github.com/gorilla/mux"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var si *schema.SchemaInterface
	var emails []schema.EmailProfile
	var vars = mux.Vars(r)

	addHeaders(w, r)

	// search for the affiliate token
	token, err := getToken(vars["affiliateid"])
	if err != nil {
		con.Error("Token  %v", err)
		b := responseFormat(true, w, "Token  %v", err)
		fmt.Fprintf(w, string(b))
		return
	}

	url := os.Getenv("URL") + "/account/emailaddress?email=" + vars["email"]
	body, errs := makeRequest(url, token, con)
	if errs != nil {
		con.Error(" %v", errs)
		b := responseFormat(true, w, " %v", errs)
		fmt.Fprintf(w, string(b))
		return
	}

	errs = json.Unmarshal(body, &emails)
	if errs != nil {
		msg := "Could not unmarshal email message data to schema %v"
		con.Error(msg, errs)
		b := responseFormat(true, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	password := getPassword(emails, vars["email"])
	// now make the call to get all data
	url = os.Getenv("URL") + "data/username/" + emails[0].ID.UserName + "/password/" + password
	body, errs = makeRequest(url, token, con)
	if err != nil {
		con.Error(" %v", err)
		b := responseFormat(true, w, " %v", err)
		fmt.Fprintf(w, string(b))
		return
	}

	// only used in testing to intercept and inject profile data
	if os.Getenv("TESTING") != "" && os.Getenv("TESTING") == "true" {
		body, err = injectJsonProfile(body)
		// just report the error
		if err != nil {
			con.Error("injectJsonProfile - testing  %v", err)
		}
	}

	errs = json.Unmarshal(body, &si)
	if errs != nil {
		msg := "Could not unmarshal profile message data to schema %v"
		con.Error(msg, errs)
		b := responseFormat(true, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	msg := "MW data successfully retrieved "
	con.Trace(msg+" %v", si)
	response := &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: fmt.Sprintf(msg), Payload: si}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
	return
}

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

// responsFormat - utility function
func responseFormat(err bool, w http.ResponseWriter, msg string, val ...interface{}) []byte {
	var b []byte
	// for testing
	if err {
		response := &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf(msg, val...)}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ = json.MarshalIndent(response, "", "	")
		return b
	}
	return b
}

// makeRequest - private utility function
func makeRequest(url string, token string, con connectors.Clients) ([]byte, error) {
	var b []byte
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("token", token)
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

// getToken - private utility call
func getToken(affiliate string) (string, error) {
	var tokens []schema.TokenDetail
	var token string = ""

	if affiliate == "" {
		return "", errors.New("Affiliate parameter is empty")
	}
	errs := json.Unmarshal([]byte(os.Getenv("TOKEN")), &tokens)
	if errs != nil {
		return "", errors.New("Unmarshalling token struct")
	}
	for x, _ := range tokens {
		if affiliate == tokens[x].Name {
			token = tokens[x].Token
			break
		}
	}
	if token == "" {
		return "", errors.New("Token not found")
	}
	return token, nil
}

func getPassword(emails []schema.EmailProfile, user string) string {
	var pwd string
	for x, _ := range emails {
		if emails[x].ID.UserName == user {
			pwd = emails[x].Password
			break
		}
	}
	return pwd
}

func injectJsonProfile(data []byte) ([]byte, error) {
	// only used for testing
	var b []byte
	var err error
	b, err = ioutil.ReadFile("../../tests/payload.json")
	if err != nil {
		return b, err
	}
	return b, err
}
