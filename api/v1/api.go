package apiv1

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// TODO: Make the volume location configurable.
const (
	// onefsVolumesPath = "/ifs/data/docker/volumes"
	papiNameSpacePath = "namespace"
	papiVolumesPath   = "/ifs/volumes"
	papiExportsPath   = "platform/1/protocols/nfs/exports"
)

var debug bool

// HACK: this seems kinda fragile.  would probably be better if the caller kept track of the ID.
var exportID int

func init() {
	debug, _ = strconv.ParseBool(os.Getenv("GOISILON_DEBUG"))
}

type IsiVolume struct {
	Name         string `json:"name"`
	AttributeMap []struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"attrs"`
}

// Isi PAPI volume JSON structs
type VolumeName struct {
	Name string `json:"name"`
}

type getIsiVolumesResp struct {
	Children []*VolumeName `json:"children"`
}

// Isi PAPI Volume ACL JSON structs
type Ownership struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type AclRequest struct {
	Authoritative string     `json:"authoritative"`
	Action        string     `json:"action"`
	Owner         *Ownership `json:"owner"`
	Group         *Ownership `json:"group"`
}

// Isi PAPI volume attributes JSON struct
type getIsiVolumeAttributesResp struct {
	AttributeMap []struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"attrs"`
}

// Isi PAPI export path JSON struct
type ExportPathList struct {
	Paths []string `json:"paths"`
}

// Isi PAPI export ID JSON struct
type postIsiExportResp struct {
	ID int `json:"ID"`
}

// Isi PAPI export attributes JSON structs
type IsiExport struct {
	ID    int      `json:"ID"`
	Paths []string `json:"paths"`
}

type getIsiExportsResp struct {
	ExportList []*IsiExport `json:"exports"`
}

// Get a list of all docker volumes on the cluster
// Note: All Docker Volumes are being stored in a single directory on the cluster.
func (papi *PapiConnection) GetIsiVolumes() (resp *getIsiVolumesResp, err error) {
	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volumes/
	err = papi.query("GET", papi.nameSpacePath(),
		"", nil, nil, &resp)
	return resp, err
}

// Create a new docker volume on the cluster
// Note: A Docker Volume is just a directory on the cluster
func (papi *PapiConnection) CreateIsiVolume(name string) (resp *getIsiVolumesResp, err error) {
	// PAPI calls: PUT https://1.2.3.4:8080/namespace/path/to/volumes/volume_name
	//             x-isi-ifs-target-type: container
	//             x-isi-ifs-access-control: public_read_write
	//
	//             PUT https://1.2.3.4:8080/namespace/path/to/volumes/volume_name?acl
	//             {authoritative: "acl",
	//              action: "update",
	//              owner: {ID: "UID:65534", name: "nobody", type: "user"},
	//              group: {ID: "UID:65534", name: "nobody", type: "group"}
	//             }

	headers := map[string]string{"x-isi-ifs-target-type": "container", "x-isi-ifs-access-control": "public_read_write"}
	// TODO: This should be configurable
	var data = &AclRequest{
		"acl",
		"update",
		&Ownership{"UID:65534", "nobody", "user"},
		&Ownership{"GID:65534", "nobody", "group"},
	}

	// create the volume
	err = papi.queryWithHeaders("PUT", papi.nameSpacePath(), name, nil, headers, nil, &resp)
	if err != nil {
		return resp, err
	}

	// set the ownership of the volume
	err = papi.query("PUT", papi.nameSpacePath(), name, map[string]string{"acl": ""}, data, &resp)

	return resp, err
}

// Query the attributes of a docker volume on the cluster
// Note: A Docker Volume is just a directory on the cluster
func (papi *PapiConnection) GetIsiVolume(name string) (resp *getIsiVolumeAttributesResp, err error) {
	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volume/?metadata
	err = papi.query("GET", papi.nameSpacePath(), name, map[string]string{"metadata": ""}, nil, &resp)
	return resp, err
}

// Delete a docker volume from the cluster
// Note: A Docker Volume is just a directory on the cluster
func (papi *PapiConnection) DeleteIsiVolume(name string) (resp *getIsiVolumesResp, err error) {
	// PAPI call: DELETE https://1.2.3.4:8080/namespace/path/to/volumes/volume_name

	err = papi.queryWithHeaders("DELETE", papi.nameSpacePath(), name, nil, nil, nil, &resp)
	return resp, err
}

// Enable an NFS export on the cluster to access the docker volumes.  Return the path to the export
// so other processes can mount the docker volume directory
func (papi *PapiConnection) Export(path string) (err error) {
	// PAPI call: POST https://1.2.3.4:8080/platform/1/protocols/nfs/exports/
	//            Content-Type: application/json
	//            {paths: ["/path/to/volume"]}

	if path == "" {
		return errors.New("no path set")
	}

	var data = &ExportPathList{Paths: []string{path}}
	headers := map[string]string{"Content-Type": "application/json"}
	var resp *postIsiExportResp

	err = papi.queryWithHeaders("POST", papiExportsPath, "", nil, headers, data, &resp)

	if err != nil {
		return err
	}

	return nil
}

// Disable the NFS export on the cluster that points to the docker volumes directory.
func (papi *PapiConnection) Unexport(ID int) (err error) {
	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/protocols/nfs/exports/23

	if ID == 0 {
		return errors.New("no path ID set")
	}

	exportPath := fmt.Sprintf("%s/%d", papiExportsPath, ID)

	var resp postIsiExportResp
	err = papi.queryWithHeaders("DELETE", exportPath, "", nil, nil, nil, &resp)

	return err
}

func (papi *PapiConnection) nameSpacePath() string {
	return fmt.Sprintf("%s%s", papiNameSpacePath, papi.VolumePath)
}

func (papi *PapiConnection) exportsPath() string {
	return fmt.Sprintf("%s%s", papiExportsPath, papi.VolumePath)
}

// Get a list of all exports on the cluster
// TODO: This shouldn't be public, but I'm still researching how to do that while still being
// able to use the PapiConnection functions.
func (papi *PapiConnection) GetIsiExports() (resp *getIsiExportsResp, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/protocols/nfs/exports
	err = papi.query("GET", papiExportsPath, "", nil, nil, &resp)

	return resp, err
}
