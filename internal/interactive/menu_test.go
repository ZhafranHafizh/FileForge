package interactive

import "testing"

func TestIsCancelled(t *testing.T) {
	err := cancelIfInterrupted(assertErr("interrupted"))
	if !IsCancelled(err) {
		t.Fatal("expected cancelled error")
	}
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
