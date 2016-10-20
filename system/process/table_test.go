package process

import (
	"os"
	"testing"
)

func TestNewProcessTable(t *testing.T) {
	_, err := NewProcessTable()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestProcessTableSync(t *testing.T) {
	table, err := NewProcessTable()
	if err != nil {
		t.Errorf(err.Error())
	}

	err = table.Sync()
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(table.processes) == 0 {
		t.Errorf("not exist any process")
	}
}

func TestProcessTableAdd(t *testing.T) {
	table, err := NewProcessTable()
	if err != nil {
		t.Errorf(err.Error())
	}

	pid := os.Getpid()
	err = table.Add(pid)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, ok := table.processes[pid]
	if ok == false {
		t.Errorf("%d process does not exist", pid)
	}

}

func TestProcessTableUpdate(t *testing.T) {
	table, err := NewProcessTable()
	if err != nil {
		t.Errorf(err.Error())
	}
	table.Sync()

	err = table.Update(os.Getpid())
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestProcessTableDelete(t *testing.T) {
	table, err := NewProcessTable()
	if err != nil {
		t.Errorf(err.Error())
	}
	table.Sync()
	table.Delete(os.Getpid())

	_, ok := table.processes[os.Getpid()]
	if ok == true {
		t.Errorf("expect not existence, but exist")
	}
}

func TestProcessTableGet(t *testing.T) {
	table, err := NewProcessTable()
	if err != nil {
		t.Errorf(err.Error())
	}
	table.Sync()

	p, err := table.Get(os.Getpid())
	if err != nil {
		t.Errorf(err.Error())
	}

	if p.Pid != os.Getpid() {
		t.Errorf("pid not matched")
	}
}

func TestProcessTableClone(t *testing.T) {
	table, err := NewProcessTable()
	if err != nil {
		t.Errorf(err.Error())
	}
	table.Sync()

	clone, err := table.Clone()
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(clone.processes) != len(table.processes) {
		t.Errorf("processes length is not same")
	}

}
