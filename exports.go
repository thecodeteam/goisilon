package goisilon

import papi "github.com/emccode/goisilon/api/v1"

import "strconv"

type ExportList []*papi.IsiExport
type Export *papi.IsiExport

func (c *Client) GetIsiExports() (ExportList, error) {
	exports, err := c.api.GetIsiExports()
	if err != nil {
		return nil, err
	}

	return exports.ExportList, nil
}

func (c *Client) GetIsiExport(id, path string) (Export, error) {
	exports, err := c.GetIsiExports()
	if err != nil {
		return nil, err
	}

	for _, export := range exports {
		if strconv.Itoa(export.ID) == id {
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

func (c *Client) Export(path string) error {
	if exported, err := c.isExported(path); !exported {
		return c.api.Export(path)
	} else if err != nil {
		return err
	}
	return nil
}

func (c *Client) Unexport(path string) error {
	export, err := c.GetIsiExport("", path)
	if err != nil {
		return err
	}

	if export == nil {
		return nil
	}

	return c.api.Unexport(export.ID)
}

func (c *Client) isExported(path string) (bool, error) {
	exportList, err := c.api.GetIsiExports()
	if err != nil {
		return false, err
	}

	for _, export := range exportList.ExportList {
		for _, existingPath := range export.Paths {
			if existingPath == path {
				return true, nil
			}
		}
	}
	return false, nil
}
