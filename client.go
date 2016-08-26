package goisilon

import (
	"os"
	"strconv"

	"golang.org/x/net/context"

	"github.com/emccode/goisilon/api"
)

// Client is an Isilon client.
type Client struct {

	// API is the underlying OneFS API client.
	API api.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	insecure, _ := strconv.ParseBool(os.Getenv("GOISILON_INSECURE"))
	return NewClientWithArgs(
		ctx,
		os.Getenv("GOISILON_ENDPOINT"),
		insecure,
		os.Getenv("GOISILON_USERNAME"),
		os.Getenv("GOISILON_GROUP"),
		os.Getenv("GOISILON_PASSWORD"),
		os.Getenv("GOISILON_VOLUMEPATH"))
}

func NewClientWithArgs(
	ctx context.Context,
	endpoint string,
	insecure bool,
	user, group, pass, volumePath string) (*Client, error) {

	client, err := api.NewWithVolumesPath(
		ctx, endpoint, user, pass, group, insecure, volumePath)
	if err != nil {
		return nil, err
	}

	return &Client{client}, err
}
