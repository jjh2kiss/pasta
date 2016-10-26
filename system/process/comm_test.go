package process

import (
	"os/exec"
	"testing"
)

func TestNewCommByPid(t *testing.T) {
	testcases := []struct {
		expected []string
	}{
		{
			expected: []string{"sleep"},
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

		actual, err := NewCommByPid(pid)
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
