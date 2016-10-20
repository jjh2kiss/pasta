package process

import (
	"os"
	"runtime"
	"sort"
	"testing"
)

func GetCurrentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

func TestPids(t *testing.T) {
	pids, err := Pids()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if len(pids) == 0 {
		t.Errorf("len(pid) should be greater than 0")
		return
	}

	if sort.IntsAreSorted(pids) == false {
		t.Errorf("pids is not sorted")
		return
	}
}

func TestNewProcess(t *testing.T) {
	pid := os.Getpid()
	_, err := NewProcess(pid)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewProcessInvalid(t *testing.T) {
	_, err := NewProcess(9999999)
	if err == nil {
		t.Errorf(err.Error())
	}
}

func TestProcessEqual(t *testing.T) {
	p, _ := NewProcess(os.Getpid())
	pp, _ := NewProcess(os.Getppid())

	testcases := []struct {
		a        *Process
		b        *Process
		expected bool
	}{
		{
			a:        p,
			b:        p,
			expected: true,
		},
		{
			a:        p,
			b:        pp,
			expected: false,
		},
		{
			a:        pp,
			b:        p,
			expected: false,
		},
	}

	for _, testcase := range testcases {
		actual := testcase.a.Equal(testcase.b)
		if actual != testcase.expected {
			t.Errorf("expected\n%v\nBut\n%v", testcase.expected, actual)
		}
	}

}

func TestProcesses(t *testing.T) {
	processes, err := Processes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if len(processes) == 0 {
		t.Errorf("len(processes) should be greater than 0")
		return
	}

	_, ok := processes[os.Getpid()]
	if ok == false {
		t.Errorf("processes should self process")
		return
	}
}

func TestProcessMapList(t *testing.T) {
	processes, err := Processes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	plist := processes.List()

	if len(plist) != len(processes) {
		t.Errorf("fail to make list")
	}
}

func TestProcessListMap(t *testing.T) {
	processes, err := Processes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	plist := processes.List()
	pmap := plist.Map()

	if len(pmap) != len(processes) {
		t.Errorf("fail to make map from list")
	}
}
