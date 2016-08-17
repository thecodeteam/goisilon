package goisilon

import (
	"os"
	"strconv"

	papi "github.com/emccode/goisilon/api/v1"
)

func NamespacePath(p string) {
	papi.NamespacePath = p
}
func VolumesPath(p string) {
	papi.VolumesPath = p
}
func ExportsPath(p string) {
	papi.ExportsPath = p
}
func QuotaPath(p string) {
	papi.QuotaPath = p
}
func SnapshotsPath(p string) {
	papi.SnapshotsPath = p
}
func VolumeSnapshotsPath(p string) {
	papi.VolumeSnapshotsPath = p
}

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
