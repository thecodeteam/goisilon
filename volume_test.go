package goisilon

import (
	"fmt"
	"testing"
)

func TestGetVolumes(*testing.T) {
	volumeName1 := "test_get_volumes_name1"
	volumeName2 := "test_get_volumes_name2"

	// identify all volumes on the cluster
	volumeMap := make(map[string]bool)
	volumes, err := client.GetVolumes(defaultCtx)
	if err != nil {
		panic(err)
	}
	for _, volume := range volumes {
		volumeMap[volume.Name] = true
	}
	initialVolumeCount := len(volumes)

	// Add the test volumes
	testVolume1, err := client.CreateVolume(defaultCtx, volumeName1)
	if err != nil {
		panic(err)
	}
	testVolume2, err := client.CreateVolume(defaultCtx, volumeName2)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(defaultCtx, volumeName1)
	defer client.DeleteVolume(defaultCtx, volumeName2)

	// get the updated volume list
	volumes, err = client.GetVolumes(defaultCtx)
	if err != nil {
		panic(err)
	}

	// verify that the new volumes are there as well as all the old volumes.
	if len(volumes) != initialVolumeCount+2 {
		panic(fmt.Sprintf("Incorrect number of volumes.  Expected: %d Actual: %d\n", initialVolumeCount+2, len(volumes)))
	}
	// remove the original volumes and add the new ones.  in the end, we
	// should only have the volumes we just created and nothing more.
	for _, volume := range volumes {
		if _, found := volumeMap[volume.Name]; found == true {
			// this volume existed prior to the test start
			delete(volumeMap, volume.Name)
		} else {
			// this volume is new
			volumeMap[volume.Name] = true
		}
	}
	if len(volumeMap) != 2 {
		panic(fmt.Sprintf("Incorrect number of new volumes.  Expected: 2 Actual: %d\n", len(volumeMap)))
	}
	if _, found := volumeMap[testVolume1.Name]; found == false {
		panic(fmt.Sprintf("testVolume1 was not in the volume list\n"))
	}
	if _, found := volumeMap[testVolume2.Name]; found == false {
		panic(fmt.Sprintf("testVolume2 was not in the volume list\n"))
	}

}

func TestGetCreateVolume(*testing.T) {
	volumeName := "test_get_create_volume_name"

	// make sure the volume doesn't exist yet
	volume, err := client.GetVolume(defaultCtx, volumeName, volumeName)
	if err == nil && volume != nil {
		panic(fmt.Sprintf("Volume (%s) already exists.\n", volumeName))
	}

	// Add the test volume
	testVolume, err := client.CreateVolume(defaultCtx, volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(defaultCtx, testVolume.Name)

	// get the new volume
	volume, err = client.GetVolume(defaultCtx, volumeName, volumeName)
	if err != nil {
		panic(err)
	}
	if volume == nil {
		panic(fmt.Sprintf("Volume (%s) was not created.\n", volumeName))
	}
	if volume.Name != volumeName {
		panic(fmt.Sprintf("Volume name not set properly.  Expected: (%s) Actual: (%s)\n", volumeName, volume.Name))
	}
}

func TestDeleteVolume(*testing.T) {
	volumeName := "test_remove_volume_name"

	// make sure the volume exists
	client.CreateVolume(defaultCtx, volumeName)
	volume, err := client.GetVolume(defaultCtx, volumeName, volumeName)
	if err != nil {
		panic(err)
	}
	if volume == nil {
		panic(fmt.Sprintf("Test not setup properly.  No test volume (%s).", volumeName))
	}

	// remove the volume
	err = client.DeleteVolume(defaultCtx, volumeName)
	if err != nil {
		panic(err)
	}

	// make sure the volume was removed
	volume, err = client.GetVolume(defaultCtx, volumeName, volumeName)
	if err == nil {
		panic(fmt.Sprintf("Attempting to get a removed volume should return an error but returned nil"))
	}
	if volume != nil {
		panic(fmt.Sprintf("Volume (%s) was not removed.\n%+v\n", volumeName, volume))
	}
}

func TestCopyVolume(*testing.T) {
	sourceVolumeName := "test_copy_source_volume_name"
	destinationVolumeName := "test_copy_destination_volume_name"
	subDirectoryName := "test_sub_directory"
	sourceSubDirectoryPath := fmt.Sprintf("%s/%s", sourceVolumeName, subDirectoryName)
	destinationSubDirectoryPath := fmt.Sprintf("%s/%s", destinationVolumeName, subDirectoryName)

	// make sure the destination volume doesn't exist yet
	destinationVolume, err := client.GetVolume(
		defaultCtx, destinationVolumeName, destinationVolumeName)
	if err == nil && destinationVolume != nil {
		panic(fmt.Sprintf("Volume (%s) already exists.\n", destinationVolumeName))
	}

	// Add the test volume
	sourceTestVolume, err := client.CreateVolume(defaultCtx, sourceVolumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(defaultCtx, sourceTestVolume.Name)
	// add a sub directory to the source volume
	_, err = client.CreateVolume(defaultCtx, sourceSubDirectoryPath)
	if err != nil {
		panic(err)
	}

	// copy the source volume to the test volume
	destinationTestVolume, err := client.CopyVolume(
		defaultCtx, sourceVolumeName, destinationVolumeName)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, destinationTestVolume.Name)
	// verify the copied volume is the same as the source volume
	if destinationTestVolume == nil {
		panic(fmt.Sprintf("Destination volume (%s) was not created.\n", destinationVolumeName))
	}
	if destinationTestVolume.Name != destinationVolumeName {
		panic(fmt.Sprintf("Destination volume name not set properly.  Expected: (%s) Actual: (%s)\n", destinationVolumeName, destinationTestVolume.Name))
	}
	// make sure the destination volume contains the sub-directory
	subTestVolume, err := client.GetVolume(
		defaultCtx, "", destinationSubDirectoryPath)
	if err != nil {
		panic(err)
	}
	// verify the copied subdirectory is the same as int the source volume
	if subTestVolume == nil {
		panic(fmt.Sprintf("Destination sub directory (%s) was not created.\n", subDirectoryName))
	}
	if subTestVolume.Name != destinationSubDirectoryPath {
		panic(fmt.Sprintf("Destination sub directory name not set properly.  Expected: (%s) Actual: (%s)\n", destinationSubDirectoryPath, subTestVolume.Name))
	}

}

func TestExportVolume(*testing.T) {
	// TODO: Make this more robust
	_, err := client.ExportVolume(defaultCtx, "testing")
	if err != nil {
		panic(err)
	}

}

func TestUnexportVolume(*testing.T) {
	// TODO: Make this more robust
	err := client.UnexportVolume(defaultCtx, "testing")
	if err != nil {
		panic(err)
	}

}

func TestPath(*testing.T) {
	// TODO: Make this more robust
	fmt.Println(client.API.VolumePath("testing"))
}

func TestGetVolumeExportMap(t *testing.T) {
	// TODO: Make this more robust
	volExMap, err := client.GetVolumeExportMap(defaultCtx, false)
	assertNoError(t, err)
	for v := range volExMap {
		t.Logf("volName=%s, volPath=%s", v.Name, client.API.VolumePath(v.Name))
	}
}
