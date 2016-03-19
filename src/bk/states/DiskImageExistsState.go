package states

type DiskImageExistsState struct {
	BaseState
}

func (self DiskImageExistsState) assess() bool {
	return false
}
