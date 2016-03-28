package states

import (
	"beastkeeper/src/bk/instanceTypes"
	//"fmt"
	//"os"
	//"os/exec"
)

type InstallImageExistsState struct {
	BaseState
}

func (self InstallImageExistsState) Assess(instance instanceTypes.BaseInstance) bool {
	return false
}

func (self InstallImageExistsState) Enforce(instance instanceTypes.BaseInstance) {
}
