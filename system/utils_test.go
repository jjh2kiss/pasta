package system

import (
	"os"
	"testing"
)

func TestKernelThread(t *testing.T) {
	testcases := []struct {
		pid                   int
		can_test_in_container bool
		expected              bool
	}{
		{1, true, false}, //init process
		{2, false, true}, //컨테이너에서 테스트하면 실해함.
		{os.Getpid(), true, false},
	}

	for _, testcase := range testcases {
		if RunningInContainer() && !testcase.can_test_in_container {
			continue
		}

		actual := KernelThread(testcase.pid)

		if actual != testcase.expected {
			t.Errorf("expected %v but %v (%v)\n", testcase.expected, actual, testcase)
		}
	}
}
