package v2

import (
	"os"
	"path"
	"strconv"

	"github.com/thecodeteam/goisilon/api"
)

const (
	namespacePath       = "namespace"
	exportsPath         = "platform/2/protocols/nfs/exports"
	quotaPath           = "platform/2/quota/quotas"
	snapshotsPath       = "platform/2/snapshot/snapshots"
	volumeSnapshotsPath = ".snapshot"
)

var (
	debug, _   = strconv.ParseBool(os.Getenv("GOISILON_DEBUG"))
	colonBytes = []byte{byte(':')}
)

func realNamespacePath(c api.Client) string {
	return path.Join(namespacePath, c.VolumesAccessPath())
}

func realExportsPath(c api.Client) string {
	return path.Join(exportsPath, c.VolumesAccessPath())
}

func realVolumeSnapshotPath(c api.Client, name string) string {
	return path.Join(realNamespacePath(c), volumeSnapshotsPath, name)
}
