package system

import "testing"

func TestClockTick(t *testing.T) {
	ct := ClockTick()

	if ct == 0 {
		t.Errorf("Fail to get ClockTick")
	}
}
