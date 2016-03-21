package states

import (
	"beastkeeper/src/bk/instanceTypes"
	"fmt"
	"os"
	"os/exec"
)

type DiskImageExistsState struct {
	BaseState
}

func (self DiskImageExistsState) Assess(instance instanceTypes.BaseInstance) bool {
	if instance.Type == instanceTypes.VM {
		if _, err := os.Stat(instance.GetDiskImageFileName()); err == nil {
			return true
		}
	}
	return false
}

func (self DiskImageExistsState) Enforce(instance instanceTypes.BaseInstance) {
	if _, fileErr := os.Stat(instance.GetDiskImageFileName()); os.IsNotExist(fileErr) {

		cmd := exec.Command("truncate", "-s", "1GB", instance.GetDiskImageFileName())

		cmdErr := cmd.Run()

		if cmdErr != nil {
			fmt.Printf("ERROR: %s\n", cmdErr)
		}
	}
}
