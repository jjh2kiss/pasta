package utils

import (
	"strings"

	"github.com/jjh2kiss/netlinkconnector/cnproc"
)

func ProcEventUint32(name string) uint32 {
	name = strings.ToLower(name)
	switch name {
	case "none":
		return cnproc.PROC_EVENT_NONE
	case "fork":
		return cnproc.PROC_EVENT_FORK
	case "exec":
		return cnproc.PROC_EVENT_EXEC
	case "uid":
		return cnproc.PROC_EVENT_UID
	case "gid":
		return cnproc.PROC_EVENT_GID
	case "sid":
		return cnproc.PROC_EVENT_SID
	case "ptrace":
		return cnproc.PROC_EVENT_PTRACE
	case "comm":
		return cnproc.PROC_EVENT_COMM
	case "coredump":
		return cnproc.PROC_EVENT_COREDUMP
	case "exit":
		return cnproc.PROC_EVENT_EXIT
	}
	return cnproc.PROC_EVENT_NONE
}
