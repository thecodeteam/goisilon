package goisilon

import (
	"fmt"
	"testing"
)

func init() {
	testClient()
}

// Test both GetQuota() and SetQuota()
func TestGetSetQuota(*testing.T) {
	volumeName := "test_get_set_quota"
	quotaSize := int64(12345)

	// Setup the test
	_, err := client.CreateVolume(volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(volumeName)
	defer client.ClearQuota(volumeName)

	// Make sure there is no quota yet
	quota, err := client.GetQuota(volumeName)
	if quota != nil {
		panic(fmt.Sprintf("Quota should be nil: %v", quota))
	}
	if err == nil {
		panic(fmt.Sprintf("GetQuota should return an error when there isn't a quota."))
	}

	// Set the quota
	err = client.SetQuotaSize(volumeName, quotaSize)
	if err != nil {
		panic(err)
	}

	// Make sure the quota was set
	quota, err = client.GetQuota(volumeName)
	if err != nil {
		panic(err)
	}
	if quota == nil {
		panic("Quota should not be nil")
	}
	if quota.Thresholds.Hard != quotaSize {
		panic(fmt.Sprintf("Unexpected new quota.  Expected: %d Actual: %d", quotaSize, quota.Thresholds.Hard))
	}

}

// Test UpdateQuota()
func TestUpdateQuota(*testing.T) {
	volumeName := "test_update_quota"
	quotaSize := int64(12345)
	updatedQuotaSize := int64(22345000)

	// Setup the test
	_, err := client.CreateVolume(volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(volumeName)
	defer client.ClearQuota(volumeName)
	// Set the quota
	err = client.SetQuotaSize(volumeName, quotaSize)
	if err != nil {
		panic(err)
	}
	// Make sure the quota is initialized
	quota, err := client.GetQuota(volumeName)
	if err != nil {
		panic(err)
	}
	if quota == nil {
		panic(fmt.Sprintf("Quota should not be nil: %v", quota))
	}
	if quota.Thresholds.Hard != quotaSize {
		panic(fmt.Sprintf("Initial quota not set properly.  Expected: %d Actual: %d", quotaSize, quota.Thresholds.Hard))
	}

	// Update the quota
	err = client.UpdateQuotaSize(volumeName, updatedQuotaSize)
	if err != nil {
		panic(err)
	}

	// Make sure the quota is updated
	quota, err = client.GetQuota(volumeName)
	if err != nil {
		panic(err)
	}
	if quota == nil {
		panic(fmt.Sprintf("Updated quota should not be nil: %v", quota))
	}
	if quota.Thresholds.Hard != updatedQuotaSize {
		panic(fmt.Sprintf("Updated quota not set properly.  Expected: %d Actual: %d", updatedQuotaSize, quota.Thresholds.Hard))
	}

}

// Test ClearQuota()
func TestClearQuota(*testing.T) {
	volumeName := "test_clear_quota"
	quotaSize := int64(12345)

	// Setup the test
	_, err := client.CreateVolume(volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(volumeName)
	defer client.ClearQuota(volumeName)
	// Set the quota
	err = client.SetQuotaSize(volumeName, quotaSize)
	if err != nil {
		panic(err)
	}
	// Make sure the quota is initialized
	quota, err := client.GetQuota(volumeName)
	if err != nil {
		panic(err)
	}
	if quota == nil {
		panic(fmt.Sprintf("Quota should not be nil: %v", quota))
	}
	if quota.Thresholds.Hard != quotaSize {
		panic(fmt.Sprintf("Initial quota not set properly.  Expected: %d Actual: %d", quotaSize, quota.Thresholds.Hard))
	}

	// Update the quota
	err = client.ClearQuota(volumeName)
	if err != nil {
		panic(err)
	}

	// Make sure the quota is gone
	quota, err = client.GetQuota(volumeName)
	if err == nil {
		panic(fmt.Sprintf("Attempting to get a cleared quota should return an error but returned nil"))
	}
	if quota != nil {
		panic(fmt.Sprintf("Cleared quota should be nil: %v", quota))
	}

}
