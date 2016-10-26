package process

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type ProcessTable struct {
	processes ProcessMap
	sync.RWMutex
}

func NewProcessTable() (*ProcessTable, error) {
	return &ProcessTable{
		processes: make(ProcessMap, 200),
	}, nil
}

// Load All process from /proc/*
func (self *ProcessTable) Sync() error {
	processes, err := Processes()
	if err != nil {
		return err
	}

	self.Lock()
	self.processes = processes
	self.Unlock()
	return nil
}

func (self *ProcessTable) Add(pid int) error {
	p, err := NewProcess(pid)
	if err != nil {
		return err
	}

	self.Lock()
	defer self.Unlock()

	self.processes[pid] = p
	return nil
}

//동일한 프로세스가 이미 존재할 경우에만 update 를 수행한다.
func (self *ProcessTable) Update(pid int) error {
	p, err := NewProcess(pid)
	if err != nil {
		return err
	}

	self.Lock()
	defer self.Unlock()

	oldp, ok := self.processes[pid]
	if ok == true {
		if oldp.Equal(p) {
			self.processes[pid] = p
			return nil
		}
		return fmt.Errorf("Not the same process")
	}
	return fmt.Errorf("Process not exist")
}

func (self *ProcessTable) Delete(pid int) {
	self.Lock()
	defer self.Unlock()
	delete(self.processes, pid)
}

func (self *ProcessTable) Get(pid int) (*Process, error) {
	self.RLock()
	defer self.RUnlock()

	p, ok := self.processes[pid]
	if ok == false {
		return nil, fmt.Errorf("Invalid pid: %d", pid)
	}
	return p, nil
}

func (self *ProcessTable) GetOrDefault(pid int) *Process {
	process, err := self.Get(pid)
	if err != nil {
		cmdline, _ := NewCmdline(strings.NewReader("<unknown>"))
		process = &Process{
			Pid:          999999,
			Cmdline:      cmdline,
			Stat:         &ProcessStat{Pid: 999999},
			TimeStamp:    time.Now(),
			KernelThread: false,
		}
	}

	return process

}

func (self *ProcessTable) Clone() (*ProcessTable, error) {
	new_table, err := NewProcessTable()
	if err != nil {
		return nil, err
	}

	self.RLock()
	defer self.RUnlock()

	for k, v := range self.processes {
		//v 는 process 포인터이다. 사용하는쪽에서 변경하게 되면 같이 변경된다.
		new_table.processes[k] = v
	}
	return new_table, nil
}
