package goisilon

import (
	"os"
	"strconv"

	papi "github.com/emccode/goisilon/api/v1"
)

type Client1 struct {
	api *papi.PapiConnection
}

type Client Client1

func NewClient() (*Client, error) {
	insecure, _ := strconv.ParseBool(os.Getenv("GOISILON_INSECURE"))
	return NewClientWithArgs(
		os.Getenv("GOISILON_ENDPOINT"),
		insecure,
		os.Getenv("GOISILON_USERNAME"),
		os.Getenv("GOISILON_GROUP"),
		os.Getenv("GOISILON_PASSWORD"),
		os.Getenv("GOISILON_VOLUMEPATH"))
}

func NewClientWithArgs(
	endpoint string,
	insecure bool,
	username, group, password, volumePath string) (*Client, error) {

	api, err := papi.New(endpoint, insecure, username, group, password, volumePath)
	if err != nil {
		return nil, err
	}

	return &Client{api}, nil
}
