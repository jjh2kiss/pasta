package monitor

import (
	"log"

	"github.com/jjh2kiss/netlinkconnector/cnproc"
	"github.com/jjh2kiss/pstat/config"
	"github.com/jjh2kiss/pstat/stats"
	"github.com/jjh2kiss/pstat/system/process"
)

func processCommEvent(event *cnproc.ProcEvent, process_table *process.ProcessTable, stats_table *stats.StatsTable, config *config.Config, logger *log.Logger) error {
	ev, err := event.CommEvent()
	if err != nil {
		return err
	}
	stats_table.Update(int(ev.ProcessPid), stats.STAT_COMM)

	process := process_table.GetOrDefault(int(ev.ProcessPid))

	if config.Quiet == false {
		if config.Events&event.What != 0 {
			name := process.Cmdline.CombinedString(process.KernelThread, config.Shortname, config.Dirstrip)
			logger.Printf("comm %5d        %8s %s -> %s\n",
				ev.ProcessPid,
				"",
				name,
				ev.Comm,
			)
		}
	}
	return nil
}
