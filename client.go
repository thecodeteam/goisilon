package goisilon

import (
	"os"
//	"strconv"

	papi "goisilon/api/v1"
)

type Client1 struct {
	api *papi.PapiConnection
}

type Client Client1

func NewClient() (*Client, error) {
//	insecure, _ := strconv.ParseBool(os.Getenv("GOISILON_INSECURE"))
	return NewClientWithArgs(
		os.Getenv("GOISILON_ENDPOINT"),
		false,
		os.Getenv("GOISILON_USERNAME"),
		os.Getenv("GOISILON_PASSWORD"))
}

func NewClientWithArgs(
	endpoint string,
	insecure bool,
	username, password string) (*Client, error) {

	api, err := papi.New(endpoint, insecure, username, password)
	if err != nil {
		return nil, err
	}

	return &Client{api}, nil
}