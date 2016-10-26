package process

import (
	"os"
	"path"
	"strconv"
)

func NewCommByPid(pid int) (*Cmdline, error) {
	filepath := path.Join("/proc", strconv.Itoa(pid), "comm")
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewCmdline(f)
}
