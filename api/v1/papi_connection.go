// Note: The bulk of this was taken from the goextremio package.
package apiv1

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type PapiConnection struct {
	endpoint   string
	insecure   bool
	username   string
	group      string
	password   string
	httpClient *http.Client
	VolumePath string
}

// Isi PAPI error JSON structs
type PapiError struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Error struct {
	StatusCode int
	Err        []PapiError `json:"errors"`
}

// Create a new HTTP connection
func New(endpoint string, insecure bool, username, group, password, volumePath string) (*PapiConnection, error) {
	if endpoint == "" || username == "" || password == "" {
		return nil, errors.New("Missing endpoint, username, or password")
	}

	if volumePath == "" {
		volumePath = papiVolumesPath
	} else if volumePath[0] == '/' {
		volumePath = fmt.Sprintf("%s%s", papiVolumesPath, volumePath)
	} else {
		volumePath = fmt.Sprintf("%s/%s", papiVolumesPath, volumePath)
	}

	var client *http.Client
	if insecure {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecure,
				},
			},
		}
	} else {
		client = &http.Client{}
	}

	return &PapiConnection{endpoint, insecure, username, group, password, client, volumePath}, nil
}

func multimap(p map[string]string) url.Values {
	q := make(url.Values, len(p))
	for k, v := range p {
		q[k] = []string{v}
	}
	return q
}

// Extract the error string from a received error message
func (err *Error) Error() string {
	// I've only seen PAPI return a single error, but, technically, it can be a list
	return err.Err[0].Message
}

// Parse a PAPI error message sent by the cluster
func buildError(r *http.Response) error {
	jsonError := Error{}
	json.NewDecoder(r.Body).Decode(&jsonError)

	jsonError.StatusCode = r.StatusCode
	// I've only seen PAPI return a single error, but, technically, it can be a list
	if jsonError.Err[0].Message == "" {
		jsonError.Err[0].Message = r.Status
	}

	return &jsonError
}

// Send an HTTP query to the cluster
func (xms *PapiConnection) query(method string, path string, id string,
	params map[string]string, body interface{}, resp interface{}) error {

	return xms.queryWithHeaders(method, path, id, params, nil, body, resp)
}

// Send an HTTP query that includes headers to the cluster
func (xms *PapiConnection) queryWithHeaders(method string, path string, id string,
	params map[string]string, headers map[string]string, body interface{},
	resp interface{}) error {

	// build the URI
	endpoint := fmt.Sprintf("%s/%s", xms.endpoint, path)
	if id != "" {
		endpoint = fmt.Sprintf("%s/%s", endpoint, id)
	}

	// add parameters to the URI
	encodedParams := multimap(params).Encode()
	if encodedParams != "" {
		endpoint = fmt.Sprintf("%s?%s", endpoint, encodedParams)
	}

	// marshal the message body (assumes json format)
	var byteBuffer bytes.Buffer
	if body != nil {
		var bodyBytes []byte
		bodyBytes, _ = json.Marshal(body)
		byteBuffer.Write(bodyBytes)
	}

	req, err := http.NewRequest(method, endpoint, &byteBuffer)
	if err != nil {
		return err
	}

	// add headers to the request
	if headers != nil {
		for header, value := range headers {
			req.Header.Add(header, value)
		}
	}

	// set the username and password
	req.SetBasicAuth(xms.username, xms.password)

	// send the request
	r, err := xms.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// parse the response
	switch {
	case resp == nil:
		return nil
	case r.StatusCode >= 200 && r.StatusCode <= 299:
		decoder := json.NewDecoder(r.Body)
		if decoder.More() {
			err = decoder.Decode(resp)
			return err
		}
		return nil
	default:
		return buildError(r)
	}

}
