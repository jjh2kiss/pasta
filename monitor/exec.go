package monitor

import (
	"log"

	"github.com/jjh2kiss/netlinkconnector/cnproc"
	"github.com/jjh2kiss/pasta/config"
	"github.com/jjh2kiss/pasta/stats"
	"github.com/jjh2kiss/pasta/system/process"
)

func processExecEvent(event *cnproc.ProcEvent, process_table *process.ProcessTable, stats_table *stats.StatsTable, config *config.Config, logger *log.Logger) error {
	ev, err := event.ExecEvent()
	if err != nil {
		return err
	}
	stats_table.Update(int(ev.ProcessPid), stats.STAT_EXEC)
	process_table.Update(int(ev.ProcessPid))

	process := process_table.GetOrDefault(int(ev.ProcessPid))

	if config.Quiet == false {
		if config.Events&event.What != 0 {
			name := process.Cmdline.CombinedString(process.KernelThread, config.Shortname, config.Dirstrip)

			logger.Printf("exec %5d        %8s %s\n",
				ev.ProcessPid,
				"",
				name,
			)
		}
	}
	return nil
}
