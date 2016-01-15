package apiv1

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	papiNameSpacePath = "namespace"
	papiVolumesPath   = "/ifs/volumes"
	papiExportsPath   = "platform/1/protocols/nfs/exports"
	papiQuotaPath     = "platform/1/quota/quotas"
)

var debug bool

// HACK: this seems kinda fragile.  would probably be better if the caller kept track of the Id.
var exportId int

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
	Name string `json:"name"`
	Type string `json:"type"`
}

type AclRequest struct {
	Authoritative string     `json:"authoritative"`
	Action        string     `json:"action"`
	Owner         *Ownership `json:"owner"`
	Group         *Ownership `json:"group,omitempty"`
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
	Paths  []string `json:"paths"`
	MapAll struct {
		User   string   `json:"user"`
		Groups []string `json:"groups,omitempty"`
	} `json:"map_all"`
}

// Isi PAPI export clients JSON struct
type ExportClientList struct {
	Clients []string `json:"clients"`
}

// Isi PAPI export Id JSON struct
type postIsiExportResp struct {
	Id int `json:"id"`
}

// Isi PAPI export attributes JSON structs
type IsiExport struct {
	Id      int      `json:"id"`
	Paths   []string `json:"paths"`
	Clients []string `json:"clients"`
}

type getIsiExportsResp struct {
	ExportList []*IsiExport `json:"exports"`
}

type isiThresholds struct {
	Advisory             int64       `json:"advisory"`
	AdvisoryExceeded     bool        `json:"advisory_exceeded"`
	AdvisoryLastExceeded interface{} `json:"advisory_last_exceeded"`
	Hard                 int64       `json:"hard"`
	HardExceeded         bool        `json:"hard_exceeded"`
	HardLastExceeded     interface{} `json:"hard_last_exceeded"`
	Soft                 int64       `json:"soft"`
	SoftExceeded         bool        `json:"soft_exceeded"`
	SoftLastExceeded     interface{} `json:"soft_last_exceeded"`
}

type IsiQuota struct {
	Container                 bool          `json:"container"`
	Enforced                  bool          `json:"enforced"`
	Id                        string        `json:"id"`
	IncludeSnapshots          bool          `json:"include_snapshots"`
	Linked                    interface{}   `json:"linked"`
	Notifications             string        `json:"notifications"`
	Path                      string        `json:"path"`
	Persona                   interface{}   `json:"persona"`
	Ready                     bool          `json:"ready"`
	Thresholds                isiThresholds `json:"thresholds"`
	ThresholdsIncludeOverhead bool          `json:"thresholds_include_overhead"`
	Type                      string        `json:"type"`
	Usage                     struct {
		Inodes   int64 `json:"inodes"`
		Logical  int64 `json:"logical"`
		Physical int64 `json:"physical"`
	} `json:"usage"`
}

type isiThresholdsReq struct {
	Advisory interface{} `json:"advisory"`
	Hard     interface{} `json:"hard"`
	Soft     interface{} `json:"soft"`
}

type IsiQuotaReq struct {
	Enforced                  bool             `json:"enforced"`
	IncludeSnapshots          bool             `json:"include_snapshots"`
	Path                      string           `json:"path"`
	Thresholds                isiThresholdsReq `json:"thresholds"`
	ThresholdsIncludeOverhead bool             `json:"thresholds_include_overhead"`
	Type                      string           `json:"type"`
}

type IsiUpdateQuotaReq struct {
	Enforced                  bool             `json:"enforced"`
	Thresholds                isiThresholdsReq `json:"thresholds"`
	ThresholdsIncludeOverhead bool             `json:"thresholds_include_overhead"`
}

type isiQuotaListResp struct {
	Quotas []IsiQuota `json:"quotas"`
}

