package stats

import "github.com/jjh2kiss/pstat/system/process"

const (
	STAT_FORK = iota
	STAT_EXEC = iota
	STAT_EXIT = iota
	STAT_CORE = iota
	STAT_COMM = iota
	STAT_LAST = iota
)

type Stats struct {
	Cmdline *process.Cmdline
	Total   uint64
	Count   [STAT_LAST]uint64
}

type StatsMap map[string]*Stats
type StatsList []*Stats

func NewStats(cmdline *process.Cmdline) *Stats {
	if cmdline == nil {
		return nil
	}

	return &Stats{Cmdline: cmdline}
}

func (self *Stats) Update(event int) {
	if event < STAT_LAST {
		self.Total++
		self.Count[event]++
	}
}

func (self StatsMap) List() StatsList {
	stats := make(StatsList, 0, len(self))

	for _, v := range self {
		stats = append(stats, v)
	}
	return stats
}

//ByTotal implements sort.Interface for []*Stats based on the Total field.
type ByTotal StatsList

func (self ByTotal) Len() int {
	return len(self)
}

func (self ByTotal) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self ByTotal) Less(i, j int) bool {
	return self[i].Total < self[j].Total
}
