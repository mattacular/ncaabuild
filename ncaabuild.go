package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	buildEnv string
	planId   string
	decoded  map[string]interface{}
	username string
	password string
)

func applyEnvOptions() {
	envUser := os.Getenv("NCAA_BARCA_BAMBOO_USER")
	envPass := os.Getenv("NCAA_BARCA_BAMBOO_PASS")

	if username == "" && envUser != "" {
		username = envUser
	}

	if password == "" && envPass != "" {
		password = envPass
	}
}

func init() {
	flag.StringVar(&username, "username", "", "Your NCAA Bamboo5 username")
	flag.StringVar(&password, "password", "", "Your NCAA Bamboo5 password")
}

func main() {
	username = ""
	password = ""

	flag.Parse()
	applyEnvOptions()

	if (username == "" || password == "") {
		fmt.Println("You must provide a NCAA Bamboo5 username and password, either by setting environment variables \"NCAA_BARCA_BAMBOO_USER\" and \"NCAA_BARCA_BAMBOO_PASS\" or by passing these values as option flags \"--user\" and \"--pass\"")
		return
	}

	if args := flag.Args(); len(args) > 0 {
		buildEnv = args[0:1][0]
		planId = "barcelona-" + buildEnv
	}

	// validate environment arg
	switch buildEnv {
	case "staging":
		planId = "barcelona-prod"
	case "qa":
	case "dev":
		break
	default:
		buildEnv = "qa"
		planId = "barcelona-qa"
	}

	// set up the GET request to check if the plan is already running
	reqUrl := "http://ncaa-build.services.56m.vgtf.net/rest/api/latest/plan/" + planId + ".json"
	client := &http.Client{}

	req, _ := http.NewRequest("GET", reqUrl, nil)
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	} else {
		defer resp.Body.Close() // close the connection once this function finishes
	}

	// read in the JSON response body and decode it
	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &decoded); err != nil {
		panic(err)
	}

	// make sure the requested plan is not already running
	if decoded["isActive"] != false {
		fmt.Println("Plan \"" + planId + "\" already running. Please wait for it to finish before running the plan again.")
		return
	}

	// set up the POST request to trigger a new build on the given plan
	reqUrl = "http://ncaa-build.services.56m.vgtf.net/rest/api/latest/queue/" + planId
	postData := url.Values{}
	postData.Set("os_authType", "basic")
	postData.Set("field2", "value")

	buildReq, _ := http.NewRequest("POST", reqUrl, bytes.NewBufferString(postData.Encode()))
	credentials := []byte(username + ":" + password)
	buildReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	buildReq.Header.Add("Authorization", "Basic " + base64.StdEncoding.EncodeToString(credentials))

	buildResp, buildErr := client.Do(buildReq)

	if buildErr != nil {
		panic(err)
	}

	if (buildResp.Status == "200 OK") {
		buildResp.Body.Close() // close the connection once this function finishes
		fmt.Println("Build trigger sent, \"" + planId + "\" will start building momentarily.")
	} else {
		fmt.Println("Build trigger was not sent. Response:", buildResp.Status)
	}
}
