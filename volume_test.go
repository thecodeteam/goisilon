package goisilon

import (
	"fmt"
	"testing"
)

func init() {
	testClient()
}

func TestGetVolumes(*testing.T) {
	volumes, err := client.GetVolumes()
	if err != nil {
		panic(err)
	}

	for _, volume := range volumes {
		fmt.Println(fmt.Sprintf("%+v", volume))
	}
}

func TestCreateVolume(*testing.T) {
	volume, err := client.CreateVolume("testing")
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v", volume))

}

func TestGetVolume(*testing.T) {
	volume, err := client.GetVolume("", "testing")
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v", volume))

}

func TestGetVolumeNone(*testing.T) {
	volume, err := client.GetVolume("invalidvolume", "")
	if err != nil {
		panic(err)
	}

	if volume != nil {
		panic("invalid volume returned")
	}

}

func TestDeleteVolume(*testing.T) {
	err := client.DeleteVolume("testing")
	if err != nil {
		panic(err)
	}

}

func TestExportVolume(*testing.T) {
	err := client.ExportVolume("testing")
	if err != nil {
		panic(err)
	}

}

func TestUnexportVolume(*testing.T) {
	err := client.UnexportVolume("testing")
	if err != nil {
		panic(err)
	}

}

func TestPath(*testing.T) {
	fmt.Println(client.Path("testing"))

}

func TestGetVolumeExports(*testing.T) {
	volumeExports, err := client.GetVolumeExports()
	if err != nil {
		panic(err)
	}

	for _, volumeExport := range volumeExports {
		fmt.Println(fmt.Sprintf("%+v", volumeExport))
	}
}
