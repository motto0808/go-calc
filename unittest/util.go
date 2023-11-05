package unittest

import (
	"testing"
)

func assert(t *testing.T, cond bool, msg ...string) {
	if !cond {
		if len(msg) > 0 {
			t.Errorf("assert failed %s", msg[0])
		} else {
			t.Error("assert failed")
		}
	}
}
