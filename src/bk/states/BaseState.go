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

func (self *BaseState) Advance() {
	self.attempts = self.attempts + 1
}

func (self BaseState) Assess(instance instanceTypes.BaseInstance) bool {
	return false
}

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
