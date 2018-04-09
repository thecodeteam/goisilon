package v1

import (
	"os"
	"path"
	"strconv"

	"github.com/thecodeteam/goisilon/api"
)

const (
	namespacePath       = "namespace"
	exportsPath         = "platform/1/protocols/nfs/exports"
	quotaPath           = "platform/1/quota/quotas"
	snapshotsPath       = "platform/1/snapshot/snapshots"
	volumesnapshotsPath = ".snapshot"
)

var (
	debug, _ = strconv.ParseBool(os.Getenv("GOISILON_DEBUG"))
)

func realNamespacePath(client api.Client) string {
	return path.Join(namespacePath, client.VolumesAccessPath())
}

func realexportsPath(client api.Client) string {
	return path.Join(exportsPath, client.VolumesAccessPath())
}

func realVolumeSnapshotPath(client api.Client, name string) string {
	return path.Join(realNamespacePath(client), volumesnapshotsPath, name)
}
