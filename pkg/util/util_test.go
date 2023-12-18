package util

import "testing"

func TestRandomSize(t *testing.T) {
	randomNumeric := RandomNumeric(6)
	t.Log(randomNumeric)
}
