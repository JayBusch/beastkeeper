package states

type T_State interface {
	Assess() bool
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
func (self BaseState) Assess() bool {
	return true
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
