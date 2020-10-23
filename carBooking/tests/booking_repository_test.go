package tests

import "testing"


func TestFindAll(t *testing.T) {
	var test = 2
	if test != 2 {
		// `t.Error*` will report test failures but continue
		// executing the test. `t.Fatal*` will report test
		// failures and stop the test immediately.
		t.Errorf("IntMin(2, -2) = %d; want -2", test)
	}
}
