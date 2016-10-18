package sched

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type InitProcess struct {
	Name    string
	Pid     int
	Threads int
}

func newInitProcess(reader io.Reader) (*InitProcess, error) {
	scanner := bufio.NewScanner(reader)

	if !scanner.Scan() {
		return nil, fmt.Errorf("Can't read data from reader")
	}

	first_line := scanner.Text()

	fields := strings.SplitN(first_line, " ", 2)
	if len(fields) != 2 {
		return nil, fmt.Errorf("fail to parse first_line(%s)", first_line)
	}

	name := strings.TrimSpace(fields[0])
	other := strings.TrimSpace(fields[1])
	other = strings.TrimRight(strings.TrimLeft(other, "("), ")")

	//extract pid
	fields = strings.SplitN(other, ",", 2)
	if len(fields) != 2 {
		return nil, fmt.Errorf("fail to parse other fields(%s)", other)
	}

	pid, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}

	//extract threads
	other = strings.TrimSpace(fields[1])
	fields = strings.SplitN(other, ":", 2)
	if len(fields) != 2 {
		return nil, fmt.Errorf("fail to parse threads field(%s)", other)
	}

	threads, err := strconv.Atoi(strings.TrimSpace(fields[1]))
	if err != nil {
		return nil, err
	}

	return &InitProcess{
		Name:    name,
		Pid:     pid,
		Threads: threads,
	}, nil

}

func GetInitProcess() (*InitProcess, error) {
	f, err := os.Open("/proc/1/sched")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return newInitProcess(f)
}

func (self *InitProcess) RunningInContainer() bool {
	if self == nil {
		return false
	}

	if self.Pid == 1 {
		return false
	}

	return true
}
