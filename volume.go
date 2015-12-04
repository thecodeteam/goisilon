package goisilon

import papi "goisilon/api/v1"

type Volume *papi.IsiVolume
//type NewVolumeOptions papi.PostVolumesReq
//type NewVolumeResult *papi.PostVolumesResp

//GetVolume returns a specific volume by name or ID
func (c *Client) GetVolume(id string, name string) (Volume, error) {
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
