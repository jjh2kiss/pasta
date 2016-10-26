package utils

import (
	"testing"

	"github.com/jjh2kiss/netlinkconnector/cnproc"
)

func TestProcEventUint32(t *testing.T) {
	testcases := []struct {
		name     string
		expected uint32
	}{
		{"none", cnproc.PROC_EVENT_NONE},
		{"fork", cnproc.PROC_EVENT_FORK},
		{"exec", cnproc.PROC_EVENT_EXEC},
		{"uid", cnproc.PROC_EVENT_UID},
		{"gid", cnproc.PROC_EVENT_GID},
		{"sid", cnproc.PROC_EVENT_SID},
		{"ptrace", cnproc.PROC_EVENT_PTRACE},
		{"comm", cnproc.PROC_EVENT_COMM},
		{"coredump", cnproc.PROC_EVENT_COREDUMP},
		{"exit", cnproc.PROC_EVENT_EXIT},
		{"abcd", cnproc.PROC_EVENT_NONE},
		{"Fork", cnproc.PROC_EVENT_FORK},
		{"Exec", cnproc.PROC_EVENT_EXEC},
		{"Uid", cnproc.PROC_EVENT_UID},
		{"Gid", cnproc.PROC_EVENT_GID},
		{"Sid", cnproc.PROC_EVENT_SID},
		{"Ptrace", cnproc.PROC_EVENT_PTRACE},
		{"Comm", cnproc.PROC_EVENT_COMM},
		{"Coredump", cnproc.PROC_EVENT_COREDUMP},
		{"Exit", cnproc.PROC_EVENT_EXIT},
		{"Abcd", cnproc.PROC_EVENT_NONE},
	}

	for _, testcase := range testcases {
		actual := ProcEventUint32(testcase.name)
		if actual != testcase.expected {
			t.Errorf("expected %x but %x\n", testcase.expected, actual)
		}
	}
}
