package goisilon

import papi "github.com/emccode/goisilon/api/v1"

import "strconv"

type ExportList []*papi.IsiExport
type Export *papi.IsiExport

// Gets a list of all exports on the cluster
func (c *Client) GetIsiExports() (ExportList, error) {
	exports, err := c.api.GetIsiExports()
	if err != nil {
		return nil, err
	}

	return exports.ExportList, nil
}

// Gets a specific export (by volume id or name) on a cluster
func (c *Client) GetIsiExport(id, name string) (Export, error) {
	exports, err := c.GetIsiExports()
	if err != nil {
		return nil, err
	}

	path := c.Path(name)
	for _, export := range exports {
		if strconv.Itoa(export.Id) == id {
			return export, nil
		} else {
			for _, existingPath := range export.Paths {
				if existingPath == path {
					return export, nil
				}
			}
		}
	}
	return nil, nil
}

// Export the volume with a given name on the cluster
func (c *Client) Export(name string) error {
	if exported, err := c.isExported(name); !exported {
		return c.api.Export(c.Path(name))
	} else if err != nil {
		return err
	}
	return nil
}

// Get a list of all clients allowed to connect to a given volume
func (c *Client) GetExportClients(name string) ([]string, error) {
	export, err := c.GetIsiExport("", name)
	if err != nil {
		return nil, err
	}

	if export == nil {
		return nil, nil
	}
	return export.Clients, nil
}

// Set the list of clients allowed to connect to a given volume
func (c *Client) SetExportClients(name string, client []string) error {
	export, err := c.GetIsiExport("", name)
	if err != nil {
		return err
	}

	if export == nil {
		return nil
	}

	return c.api.SetExportClients(export.Id, client)
}

// Clear the list of all clients that can connect to a given volume.  This essentially
// makes the volume accessible by all.
func (c *Client) ClearExportClients(name string) error {
	return c.SetExportClients(name, make([]string, 0))
}

// Stop exporting a given volume from the cluster
func (c *Client) Unexport(name string) error {
	export, err := c.GetIsiExport("", name)
	if err != nil {
		return err
	}

	if export == nil {
		return nil
	}

	return c.api.Unexport(export.Id)
}

// Check if a volume is currently being exported
func (c *Client) isExported(name string) (bool, error) {
	exportList, err := c.api.GetIsiExports()
	if err != nil {
		return false, err
	}

	for _, export := range exportList.ExportList {
		for _, existingPath := range export.Paths {
			if existingPath == c.Path(name) {
				return true, nil
			}
		}
	}
	return false, nil
}
