package system

import "testing"

func TestBootTimeInternal(t *testing.T) {
	btime, err := bootTime()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if btime == 0 {
		t.Errorf("Fail to get boot time from /proc/stat")
	}
}

func TestBootTime(t *testing.T) {
	for i := 0; i < 10; i++ {
		btime := BootTime()

		if btime == 0 {
			t.Errorf("Fail to get boot time")
		}
	}
}
