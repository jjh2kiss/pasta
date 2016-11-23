package stats

import (
	"sync"

	"github.com/jjh2kiss/pasta/config"
	"github.com/jjh2kiss/pasta/system/process"
)

type StatsTable struct {
	processes    *process.ProcessTable
	table        StatsMap
	config       *config.Config
	sync.RWMutex //protect table
}

//make new Stats table
func NewStatsTable(process_table *process.ProcessTable, config *config.Config) *StatsTable {
	return &StatsTable{
		processes: process_table,
		config:    config,
		table:     make(StatsMap, 1024),
	}
}

// update statistics
func (self *StatsTable) Update(pid int, event int) error {
	process, err := self.processes.Get(pid)
	if err != nil {
		return err
	}

	name := process.Cmdline.CombinedString(process.KernelThread, self.config.Shortname, self.config.Dirstrip)

	self.RLock()
	stats, ok := self.table[name]
	self.RUnlock()

	if ok == false {
		stats = NewStats(name)

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
