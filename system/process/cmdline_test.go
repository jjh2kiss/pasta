package process

import (
	"os/exec"
	"strings"
	"testing"
)

func StringListEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, _ := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNewCmdline(t *testing.T) {
	testcases := []struct {
		cmdline  string
		expected []string
	}{
		{
			cmdline:  "hello\x00world\x00",
			expected: []string{"hello", "world"},
		},
		{
			cmdline:  "hello\x00world",
			expected: []string{"hello", "world"},
		},
		{
			cmdline:  "hello",
			expected: []string{"hello"},
		},
		{
			cmdline:  "",
			expected: []string{},
		},
		{
			cmdline:  "\x00",
			expected: []string{},
		},
	}

	for _, testcase := range testcases {
		reader := strings.NewReader(testcase.cmdline)
		actual, err := NewCmdline(reader)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		if StringListEqual(actual.slice, testcase.expected) == false {
			t.Errorf("expected\n%s(%d)\nBut\n%s(%d)\n",
				testcase.expected,
				len(testcase.expected),
				actual.slice,
				len(actual.slice),
			)
		}
	}
}

func TestNewCmdlineByPid(t *testing.T) {
	testcases := []struct {
		expected []string
	}{
		{
			expected: []string{"sleep", "5"},
		},
	}

	for _, testcase := range testcases {
		//sleep 10으로 임으의 프로세스를 생성한다.
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

		actual, err := NewCmdlineByPid(pid)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		if StringListEqual(actual.slice, testcase.expected) == false {
			t.Errorf("expected\n%s(%d)\nBut\n%s(%d)\n",
				testcase.expected,
				len(testcase.expected),
				actual.slice,
				len(actual.slice),
			)
		}
	}
}

func TestCmdlineString(t *testing.T) {
	testcases := []struct {
		cmdline  *Cmdline
		expected string
	}{
		{
			cmdline:  &Cmdline{slice: []string{"hello", "world"}},
			expected: "hello world",
		},
		{
			cmdline:  &Cmdline{slice: []string{"hello"}},
			expected: "hello",
		},
		{
			cmdline:  &Cmdline{slice: []string{}},
			expected: "",
		},
	}

	for _, testcase := range testcases {
		actual := testcase.cmdline.String()

		if actual != testcase.expected {
			t.Errorf("expected\n%s\nBut\n%s\n",
				testcase.expected,
				actual,
			)
		}
	}
}

func TestCmdlineShortString(t *testing.T) {
	testcases := []struct {
		cmdline  *Cmdline
		count    int
		expected string
	}{
		{
			cmdline:  &Cmdline{slice: []string{"hello", "world", "1", "2", "3"}},
			count:    1,
			expected: "hello",
		},
		{
			cmdline:  &Cmdline{slice: []string{"hello", "world", "1", "2", "3"}},
			count:    2,
			expected: "hello world",
		},
		{
			cmdline:  &Cmdline{slice: []string{"hello", "world", "1", "2", "3"}},
			count:    3,
			expected: "hello world 1",
		},
		{
			cmdline:  &Cmdline{slice: []string{""}},
			count:    3,
			expected: "",
		},
		{
			cmdline:  &Cmdline{slice: []string{""}},
			count:    2,
			expected: "",
		},
		{
			cmdline:  &Cmdline{slice: []string{""}},
			count:    1,
			expected: "",
		},
	}

	for _, testcase := range testcases {
		actual := testcase.cmdline.ShortString(testcase.count)

		if actual != testcase.expected {
			t.Errorf("expected\n%s\nBut\n%s\n",
				testcase.expected,
				actual,
			)
		}
	}
}

func TestCmdlineCombinedString(t *testing.T) {
	testcases := []struct {
		cmdline      *Cmdline
		kernelThread bool
		short        bool
		dirstrip     bool
		expected     string
	}{
		//kernel thread
		{
			cmdline:      &Cmdline{slice: []string{"kworker"}},
			kernelThread: true,
			short:        false,
			dirstrip:     false,
			expected:     "[kworker]",
		},
		{
			cmdline:      &Cmdline{slice: []string{"kworker"}},
			kernelThread: false,
			short:        false,
			dirstrip:     false,
			expected:     "kworker",
		},

		//short
		{
			cmdline:      &Cmdline{slice: []string{"kworker", "1", "2"}},
			kernelThread: false,
			short:        true,
			dirstrip:     false,
			expected:     "kworker",
		},
		{
			cmdline:      &Cmdline{slice: []string{"kworker", "1", "2"}},
			kernelThread: false,
			short:        false,
			dirstrip:     false,
			expected:     "kworker 1 2",
		},

		//dirstrip
		{
			cmdline:      &Cmdline{slice: []string{"/bin/bash"}},
			kernelThread: false,
			short:        false,
			dirstrip:     true,
			expected:     "bash",
		},
		{
			cmdline:      &Cmdline{slice: []string{"/bin/bash"}},
			kernelThread: false,
			short:        false,
			dirstrip:     false,
			expected:     "/bin/bash",
		},

		//short on + dirstrip on
		{
			cmdline:      &Cmdline{slice: []string{"/bin/bash", "1", "2"}},
			kernelThread: false,
			short:        true,
			dirstrip:     true,
			expected:     "bash",
		},
		//short off + dirstrip on
		{
			cmdline:      &Cmdline{slice: []string{"/bin/bash", "1", "2"}},
			kernelThread: false,
			short:        false,
			dirstrip:     true,
			expected:     "bash 1 2",
		},
	}

	for _, testcase := range testcases {
		actual := testcase.cmdline.CombinedString(testcase.kernelThread,
			testcase.short,
			testcase.dirstrip)

		if actual != testcase.expected {
			t.Errorf("expected\n%s\nBut\n%s\n",
				testcase.expected,
				actual,
			)
		}
	}
}
