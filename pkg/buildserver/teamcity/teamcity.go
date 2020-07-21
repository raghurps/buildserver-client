package teamcity

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// TCClient is client object to talk to teamcity
type TCClient struct {
	client    *http.Client
	token     string
	serverURL string
}

// NewTeamcityClient ...
func NewTeamcityClient(
	requestTimeout, dialTimeout, tlsHandshakeTimeout time.Duration,
	serverURL, token string,
	insecure bool,
) *TCClient {
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: dialTimeout,
		}).Dial,
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: insecure},
	}

	client := &http.Client{
		Timeout:   requestTimeout,
		Transport: tr,
	}

	return &TCClient{
		client:    client,
		serverURL: serverURL,
		token:     token,
	}
}

// GetBuild returns build details
// for the provided id
func (t *TCClient) GetBuild(id int, buildDetails interface{}) (err error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/app/rest/builds/id:%d", t.serverURL, id), nil)
	req.Header.Add("Authorization", t.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = json.Unmarshal(body, &buildDetails)
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}

/*
StartBuild adds a build to the build queue

buildTypeID is the unique ID of a build pipeline

branch is the branch name on which build will be triggered

params is a map containing env variables and other overrides that
user wants to provide
*/
func (t *TCClient) StartBuild(buildTypeID, branch string, params map[string]string) (int, error) {
	var buildDetails TCBuildDetails

	payload := TCBuildPayload{
		BuildType: TCBuildType{
			ID: buildTypeID,
		},
		Comment: TCBuildComment{
			Text: "Build started by Prow",
		},
		Properties: TCBuildProperties{
			Property: []TCBuildProperty{},
		},
		Personal:   "False",
		BranchName: branch,
	}

	// Add params to properties
	for k, v := range params {
		payload.Properties.Property = append(payload.Properties.Property, TCBuildProperty{k, v})
	}

	requestPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	log.Println(string(requestPayload))

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/app/rest/buildQueue", t.serverURL),
		bytes.NewBuffer(requestPayload))
	req.Header.Add("Authorization", t.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	err = json.Unmarshal(body, &buildDetails)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	log.Println(buildDetails)
	return buildDetails.ID, nil
}

// CancelQueuedBuild cancels a build that is currently
// queued in the BuildQueue
// If the build has already started or finished,
// this call will fail
func (t *TCClient) CancelQueuedBuild(id int) error {
	//var buildDetails TCBuildDetails

	payload := TCBuildStopPayload{
		Comment:        "Build cancelled by prow",
		ReaddIntoQueue: "false",
	}

	requestPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println(string(requestPayload))

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/app/rest/buildQueue/%d", t.serverURL, id),
		bytes.NewBuffer(requestPayload))
	req.Header.Add("Authorization", t.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	/* err = json.Unmarshal(body, &buildDetails)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println(buildDetails) */
	log.Println(string(body))
	return nil
}

// StopBuild stops a running build
func (t *TCClient) StopBuild(id int) error {
	//var buildDetails TCBuildDetails

	payload := TCBuildStopPayload{
		Comment:        "Build stopped by prow",
		ReaddIntoQueue: "false",
	}

	requestPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println(string(requestPayload))

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/app/rest/builds/%d", t.serverURL, id),
		bytes.NewBuffer(requestPayload))
	req.Header.Add("Authorization", t.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	/* err = json.Unmarshal(body, &buildDetails)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println(buildDetails) */
	log.Println(string(body))
	return nil
}

/*
GetArtifactTextFile fetches the content of an artifact file

path is the relative path of the file in teamcity artifacts

id is the build id from which the artifact will be fetched

It returns content of the file as array of bytes, content type of that file and error object if any
*/
func (t *TCClient) GetArtifactTextFile(path string, id int) ([]byte, string, error) {
	var fileContent []byte
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/app/rest/builds/id:%d/artifacts/content/%s", t.serverURL, id, path), nil)
	req.Header.Add("Authorization", t.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return fileContent, "", err
	}

	defer resp.Body.Close()
	fileContent, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return fileContent, "", err
	}
	return fileContent, resp.Header.Get("Content-Type"), nil
}
