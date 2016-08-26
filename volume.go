package goisilon

import (
	"golang.org/x/net/context"

	api "github.com/emccode/goisilon/api/v1"
)

type Volume *api.IsiVolume
type VolumeExport struct {
	Volume     Volume
	ExportPath string
	Clients    []string
}

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

//GetVolumeExports return a list of volume exports
func (c *Client) GetVolumeExports(
	ctx context.Context) ([]*VolumeExport, error) {

	exports, err := c.GetExports(ctx)
	if err != nil {
		return nil, err
	}

	exportPaths := make(map[string]bool)
	exportClients := make(map[string]([]string))
	for _, export := range exports {
		for _, path := range *export.Paths {
			exportPaths[path] = true
			exportClients[path] = *export.Clients
		}
	}

	volumes, err := c.GetVolumes(ctx)
	if err != nil {
		return nil, err
	}

	var volumeExports []*VolumeExport
	for _, volume := range volumes {
		if _, ok := exportPaths[c.API.VolumePath(volume.Name)]; ok {
			volumeExports = append(volumeExports, &VolumeExport{
				Volume:     volume,
				ExportPath: c.API.VolumePath(volume.Name),
				Clients:    exportClients[c.API.VolumePath(volume.Name)],
			})
		}
	}

	return volumeExports, nil
}
