package goisilon

import (
	"fmt"
	"sort"
	"testing"
)

func init() {
	testClient()
}

// Test GetIsiExports()
func TestGetIsiExports(*testing.T) {
	volumeName1 := "test_get_exports1"
	volumeName2 := "test_get_exports2"
	volumeName3 := "test_get_exports3"
	volumePath1 := client.Path(volumeName1)
	volumePath2 := client.Path(volumeName2)
	volumePath3 := client.Path(volumeName3)

	// Identify all exports currently on the cluster
	exportMap := make(map[int]string)
	exports, err := client.GetIsiExports()
	if err != nil {
		panic(err)
	}
	for _, export := range exports {
		exportMap[export.Id] = export.Paths[0]
	}
	initialExportCount := len(exports)
	// Add the test exports
	_, err = client.CreateVolume(volumeName1)
	if err != nil {
		panic(err)
	}
	_, err = client.CreateVolume(volumeName2)
	if err != nil {
		panic(err)
	}
	_, err = client.CreateVolume(volumeName3)
	if err != nil {
		panic(err)
	}
	err = client.Export(volumeName1)
	if err != nil {
		panic(err)
	}
	err = client.Export(volumeName2)
	if err != nil {
		panic(err)
	}
	err = client.Export(volumeName3)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.Unexport(volumeName1)
	defer client.Unexport(volumeName2)
	defer client.Unexport(volumeName3)
	defer client.DeleteVolume(volumeName1)
	defer client.DeleteVolume(volumeName2)
	defer client.DeleteVolume(volumeName3)

	// Get the updated export list
	exports, err = client.GetIsiExports()
	if err != nil {
		panic(err)
	}

	// Verify that the new exports are there as well as all the old exports.
	if len(exports) != initialExportCount+3 {
		panic(fmt.Sprintf("Incorrect number of exports.  Expected: %d Actual: %d", initialExportCount+3, len(exports)))
	}
	// Remove the original exports and add the new ones.  In the end, we should only have the
	// exports we just created and nothing more.
	for _, export := range exports {
		if _, found := exportMap[export.Id]; found == true {
			// this export was exported prior to the test start
			delete(exportMap, export.Id)
		} else {
			// this export is new
			exportMap[export.Id] = export.Paths[0]
		}
	}
	if len(exportMap) != 3 {
		panic(fmt.Sprintf("Incorrect number of new exports.  Exptected: 3 Actual: %d", len(exportMap)))
	}
	volumeBitmap := 0
	for _, path := range exportMap {
		if path == volumePath1 {
			volumeBitmap += 1
		} else if path == volumePath2 {
			volumeBitmap += 2
		} else if path == volumePath3 {
			volumeBitmap += 4
		}
	}
	if volumeBitmap != 7 {
		panic(fmt.Sprintf("Incorrect new exports: %v", exportMap))
	}
}

