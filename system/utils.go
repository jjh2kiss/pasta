package system

import (
	"fmt"
	"os"

	"github.com/jjh2kiss/pasta/system/sched"
)

var running_in_container bool

func RunningInContainer() bool {
	return running_in_container
}

func KernelThread(pid int) bool {
	path := fmt.Sprintf("/proc/%d/exe", pid)
	_, err := os.Readlink(path)
	if err != nil {
		return true
	}
	return false
}

func init() {
	//caching running_in_container
	init_process, err := sched.GetInitProcess()
	if err != nil {
		running_in_container = true
	} else {
		running_in_container = init_process.RunningInContainer()
	}
}
