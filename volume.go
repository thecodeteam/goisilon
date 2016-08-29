package goisilon

import (
	log "github.com/emccode/gournal"
	"golang.org/x/net/context"

	api "github.com/emccode/goisilon/api/v1"
)

type Volume *api.IsiVolume

//GetVolume returns a specific volume by name or ID
func (c *Client) GetVolume(
	ctx context.Context, id, name string) (Volume, error) {

	if id != "" {
		name = id
	}
	volume, err := api.GetIsiVolume(ctx, c.API, name)
	if err != nil {
		return nil, err
	}
	var isiVolume = &api.IsiVolume{Name: name, AttributeMap: volume.AttributeMap}
	return isiVolume, nil
}

//GetVolumes returns a list of volumes
func (c *Client) GetVolumes(ctx context.Context) ([]Volume, error) {

	volumes, err := api.GetIsiVolumes(ctx, c.API)
	if err != nil {
		return nil, err
	}
	var isiVolumes []Volume
	for _, volume := range volumes.Children {
		newVolume := &api.IsiVolume{Name: volume.Name}
		isiVolumes = append(isiVolumes, newVolume)
	}
	return isiVolumes, nil
}

//CreateVolume creates a volume
func (c *Client) CreateVolume(
	ctx context.Context, name string) (Volume, error) {

	_, err := api.CreateIsiVolume(ctx, c.API, name)
	if err != nil {
		return nil, err
	}

	var isiVolume = &api.IsiVolume{Name: name, AttributeMap: nil}
	return isiVolume, nil
}

//DeleteVolume deletes a volume
func (c *Client) DeleteVolume(
	ctx context.Context, name string) error {

	_, err := api.DeleteIsiVolume(ctx, c.API, name)
	return err
}

//CopyVolume creates a volume based on an existing volume
func (c *Client) CopyVolume(
	ctx context.Context, src, dest string) (Volume, error) {

	_, err := api.CopyIsiVolume(ctx, c.API, src, dest)
	if err != nil {
		return nil, err
	}

	return c.GetVolume(ctx, dest, dest)
}

//ExportVolume exports a volume
func (c *Client) ExportVolume(
	ctx context.Context, name string) (int, error) {

	return c.Export(ctx, name)
}

//UnexportVolume stops exporting a volume
func (c *Client) UnexportVolume(
	ctx context.Context, name string) error {

	return c.Unexport(ctx, name)
}

// GetVolumeExportMap returns a map that relates Volumes to their corresponding
// Exports. This function uses an Export's "clients" property to define the
// relationship. The flag "includeRootClients" can be set to "true" in order to
// also inspect the "root_clients" property of an Export when determining the
// Volume-to-Export relationship.
func (c *Client) GetVolumeExportMap(
	ctx context.Context,
	includeRootClients bool) (map[Volume]Export, error) {

	volumes, err := c.GetVolumes(ctx)
	if err != nil {
		return nil, err
	}
	exports, err := c.GetExports(ctx)
	if err != nil {
		return nil, err
	}

	volToExpMap := map[Volume]Export{}

	for _, v := range volumes {
		vp := c.API.VolumePath(v.Name)
		for _, e := range exports {
			if e.Clients == nil {
				continue
			}
			for _, p := range *e.Clients {
				if vp == p {
					if _, ok := volToExpMap[v]; ok {
						log.WithFields(map[string]interface{}{
							"volumeName": v.Name,
							"volumePath": vp,
						}).Info(ctx, "vol-ex client map already defined")
						break
					}
					volToExpMap[v] = e
				}
			}
			if !includeRootClients || e.RootClients == nil {
				continue
			}
			for _, p := range *e.RootClients {
				if vp == p {
					if _, ok := volToExpMap[v]; ok {
						log.WithFields(map[string]interface{}{
							"volumeName": v.Name,
							"volumePath": vp,
						}).Info(ctx, "vol-ex root client map already defined")
						break
					}
					volToExpMap[v] = e
				}
			}
		}
	}

	return volToExpMap, nil
}
