package monitor

import (
	"log"

	"github.com/jjh2kiss/netlinkconnector/cnproc"
	"github.com/jjh2kiss/pasta/config"
	"github.com/jjh2kiss/pasta/stats"
	"github.com/jjh2kiss/pasta/system/process"
)

func processForkEvent(event *cnproc.ProcEvent, process_table *process.ProcessTable, stats_table *stats.StatsTable, config *config.Config, logger *log.Logger) error {
	ev, err := event.ForkEvent()
	if err != nil {
		return err
	}
	stats_table.Update(int(ev.ParentPid), stats.STAT_FORK)
	process_table.Add(int(ev.ChildPid))

	parent := process_table.GetOrDefault(int(ev.ParentPid))
	child := process_table.GetOrDefault(int(ev.ChildPid))

	if config.Quiet == false {
		if config.Events&event.What != 0 {
			name := parent.Cmdline.CombinedString(parent.KernelThread, config.Shortname, config.Dirstrip)
			logger.Printf("fork %5d parent %8s %s\n",
				ev.ParentPid,
				"",
				name,
			)

			name = child.Cmdline.CombinedString(parent.KernelThread, config.Shortname, config.Dirstrip)
			logger.Printf("fork %5d child  %8s %s\n",
				ev.ChildPid,
				"",
				name,
			)
		}
	}
	return nil
}