// Test Export()
func TestCreateExport(*testing.T) {
	volumeName := "test_create_export"
	volumePath := client.Path(volumeName)

	// setup the test
	_, err := client.CreateVolume(volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.Unexport(volumeName)
	defer client.DeleteVolume(volumeName)
	// verify the volume isn't already exported
	export, err := client.GetIsiExport(volumeName, volumeName)
	if err != nil {
		panic(fmt.Sprintf("Unable to query volume (%s) to be exported.  Error: %v", volumeName, err))
	}
	if export != nil {
		panic(fmt.Sprintf("Volume is already exported (%s)", volumeName))
	}

	// export the volume
	err = client.Export(volumeName)
	if err != nil {
		panic(fmt.Sprintf("Error exporting volume (%s).  Error: %v", volumeName, err))
	}

	// verify the volume has been exported
	export, err = client.GetIsiExport(volumeName, volumeName)
	if err != nil {
		panic(fmt.Sprintf("Unable to query volume (%s) after exporting.  Error: %v", volumeName, err))
	}
	if export == nil {
		panic(fmt.Sprintf("Volume was not exported (%s)", volumeName))
	}
	found := false
	for _, path := range export.Paths {
		if path == volumePath {
			found = true
			break
		}
	}
	if found == false {
		panic(fmt.Sprintf("Export does not include volume path. Expected: %s Actual: %v", volumePath, export))
	}
}

// Test Unexport()
func TestRemoveExport(*testing.T) {
	volumeName := "test_unexport_volume"

	// initialize the export
	_, err := client.CreateVolume(volumeName)
	if err != nil {
		panic(err)
	}
	err = client.Export(volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(volumeName)

	// verify the volume is exported
	export, err := client.GetIsiExport(volumeName, volumeName)
	if err != nil {
		panic(fmt.Sprintf("Unable to query volume (%s) to be exported.  Error: %v", volumeName, err))
	}
	if export == nil {
		panic(fmt.Sprintf("Volume wasn't exported (%s)", volumeName))
	}

	// Unexport the volume
	err = client.Unexport(volumeName)
	if err != nil {
		panic(fmt.Sprintf("Error Unexporting volume (%s).  Error: %v", volumeName, err))
	}

	// verify the volume is no longer exported
	export, err = client.GetIsiExport(volumeName, volumeName)
	if err != nil {
		panic(fmt.Sprintf("Unable to query volume (%s) after exporting.  Error: %v", volumeName, err))
	}
	if export != nil {
		panic(fmt.Sprintf("Volume is still exported (%s)", volumeName))
	}
}

// Test GetExportClients()
func TestGetExportClients(*testing.T) {
	volumeName := "test_get_export_clients"
	clientList := []string{"1.2.3.4", "1.2.3.5"}

	// initialize the export
	_, err := client.CreateVolume(volumeName)
	if err != nil {
		panic(err)
	}
	err = client.Export(volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.Unexport(volumeName)
	defer client.DeleteVolume(volumeName)
	// set the export client
	err = client.SetExportClients(volumeName, clientList)
	if err != nil {
		panic(err)
	}

	// test getting the client list
	currentClients, err := client.GetExportClients(volumeName)
	if err != nil {
		panic(fmt.Sprintf("Unexpected error in GetExportClients: %v", err))
	}
	// verify we received the correct clients
	if len(currentClients) != len(clientList) {
		panic(fmt.Sprintf("Unexpected number of clients returned.  Expected: %d Actual: %d", len(clientList), len(currentClients)))
	}
	sort.Strings(currentClients)
	sort.Strings(clientList)
	for i := range currentClients {
		if currentClients[i] != clientList[i] {
			panic(fmt.Sprintf("Unexpected client returned.  Expected: %v Actual: %v", clientList, currentClients))
		}
	}
}

// Test SetExportClients()
func TestSetExportClients(*testing.T) {
	volumeName := "test_set_export"
	volumePath := client.Path(volumeName)
	clientList := []string{"1.2.3.4", "1.2.3.5"}

	// initialize the export
	_, err := client.CreateVolume(volumeName)
	if err != nil {
		panic(err)
	}
	err = client.Export(volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.Unexport(volumeName)
	defer client.DeleteVolume(volumeName)
	// verify we aren't already exporting the volume to any of the clients
	exports, err := client.GetIsiExports()
	for _, export := range exports {
		if export.Paths[0] == volumePath {
			for _, currentClient := range export.Clients {
				for _, newClient := range clientList {
					if currentClient == newClient {
						panic(fmt.Sprintf("Volume already exporting to %s: %v", newClient, export.Clients))
					}
				}
			}
		}
	}

	// test setting the export client
	err = client.SetExportClients(volumeName, clientList)
	if err != nil {
		panic(err)
	}

	// verify the export client was set
	sort.Strings(clientList)
	exports, err = client.GetIsiExports()
	for _, export := range exports {
		if export.Paths[0] == volumePath {
			// verify we received the correct clients
			if len(export.Clients) != len(clientList) {
				panic(fmt.Sprintf("Unexpected number of clients returned.  Expected: %d Actual: %d", len(clientList), len(export.Clients)))
			}
			sort.Strings(export.Clients)
			for i := range export.Clients {
				if export.Clients[i] != clientList[i] {
					panic(fmt.Sprintf("Unexpected client returned.  Expected: %v Actual: %v", clientList, export.Clients))
				}
			}
			// clients match so return.
			return
		}
	}
	panic(fmt.Sprintf("Volume %s not found in export list", volumePath))
}

// Test ClearExportClients()
func TestClearExportClients(*testing.T) {
	volumeName := "test_clear_export"
	volumePath := client.Path(volumeName)
	clientList := []string{"1.2.3.4", "1.2.3.5"}

	// initialize the export
	_, err := client.CreateVolume(volumeName)
	if err != nil {
		panic(err)
	}
	err = client.Export(volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.Unexport(volumeName)
	defer client.DeleteVolume(volumeName)
	// verify we are exporting the volume
	err = client.SetExportClients(volumeName, clientList)
	if err != nil {
		panic(err)
	}
	exports, err := client.GetIsiExports()
	sort.Strings(clientList)
	for _, export := range exports {
		if export.Paths[0] == volumePath {
			if len(export.Clients) != len(clientList) {
				panic(fmt.Sprintf("Unexpected number of clients returned.  Expected: %d Actual: %d", len(clientList), len(export.Clients)))
			}
			sort.Strings(export.Clients)
			for i := range export.Clients {
				if export.Clients[i] != clientList[i] {
					panic(fmt.Sprintf("Unexpected client returned.  Expected: %v Actual: %v", clientList, export.Clients))
				}
			}
		}
	}

	// test clearing the export client
	err = client.ClearExportClients(volumeName)
	if err != nil {
		panic(err)
	}

	// verify the export client was cleared
	exports, err = client.GetIsiExports()
	for _, export := range exports {
		if export.Paths[0] == volumePath {
			if len(export.Clients) > 0 {
				panic(fmt.Sprintf("Unexpected client address.  Expected: () Actual: (%s)", export.Clients[0]))
			}
		}
	}
}
