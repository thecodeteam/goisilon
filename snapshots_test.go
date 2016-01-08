package goisilon

import (
	"fmt"
	"testing"
)

func init() {
	testClient()
}

func TestGetSnapshots(*testing.T) {
	snapshotPath := "/ifs/volumes"
	snapshotName1 := "test_get_snapshots_name1"
	snapshotName2 := "test_get_snapshots_name2"

	// identify all snapshots on the cluster
	snapshotMap := make(map[int64]string)
	snapshots, err := client.GetSnapshots()
	if err != nil {
		panic(err)
	}
	for _, snapshot := range snapshots {
		snapshotMap[snapshot.Id] = snapshot.Name
	}
	initialSnapshotCount := len(snapshots)

	// Add the test snapshots
	testSnapshot1, err := client.CreateSnapshot(snapshotPath, snapshotName1)
	if err != nil {
		panic(err)
	}
	testSnapshot2, err := client.CreateSnapshot(snapshotPath, snapshotName2)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.RemoveSnapshot(testSnapshot1.Id, snapshotName1)
	defer client.RemoveSnapshot(testSnapshot2.Id, snapshotName2)

	// get the updated snapshot list
	snapshots, err = client.GetSnapshots()
	if err != nil {
		panic(err)
	}

	// verify that the new snapshots are there as well as all the old snapshots.
	if len(snapshots) != initialSnapshotCount+2 {
		panic(fmt.Sprintf("Incorrect number of snapshots.  Expected: %d Actual: %d\n", initialSnapshotCount+2, len(snapshots)))
	}
	// remove the original snapshots and add the new ones.  in the end, we
	// should only have the snapshots we just created and nothing more.
	for _, snapshot := range snapshots {
		if _, found := snapshotMap[snapshot.Id]; found == true {
			// this snapshot existed prior to the test start
			delete(snapshotMap, snapshot.Id)
		} else {
			// this snapshot is new
			snapshotMap[snapshot.Id] = snapshot.Name
		}
	}
	if len(snapshotMap) != 2 {
		panic(fmt.Sprintf("Incorrect number of new exports.  Expected: 2 Actual: %d\n", len(snapshotMap)))
	}
	if _, found := snapshotMap[testSnapshot1.Id]; found == false {
		panic(fmt.Sprintf("testSnapshot1 was not in the snapshot list\n"))
	}
	if _, found := snapshotMap[testSnapshot2.Id]; found == false {
		panic(fmt.Sprintf("testSnapshot2 was not in the snapshot list\n"))
	}

}

func TestGetSnapshotsByPath(*testing.T) {
	snapshotPath1 := "/ifs/volumes"
	snapshotPath2 := "/ifs/data"
	snapshotName1 := "test_get_snapshots_by_path_name1"
	snapshotName2 := "test_get_snapshots_by_path_name2"
	snapshotName3 := "test_get_snapshots_by_path_name3"

	// identify all snapshots on the cluster
	snapshotMap := make(map[int64]string)
	snapshots, err := client.GetSnapshotsByPath(snapshotPath1)
	if err != nil {
		panic(err)
	}
	for _, snapshot := range snapshots {
		snapshotMap[snapshot.Id] = snapshot.Name
	}
	initialSnapshotCount := len(snapshots)

	// Add the test snapshots
	testSnapshot1, err := client.CreateSnapshot(snapshotPath1, snapshotName1)
	if err != nil {
		panic(err)
	}
	testSnapshot2, err := client.CreateSnapshot(snapshotPath2, snapshotName2)
	if err != nil {
		panic(err)
	}
	testSnapshot3, err := client.CreateSnapshot(snapshotPath1, snapshotName3)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.RemoveSnapshot(testSnapshot1.Id, snapshotName1)
	defer client.RemoveSnapshot(testSnapshot2.Id, snapshotName2)
	defer client.RemoveSnapshot(testSnapshot3.Id, snapshotName3)

	// get the updated snapshot list
	snapshots, err = client.GetSnapshotsByPath(snapshotPath1)
	if err != nil {
		panic(err)
	}

	// verify that the new snapshots in the given path are there as well
	// as all the old snapshots in that path.
	if len(snapshots) != initialSnapshotCount+2 {
		panic(fmt.Sprintf("Incorrect number of snapshots for path (%s).  Expected: %d Actual: %d\n", snapshotPath1, initialSnapshotCount+2, len(snapshots)))
	}
	// remove the original snapshots and add the new ones.  in the end, we
	// should only have the snapshots we just created and nothing more.
	for _, snapshot := range snapshots {
		if _, found := snapshotMap[snapshot.Id]; found == true {
			// this snapshot existed prior to the test start
			delete(snapshotMap, snapshot.Id)
		} else {
			// this snapshot is new
			snapshotMap[snapshot.Id] = snapshot.Name
		}
	}
	if len(snapshotMap) != 2 {
		panic(fmt.Sprintf("Incorrect number of new exports.  Expected: 2 Actual: %d\n", len(snapshotMap)))
	}
	if _, found := snapshotMap[testSnapshot1.Id]; found == false {
		panic(fmt.Sprintf("testSnapshot1 was not in the snapshot list\n"))
	}
	if _, found := snapshotMap[testSnapshot3.Id]; found == false {
		panic(fmt.Sprintf("testSnapshot3 was not in the snapshot list\n"))
	}
}

func TestCreateSnapshot(*testing.T) {
	snapshotPath := "/ifs/volumes"
	snapshotName := "test_get_create_snapshot_name"

	// make sure the snapshot doesn't exist yet
	snapshot, err := client.GetSnapshot(-1, snapshotName)
	if err == nil && snapshot != nil {
		panic(fmt.Sprintf("Snapshot (%s) already exists.\n", snapshotName))
	}

	// Add the test snapshot
	testSnapshot, err := client.CreateSnapshot(snapshotPath, snapshotName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.RemoveSnapshot(testSnapshot.Id, snapshotName)

	// get the updated snapshot list
	snapshot, err = client.GetSnapshot(testSnapshot.Id, snapshotName)
	if err != nil {
		panic(err)
	}
	if snapshot == nil {
		panic(fmt.Sprintf("Snapshot (%s) was not created.\n", snapshotName))
	}
	if snapshot.Name != snapshotName {
		panic(fmt.Sprintf("Snapshot name not set properly.  Expected: (%s) Actual: (%s)\n", snapshotName, snapshot.Name))
	}
	if snapshot.Path != snapshotPath {
		panic(fmt.Sprintf("Snapshot path not set properly.  Expected: (%s) Actual: (%s)\n", snapshotPath, snapshot.Path))
	}
}

func TestRemoveSnapshot(*testing.T) {
	snapshotPath := "/ifs/volumes"
	snapshotName := "test_remove_snapshot_name"

	// make sure the snapshot exists
	client.CreateSnapshot(snapshotPath, snapshotName)
	snapshot, err := client.GetSnapshot(-1, snapshotName)
	if err != nil {
		panic(err)
	}
	if snapshot == nil {
		panic(fmt.Sprintf("Test not setup properly.  No test snapshot (%s).", snapshotName))
	}

	// remove the snapshot
	err = client.RemoveSnapshot(snapshot.Id, snapshotName)
	if err != nil {
		panic(err)
	}

	// make sure the snapshot was removed
	snapshot, err = client.GetSnapshot(snapshot.Id, snapshotName)
	if err != nil {
		panic(err)
	}
	if snapshot != nil {
		panic(fmt.Sprintf("Snapshot (%s) was not removed.\n%+v\n", snapshotName, snapshot))
	}
}
