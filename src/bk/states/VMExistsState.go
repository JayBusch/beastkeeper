package states

import (
	"beastkeeper/src/bk/instanceTypes"
	//	"fmt"
	//	"os"
	//	"os/exec"
)

type VMExistsState struct {
	BaseState
}

func (self VMExistsState) Assess(instance instanceTypes.BaseInstance) bool {
	return false
}

func (self VMExistsState) Enforce(instance instanceTypes.BaseInstance) {
}