// Get the quota for a directory
func (papi *PapiConnection) GetIsiQuota(path string) (quota *IsiQuota, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/quota/quotas
	// This will list out all quotas on the cluster

	var quotaResp isiQuotaListResp
	err = papi.query("GET", papiQuotaPath, "", nil, nil, &quotaResp)
	if err != nil {
		return nil, err
	}

	// find the specific quota we are looking for
	for _, quota := range quotaResp.Quotas {
		if quota.Path == path {
			return &quota, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Quota not found: %s", path))
}

// TODO: Add a means to set/update more than just the hard threshold

// Set the hard threshold of a quota for a directory
func (papi *PapiConnection) SetIsiQuotaHardThreshold(path string, size int64) (err error) {
	// PAPI call: POST https://1.2.3.4:8080/platform/1/quota/quotas
	//             { "enforced" : true,
	//               "include_snapshots" : false,
	//               "path" : "/ifs/volumes/volume_name",
	//               "thresholds_include_overhead" : false,
	//               "type" : "directory",
	//               "thresholds" : { "advisory" : null,
	//                                "hard" : 1234567890,
	//                                "soft" : null
	//                              }
	//             }
	var data = &IsiQuotaReq{
		Enforced:         true,
		IncludeSnapshots: false,
		Path:             path,
		ThresholdsIncludeOverhead: false,
		Type:       "directory",
		Thresholds: isiThresholdsReq{Advisory: nil, Hard: size, Soft: nil},
	}

	var quotaResp IsiQuota
	err = papi.query("POST", papiQuotaPath, "", nil, data, &quotaResp)
	return err
}

// Update the hard threshold of a quota for a directory
func (papi *PapiConnection) UpdateIsiQuotaHardThreshold(path string, size int64) (err error) {
	// PAPI call: PUT https://1.2.3.4:8080/platform/1/quota/quotas/Id
	//             { "enforced" : true,
	//               "thresholds_include_overhead" : false,
	//               "thresholds" : { "advisory" : null,
	//                                "hard" : 1234567890,
	//                                "soft" : null
	//                              }
	//             }
	var data = &IsiUpdateQuotaReq{
		Enforced:                  true,
		ThresholdsIncludeOverhead: false,
		Thresholds:                isiThresholdsReq{Advisory: nil, Hard: size, Soft: nil},
	}

	quota, err := papi.GetIsiQuota(path)
	if err != nil {
		return err
	}

	var quotaResp IsiQuota
	err = papi.query("PUT", papiQuotaPath, quota.Id, nil, data, &quotaResp)
	return err
}

// Delete the quota for a directory
func (papi *PapiConnection) DeleteIsiQuota(path string) (err error) {
	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/quota/quotas?path=/path/to/volume
	// This will remove a the quota on a volume

	var quotaResp isiQuotaListResp
	err = papi.query("DELETE", papiQuotaPath, "", map[string]string{"path": path}, nil, &quotaResp)

	return err
}

// Get a list of all docker volumes on the cluster
// Note: All Docker Volumes are being stored in a single directory on the cluster.
func (papi *PapiConnection) GetIsiVolumes() (resp *getIsiVolumesResp, err error) {
	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volumes/
	err = papi.query("GET", papi.nameSpacePath(), "", nil, nil, &resp)
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
	//              owner: {name: "username", type: "user"},
	//              group: {name: "groupname", type: "group"}
	//             }

	headers := map[string]string{"x-isi-ifs-target-type": "container", "x-isi-ifs-access-control": "public_read_write"}
	var data = &AclRequest{
		"acl",
		"update",
		&Ownership{papi.username, "user"},
		nil,
	}
	if papi.group != "" {
		data.Group = &Ownership{papi.group, "group"}
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
	data.MapAll.User = papi.username
	if papi.group != "" {
		data.MapAll.Groups = append(data.MapAll.Groups, papi.group)
	}
	headers := map[string]string{"Content-Type": "application/json"}
	var resp *postIsiExportResp

	err = papi.queryWithHeaders("POST", papiExportsPath, "", nil, headers, data, &resp)

	if err != nil {
		return err
	}

	return nil
}

// Limit access to an NFS export on the cluster to a specific client address.
func (papi *PapiConnection) SetExportClients(Id int, clients []string) (err error) {
	// PAPI call: PUT https://1.2.3.4:8080/platform/1/protocols/nfs/exports/Id
	//            Content-Type: application/json
	//            {clients: ["client_ip_address"]}

	var data = &ExportClientList{Clients: clients}
	headers := map[string]string{"Content-Type": "application/json"}
	var resp *postIsiExportResp

	err = papi.queryWithHeaders("PUT", papiExportsPath, strconv.Itoa(Id), nil, headers, data, &resp)

	return err
}

// Disable the NFS export on the cluster that points to the docker volumes directory.
func (papi *PapiConnection) Unexport(Id int) (err error) {
	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/protocols/nfs/exports/23

	if Id == 0 {
		return errors.New("no path Id set")
	}

	exportPath := fmt.Sprintf("%s/%d", papiExportsPath, Id)

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
