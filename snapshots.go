package goisilon

import (
	"errors"
	"fmt"
	papi "github.com/emccode/goisilon/api/v1"
	"strings"
)

type SnapshotList []*papi.IsiSnapshot
type Snapshot *papi.IsiSnapshot

func (c *Client) GetSnapshots() (SnapshotList, error) {
	snapshots, err := c.api.GetIsiSnapshots()
	if err != nil {
		return nil, err
	}

	return snapshots.SnapshotList, nil
}

func (c *Client) GetSnapshotsByPath(path string) (SnapshotList, error) {
	snapshots, err := c.api.GetIsiSnapshots()
	if err != nil {
		return nil, err
	}
	// find all the snapshots with the same path
	snapshotsWithPath := make(SnapshotList, 0, len(snapshots.SnapshotList))
	for _, snapshot := range snapshots.SnapshotList {
		if snapshot.Path == c.Path(path) {
			snapshotsWithPath = append(snapshotsWithPath, snapshot)
		}
	}
	return snapshotsWithPath, nil
}

func (c *Client) GetSnapshot(id int64, name string) (Snapshot, error) {
	// if we have an id, use it to find the snapshot
	snapshot, err := c.api.GetIsiSnapshot(id)
	if err == nil {
		return snapshot, nil
	}

	// there's no id or it didn't match, iterate through all snapshots and match
	// based on name
	if name == "" {
		return nil, err
	}
	snapshotList, err := c.GetSnapshots()
	if err != nil {
		return nil, err
	}

	for _, snapshot = range snapshotList {
		if snapshot.Name == name {
			return snapshot, nil
		}
	}

	return nil, nil
}

func (c *Client) CreateSnapshot(path, name string) (Snapshot, error) {
	return c.api.CreateIsiSnapshot(c.Path(path), name)
}

func (c *Client) RemoveSnapshot(id int64, name string) error {
	snapshot, err := c.GetSnapshot(id, name)
	if err != nil {
		return err
	}

	return c.api.RemoveIsiSnapshot(snapshot.Id)
}

func (c *Client) CopySnapshot(sourceId int64, sourceName, destinationName string) (Volume, error) {
	snapshot, err := c.GetSnapshot(sourceId, sourceName)
	if err != nil {
		return nil, err
	}
	if snapshot == nil {
		return nil, errors.New(fmt.Sprintf("Snapshot doesn't exist: (%d, %s)", sourceId, sourceName))
	}

	_, err = c.api.CopyIsiSnapshot(snapshot.Name, c.NameFromPath(snapshot.Path), destinationName)
	if err != nil {
		return nil, err
	}

	return c.GetVolume(destinationName, destinationName)
}

func (c *Client) NameFromPath(path string) string {
	// the name is the last entry in the path
	tokens := strings.Split(path, "/")
	return tokens[len(tokens)-1]
}
