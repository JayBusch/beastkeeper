package states

import "beastkeeper/src/bk/instanceTypes"

type T_State interface {
	Assess(instanceTypes.BaseInstance) bool
	Enforce(instanceTypes.BaseInstance)
	Advance()
	GetAttempts() int
	SetAttempts(int)
	GetMaxAttempts() int
	SetMaxAttempts(int)
}

type BaseState struct {
	attempts    int
	maxAttempts int
}

// Simple function for tracking each attempt.
func (self *BaseState) Advance() {
	self.attempts = self.attempts + 1
}

// Available here only to provide interface compatability; derived types
// should use this function to return a boolean representing whether or
// not the state is currently enforced.
func (self BaseState) Assess(instance instanceTypes.BaseInstance) bool {
	return false
}

// Available here only to provide interface compatability; derived types
// should use this function to enforce the state; and expect that it is
// called by the state machine serially for each target system but
// concurrently across target systems.
func (self BaseState) Enforce(instance instanceTypes.BaseInstance) {

}

func (self BaseState) GetAttempts() int {
	return self.attempts
}

func (self *BaseState) SetAttempts(attempts int) {
	self.attempts = attempts
}

func (self BaseState) GetMaxAttempts() int {
	return self.maxAttempts
}

func (self *BaseState) SetMaxAttempts(maxAttempts int) {
	self.maxAttempts = maxAttempts
}
