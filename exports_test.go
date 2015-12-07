package goisilon

import (
	"fmt"
	"testing"
)

func init() {
	testClient()
}

func TestGetIsiExports(*testing.T) {
	exports, err := client.GetIsiExports()
	if err != nil {
		panic(err)
	}

	for _, export := range exports {
		fmt.Println(fmt.Sprintf("%+v", export))
	}
}

func TestExport(*testing.T) {
	err := client.Export("/ifs/data/docker/volumes")
	if err != nil {
		panic(err)
	}

}

func TestUnexport(*testing.T) {
	err := client.Unexport("/ifs/data/docker/volumes")
	if err != nil {
		panic(err)
	}
}
