package states

import "testing"

func TestAdvance(t *testing.T) {

	testBaseState := BaseState{}

	testBaseState.SetMaxAttempts(5)
	testBaseState.SetAttempts(0)

	if testBaseState.GetAttempts() != 0 {
		t.Fatalf("attempt not initialized")
	}

	testBaseState.Advance()

	if testBaseState.GetAttempts() != 1 {
		t.Fatalf("Advance not incrementing attempt count")
	}

}
