package process

import (
	"io"
	"os/exec"
	"strings"
	"testing"
)

func TestNewProcessStat(t *testing.T) {
	testcases := []struct {
		reader   io.Reader
		expected *ProcessStat
	}{
		{
			reader: strings.NewReader(`1 (systemd) S 0 1 1 0 -1 4194560 112756 4037474 61 2441 409 330 6008 1580 20 0 1 0 3 189865984 1446 18446744073709551615 1 1 0 0 0 0 671173123 4096 1260 0 0 0 17 3 0 0 7 0 0 0 0 0 0 0 0 0 0`),
			expected: &ProcessStat{
				Pid:                 1,
				Comm:                "systemd",
				State:               "S",
				PPid:                0,
				UserTime:            409,
				SystemTime:          330,
				ChildUserTime:       6008,
				ChildSystemTime:     1580,
				Threads:             1,
				StartTime:           3,
				VmSize:              189865984,
				VmRSS:               1446,
				DelayAcctBlkioTicks: 7,
			},
		},
		{
			reader: strings.NewReader(`1 ( systemd ) S 0 1 1 0 -1 4194560 112756 4037474 61 2441 409 330 6008 1580 20 0 1 0 3 189865984 1446 18446744073709551615 1 1 0 0 0 0 671173123 4096 1260 0 0 0 17 3 0 0 7 0 0 0 0 0 0 0 0 0 0`),
			expected: &ProcessStat{
				Pid:                 1,
				Comm:                "systemd",
				State:               "S",
				PPid:                0,
				UserTime:            409,
				SystemTime:          330,
				ChildUserTime:       6008,
				ChildSystemTime:     1580,
				Threads:             1,
				StartTime:           3,
				VmSize:              189865984,
				VmRSS:               1446,
				DelayAcctBlkioTicks: 7,
			},
		},
		{
			reader: strings.NewReader(`1 (systemd a) S 0 1 1 0 -1 4194560 112756 4037474 61 2441 409 330 6008 1580 20 0 1 0 3 189865984 1446 18446744073709551615 1 1 0 0 0 0 671173123 4096 1260 0 0 0 17 3 0 0 7 0 0 0 0 0 0 0 0 0 0`),
			expected: &ProcessStat{
				Pid:                 1,
				Comm:                "systemd a",
				State:               "S",
				PPid:                0,
				UserTime:            409,
				SystemTime:          330,
				ChildUserTime:       6008,
				ChildSystemTime:     1580,
				Threads:             1,
				StartTime:           3,
				VmSize:              189865984,
				VmRSS:               1446,
				DelayAcctBlkioTicks: 7,
			},
		},
	}

	for _, testcase := range testcases {
		actual, err := NewProcessStat(testcase.reader)
		if err != nil {
			t.Errorf(err.Error())
		}

		if actual.Equal(testcase.expected) == false {
			t.Errorf("expected\n%v\nBut\n%v", testcase.expected, actual)
		}
	}

}

func TestNewProcessStatNew(t *testing.T) {
	expected := &ProcessStat{
		Comm:  "sleep",
		State: "S",
	}

	cmd := exec.Command("sleep", "5")
	err := cmd.Start()

	if err != nil {
		t.Errorf(err.Error())
		return
	}
	pid := cmd.Process.Pid
	defer func() {
		_ = cmd.Process.Kill()
	}()

	go func() {
		cmd.Wait()
	}()

	actual, err := NewProcessStatByPid(pid)
	if err != nil {
		t.Errorf(err.Error())
	}

	if actual.Pid != pid || actual.Comm != expected.Comm {
		t.Errorf("expected\n%v\nBut\n%v", expected, actual)
	}

}

func TestNewProcessStatEqual(t *testing.T) {
	testcases := []struct {
		a        *ProcessStat
		b        *ProcessStat
		expected bool
	}{
		{
			a: &ProcessStat{
				Pid:                 1,
				Comm:                "systemd",
				State:               "S",
				PPid:                0,
				UserTime:            409,
				SystemTime:          330,
				ChildUserTime:       6008,
				ChildSystemTime:     1580,
				Threads:             1,
				StartTime:           3,
				VmSize:              189865984,
				VmRSS:               1446,
				DelayAcctBlkioTicks: 7,
			},
			b: &ProcessStat{
				Pid:                 1,
				Comm:                "systemd",
				State:               "S",
				PPid:                0,
				UserTime:            409,
				SystemTime:          330,
				ChildUserTime:       6008,
				ChildSystemTime:     1580,
				Threads:             1,
				StartTime:           3,
				VmSize:              189865984,
				VmRSS:               1446,
				DelayAcctBlkioTicks: 7,
			},
			expected: true,
		},
		{
			a: &ProcessStat{
				Pid:                 2,
				Comm:                "systemd",
				State:               "S",
				PPid:                0,
				UserTime:            409,
				SystemTime:          330,
				ChildUserTime:       6008,
				ChildSystemTime:     1580,
				Threads:             1,
				StartTime:           3,
				VmSize:              189865984,
				VmRSS:               1446,
				DelayAcctBlkioTicks: 7,
			},
			b: &ProcessStat{
				Pid:                 1,
				Comm:                "systemd",
				State:               "S",
				PPid:                0,
				UserTime:            409,
				SystemTime:          330,
				ChildUserTime:       6008,
				ChildSystemTime:     1580,
				Threads:             1,
				StartTime:           3,
				VmSize:              189865984,
				VmRSS:               1446,
				DelayAcctBlkioTicks: 7,
			},
			expected: false,
		},
		{
			a: &ProcessStat{
				Pid:                 1,
				Comm:                "systemda",
				State:               "S",
				PPid:                0,
				UserTime:            409,
				SystemTime:          330,
				ChildUserTime:       6008,
				ChildSystemTime:     1580,
				Threads:             1,
				StartTime:           3,
				VmSize:              189865984,
				VmRSS:               1446,
				DelayAcctBlkioTicks: 7,
			},
			b: &ProcessStat{
				Pid:                 1,
				Comm:                "systemd",
				State:               "S",
				PPid:                0,
				UserTime:            409,
				SystemTime:          330,
				ChildUserTime:       6008,
				ChildSystemTime:     1580,
				Threads:             1,
				StartTime:           3,
				VmSize:              189865984,
				VmRSS:               1446,
				DelayAcctBlkioTicks: 7,
			},
			expected: false,
		},
	}

	for _, testcase := range testcases {
		actual := testcase.a.Equal(testcase.b)
		if actual != testcase.expected {
			t.Errorf("expected\n%t\nBut\n%t", testcase.expected, actual)
		}
	}

}
