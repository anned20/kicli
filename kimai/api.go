package kimai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type API struct {
	client   *http.Client
	baseURL  string
	username string
	token    string
}

func NewAPI(baseURL string, username string, token string) *API {
	// Strip trailing slash from baseURL
	if baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	logrus.WithFields(logrus.Fields{
		"baseURL":  baseURL,
		"username": username,
		"token":    token,
	}).Debug("Creating API")

	return &API{
		client:   &http.Client{},
		baseURL:  baseURL,
		username: username,
		token:    token,
	}
}

func (a *API) get(url string) (*http.Response, error) {
	url = fmt.Sprintf("%s%s", a.baseURL, url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	return a.do(req)
}

func (a *API) post(url string, body interface{}) (*http.Response, error) {
	url = fmt.Sprintf("%s%s", a.baseURL, url)

	jsonData, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	return a.do(req)
}

func (a *API) patch(url string, body interface{}) (*http.Response, error) {
	url = fmt.Sprintf("%s%s", a.baseURL, url)

	var jsonData []byte

	if body != nil {
		data, err := json.Marshal(body)

		if err != nil {
			return nil, err
		}

		jsonData = data
	}

	logrus.WithFields(logrus.Fields{
		"url":  url,
		"body": string(jsonData),
	}).Debug("Patching")

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	return a.do(req)
}

func (a *API) do(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-AUTH-USER", a.username)
	req.Header.Add("X-AUTH-TOKEN", a.token)
	req.Header.Add("Accept", "application/json")

	logrus.WithFields(logrus.Fields{
		"url":     req.URL.String(),
		"headers": req.Header,
		"body":    req.Body,
	}).Debug("Doing request")

	response, err := a.client.Do(req)

	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		// Get the error message from the response body
		errorMessage, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s", response.Status, errorMessage)
	}

	return response, nil
}
