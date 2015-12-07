package goisilon

import (
	"fmt"

	papi "github.com/emccode/goisilon/api/v1"
)

type Volume *papi.IsiVolume
type VolumeExport struct {
	Volume     Volume
	ExportPath string
}

//type NewVolumeOptions papi.PostVolumesReq
//type NewVolumeResult *papi.PostVolumesResp

//GetVolume returns a specific volume by name or ID
func (c *Client) GetVolume(id string, name string) (Volume, error) {
	if id != "" {
		name = id
	}
	volume, err := c.api.GetIsiVolume(name)
	if err != nil {
		return nil, err
	}
	var isiVolume = &papi.IsiVolume{Name: name, AttributeMap: volume.AttributeMap}
	return isiVolume, nil
}

//GetVolumes returns a list of volumes
func (c *Client) GetVolumes() ([]Volume, error) {
	volumes, err := c.api.GetIsiVolumes()
	if err != nil {
		return nil, err
	}
	var isiVolumes []Volume
	for _, volume := range volumes.Children {
		newVolume := &papi.IsiVolume{Name: volume.Name}
		isiVolumes = append(isiVolumes, newVolume)
	}
	return isiVolumes, nil
}

//CreateVolume creates a volume
func (c *Client) CreateVolume(name string) (Volume, error) {
	_, err := c.api.CreateIsiVolume(name)
	if err != nil {
		return nil, err
	}

	var isiVolume = &papi.IsiVolume{Name: name, AttributeMap: nil}
	return isiVolume, nil
}

//DeleteVolume deletes a volume
func (c *Client) DeleteVolume(name string) error {
	_, err := c.api.DeleteIsiVolume(name)
	return err
}

func (c *Client) Path(name string) string {
	return fmt.Sprintf("%s/%s", c.api.VolumePath, name)
}

func (c *Client) ExportVolume(name string) error {
	return c.Export(c.Path(name))
}

//DeleteVolume deletes a volume
func (c *Client) UnexportVolume(name string) error {
	return c.Unexport(c.Path(name))
}
func (c *Client) GetVolumeExports() ([]*VolumeExport, error) {
	exports, err := c.GetIsiExports()
	if err != nil {
		return nil, err
	}

	exportPaths := make(map[string]bool)
	for _, export := range exports {
		for _, path := range export.Paths {
			exportPaths[path] = true
		}
	}

	volumes, err := c.GetVolumes()
	if err != nil {
		return nil, err
	}

	var volumeExports []*VolumeExport
	for _, volume := range volumes {
		if _, ok := exportPaths[c.Path(volume.Name)]; ok {
			volumeExports = append(volumeExports, &VolumeExport{
				Volume:     volume,
				ExportPath: c.Path(volume.Name),
			})
		}
	}

	return volumeExports, nil
}
