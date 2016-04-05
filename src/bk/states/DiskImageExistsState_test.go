package states

import (
	"beastkeeper/src/bk/instanceTypes"
	"fmt"
	"os"
	//"os/exec"
	"testing"
)

func Test_DiskImageExists_Assess(t *testing.T) {

	testDiskImageExistsState := DiskImageExistsState{}

	testInstance :=
		instanceTypes.NewBaseInstance("a9e91f49-cfa3-4910-a01c-83fc8c73e652",
			"Test_DiskImageExists_Assess_Image", instanceTypes.VM,
			float32(10.2), true, "../../../test/virtualMachines/")
	os.Remove(testInstance.GetDiskImageFileName())

	if testDiskImageExistsState.Assess(*testInstance) {
		t.Fatalf("Assess returning true when no image exists")
	}
	newfile, fileErr := os.Create(string(testInstance.GetDiskImageFileName()))
	if fileErr != nil {
		fmt.Printf("Error Creating File: %v\n", fileErr.Error())
	}
	newfile.Close()

	if !testDiskImageExistsState.Assess(*testInstance) {
		t.Fatalf("Asses returning false when image does exist")
	}

	delErr := os.Remove(testInstance.GetDiskImageFileName())
	if delErr != nil {
		t.Fatalf("Could not delete test image")
	}
}
