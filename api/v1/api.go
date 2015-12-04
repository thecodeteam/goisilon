package apiv1

import (
	"os"
	"strconv"
	"fmt"
)

// TODO: Make the volume location configurable.
const (
	onefsVolumesPath = "/ifs/data/docker/volumes"
	papiVolumesPath = "namespace/ifs/data/docker/volumes"
	papiExportsPath = "platform/1/protocols/nfs/exports"
)

var debug bool

// HACK: this seems kinda fragile.  would probably be better if the caller kept track of the ID.
var exportID int

func init() {
	debug, _ = strconv.ParseBool(os.Getenv("GOISILON_DEBUG"))
}

type IsiVolume struct {
	Name string `json:"name"`
	AttributeMap []struct {
		Name string `json:"name"`
		Value interface{} `json:"value"`
	}`json:"attrs"`	
}

// Isi PAPI volume JSON structs
type volumeName struct {
	Name string `json:"name"`
}

type getIsiVolumesResp struct {
	Children []volumeName `json:"children"`
}

// Isi PAPI Volume ACL JSON structs
type ownership struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type aclRequest struct {
	Authoritative string    `json:"authoritative"`
	Action        string    `json:"action"`
	Owner         ownership `json:"owner"`
	Group         ownership `json:"group"`
}

// Isi PAPI volume attributes JSON struct
type getIsiVolumeAttributesResp struct {
	AttributeMap []struct {
		Name string `json:"name"`
		Value interface{} `json:"value"`
	}`json:"attrs"`
}

// Isi PAPI export path JSON struct
type exportPathList struct {
	Paths []string `json:"paths"`
}

// Isi PAPI export id JSON struct
type postIsiExportResp struct {
	Id int `json:"id"`
}

// Isi PAPI export attributes JSON structs
type isiExport struct {
	Id int `json:"id"`
	Paths []string `json:"paths"`
}

type getIsiExportsResp struct {
	ExportList []isiExport `json:"exports"`
}

// Get a list of all docker volumes on the cluster
// Note: All Docker Volumes are being stored in a single directory on the cluster.
func (papi *PapiConnection) GetIsiVolumes() (resp *getIsiVolumesResp, err error) {
	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volumes/
	err = papi.query("GET", papiVolumesPath, "", nil, nil, &resp)
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
	//              owner: {id: "UID:65534", name: "nobody", type: "user"}, 
	//              group: {id: "UID:65534", name: "nobody", type: "group"}
	//             }
	
	headers := map[string]string{"x-isi-ifs-target-type": "container", "x-isi-ifs-access-control": "public_read_write"}
	// TODO: This should be configurable
	var data = aclRequest{
		"acl",
		"update",
		ownership{"UID:65534", "nobody", "user"},
		ownership{"GID:65534", "nobody", "group"},
	}

    // create the volume
	err = papi.queryWithHeaders("PUT", papiVolumesPath, name, nil, headers, nil, &resp)
	if err != nil {
		return resp, err
	}
	
	// set the ownership of the volume
	err = papi.query("PUT", papiVolumesPath, name, map[string]string{"acl": ""}, data, &resp)
	
	return resp, err
}

// Query the attributes of a docker volume on the cluster
// Note: A Docker Volume is just a directory on the cluster
func (papi *PapiConnection) GetIsiVolume(name string) (resp *getIsiVolumeAttributesResp, err error) {
	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volume/?metadata
	err = papi.query("GET", papiVolumesPath, name, map[string]string{"metadata": ""}, nil, &resp)
	return resp, err
}

// Delete a docker volume from the cluster
// Note: A Docker Volume is just a directory on the cluster
func (papi *PapiConnection) DeleteIsiVolume(name string) (resp *getIsiVolumesResp, err error) {
	// PAPI call: DELETE https://1.2.3.4:8080/namespace/path/to/volumes/volume_name

	err = papi.queryWithHeaders("DELETE", papiVolumesPath, name, nil, nil, nil, &resp)
	return resp, err
}

// Enable an NFS export on the cluster to access the docker volumes.  Return the path to the export
// so other processes can mount the docker volume directory
func (papi *PapiConnection) Attach() (path string, err error) {	
	// PAPI call: POST https://1.2.3.4:8080/platform/1/protocols/nfs/exports/ 
	//            Content-Type: application/json
	//            {paths: ["/path/to/volume"]}
	
	// see if the docker volumes export already exists
	exportList, err := papi.GetIsiExports()
	for _, export := range exportList.ExportList {
		for _, path := range export.Paths {
			if path == onefsVolumesPath {
				// the export already exists, grab it's id and return it's path
				exportID = export.Id
				return path, nil
			}
		}
	}
	
	// the docker volumes export doesn't exist yet, create it.
	var data = exportPathList{Paths: []string{onefsVolumesPath}}
	headers := map[string]string{"Content-Type": "application/json"}
	var resp *postIsiExportResp
	
	err = papi.queryWithHeaders("POST", papiExportsPath, "", nil, headers, data, &resp)

	if err != nil {
		return "None", err
	}
	exportID = resp.Id
	
	return data.Paths[0], err
	
}

// Disable the NFS export on the cluster that points to the docker volumes directory.
func (papi *PapiConnection) Detach() (err error) {
	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/protocols/nfs/exports/23

	exportPath := fmt.Sprintf("%s%d", papiExportsPath, exportID)
	var resp postIsiExportResp 
	
	err = papi.queryWithHeaders("DELETE", exportPath, "", nil, nil, nil, &resp)

	return err	
}

// Get a list of all exports on the cluster
// TODO: This shouldn't be public, but I'm still researching how to do that while still being
// able to use the PapiConnection functions.
func (papi *PapiConnection) GetIsiExports() (resp *getIsiExportsResp, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/protocols/nfs/exports
	err = papi.query("GET", papiExportsPath, "", nil, nil, &resp)

	return resp, err
}

