package process

import (
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/jjh2kiss/pstat/system"
)

type Process struct {
	Pid          int
	Cmdline      *Cmdline
	Stat         *ProcessStat
	TimeStamp    time.Time
	KernelThread bool
}

type ProcessList []*Process
type ProcessMap map[int]*Process

func (self *Process) Equal(other *Process) bool {
	if self.Pid == other.Pid && self.Stat.StartTime == other.Stat.StartTime {
		return true
	}
	return false
}

func NewProcess(pid int) (*Process, error) {
	now := time.Now()

	cmdline, err := NewCmdlineByPid(pid)
	if err != nil {
		return nil, err
	}

	stat, err := NewProcessStatByPid(pid)
	if err != nil {
		return nil, err
	}

	return &Process{
		Pid:          pid,
		Cmdline:      cmdline,
		TimeStamp:    now,
		Stat:         stat,
		KernelThread: system.KernelThread(pid),
	}, nil
}

func Pids() ([]int, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	fnames, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	pids := make([]int, 0, len(fnames))
	for _, fname := range fnames {
		pid, err := strconv.Atoi(fname)
		if err != nil {
			// if not numeric name, just skip
			continue
		}
		pids = append(pids, pid)
	}

	sort.Ints(pids)
	return pids, nil
}

func Processes() (ProcessMap, error) {
	pids, err := Pids()
	if err != nil {
		return nil, err
	}

	processes := make(ProcessMap, len(pids))
	for _, pid := range pids {
		process, err := NewProcess(pid)
		if err != nil {
			continue
		}
		processes[pid] = process
	}

	return processes, nil
}

//return process list
func (self ProcessList) Map() ProcessMap {
	processes := make(ProcessMap, len(self))
	for _, p := range self {
		processes[p.Pid] = p
	}
	return processes
}

//return processes map
func (self ProcessMap) List() ProcessList {
	processes := make(ProcessList, 0, len(self))

	for _, v := range self {
		processes = append(processes, v)
	}
	return processes
}
