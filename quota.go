package goisilon

import (
	papi "github.com/emccode/goisilon/api/v1"
)

type Quota *papi.IsiQuota

// GetQuota returns a specific quota by path
func (c *Client) GetQuota(name string) (Quota, error) {
	quota, err := c.api.GetIsiQuota(c.Path(name))
	if err != nil {
		return nil, err
	}

	return quota, nil
}

// TODO: Add a means to set/update more fields of the quota

// SetQuota sets the max size (hard threshold) of a quota for a volume
func (c *Client) SetQuotaSize(name string, size int64) error {
	err := c.api.SetIsiQuotaHardThreshold(c.Path(name), size)
	return err
}

// UpdateQuota modifies the max size (hard threshold) of a quota for a volume
func (c *Client) UpdateQuotaSize(name string, size int64) error {
	err := c.api.UpdateIsiQuotaHardThreshold(c.Path(name), size)
	return err
}

// ClearQuota removes the quota from a volume
func (c *Client) ClearQuota(name string) error {
	err := c.api.DeleteIsiQuota(c.Path(name))
	return err
}
