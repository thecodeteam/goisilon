package goisilon

import (
	"fmt"
	"testing"
)

func init() {
	testClient()
}

func TestGetVolumes(*testing.T) {
	// TODO: Make this more robust
	volumes, err := client.GetVolumes()
	if err != nil {
		panic(err)
	}

	for _, volume := range volumes {
		fmt.Println(fmt.Sprintf("%+v", volume))
	}
}

func TestCreateVolume(*testing.T) {
	// TODO: Make this more robust
	volume, err := client.CreateVolume("testing")
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v", volume))

}

func TestGetVolume(*testing.T) {
	// TODO: Make this more robust
	volume, err := client.GetVolume("", "testing")
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v", volume))

}

func TestGetVolumeNone(*testing.T) {
	// TODO: Make this more robust
	volume, err := client.GetVolume("invalidvolume", "")
	if err != nil {
		panic(err)
	}

	if volume != nil {
		panic("invalid volume returned")
	}

}

func TestDeleteVolume(*testing.T) {
	// TODO: Make this more robust
	err := client.DeleteVolume("testing")
	if err != nil {
		panic(err)
	}

}

func TestExportVolume(*testing.T) {
	// TODO: Make this more robust
	err := client.ExportVolume("testing")
	if err != nil {
		panic(err)
	}

}

func TestUnexportVolume(*testing.T) {
	// TODO: Make this more robust
	err := client.UnexportVolume("testing")
	if err != nil {
		panic(err)
	}

}

func TestPath(*testing.T) {
	// TODO: Make this more robust
	fmt.Println(client.Path("testing"))
}

func TestGetVolumeExports(*testing.T) {
	// TODO: Make this more robust
	volumeExports, err := client.GetVolumeExports()
	if err != nil {
		panic(err)
	}

	for _, volumeExport := range volumeExports {
		fmt.Println(fmt.Sprintf("%+v", volumeExport))
	}
}
