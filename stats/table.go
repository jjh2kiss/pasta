package stats

import (
	"sync"

	"github.com/jjh2kiss/pstat/system/process"
)

type StatsTable struct {
	processes    *process.ProcessTable
	table        StatsMap
	sync.RWMutex //protect table
}

//make new Stats table
func NewStatsTable(process_table *process.ProcessTable) *StatsTable {
	return &StatsTable{
		processes: process_table,
		table:     make(StatsMap, 1024),
	}
}

// update statistics
func (self *StatsTable) Update(pid int, event int) error {
	process, err := self.processes.Get(pid)
	if err != nil {
		return err
	}

	name := process.Cmdline.String()

	self.RLock()
	stats, ok := self.table[name]
	self.RUnlock()

	if ok == false {
		stats = NewStats(process.Cmdline)

		self.Lock()
		self.table[name] = stats
		self.Unlock()
	}

	stats.Update(event)
	return nil
}

func (self *StatsTable) List() StatsList {
	return self.table.List()
}
