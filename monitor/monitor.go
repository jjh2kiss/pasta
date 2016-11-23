package monitor

import (
	"fmt"
	"log"
	"os"

	"github.com/jjh2kiss/netlinkconnector/cnproc"
	"github.com/jjh2kiss/pasta/config"
	"github.com/jjh2kiss/pasta/stats"
	"github.com/jjh2kiss/pasta/system/process"
)

func Monitor(config *config.Config, done <-chan struct{}) error {
	header_fmt := "%-19s %-6s %-5s %s %s %s\n"

	logger := log.New(os.Stdout, "", log.LstdFlags)
	err_logger := log.New(os.Stderr, "Error", log.LstdFlags)

	subscriber, err := cnproc.NewSubscriber()
	if err != nil {
		return err
	}
	defer subscriber.Close()

	ch, err := subscriber.Subscribe()
	if err != nil {
		return err
	}

	process_table, err := process.NewProcessTable()
	if err != nil {
		return err
	}
	process_table.Sync()

	stats_table := stats.NewStatsTable(process_table, config)

	if !config.Quiet {
		fmt.Printf(header_fmt,
			"Time",
			"Event",
			"PID",
			"Info",
			"Duration",
			"Process",
		)
	}

	for {
		select {
		case event := <-ch:
			switch event.What {
			case cnproc.PROC_EVENT_FORK:
				err := processForkEvent(event, process_table, stats_table, config, logger)
				if err != nil {
					err_logger.Println(err.Error())
				}
			case cnproc.PROC_EVENT_EXEC:
				err := processExecEvent(event, process_table, stats_table, config, logger)
				if err != nil {
					err_logger.Println(err.Error())
				}
			case cnproc.PROC_EVENT_EXIT:
				err := processExitEvent(event, process_table, stats_table, config, logger)
				if err != nil {
					err_logger.Println(err.Error())
				}
			case cnproc.PROC_EVENT_COREDUMP:
				err := processCoredumpEvent(event, process_table, stats_table, config, logger)
				if err != nil {
					err_logger.Println(err.Error())
				}
			case cnproc.PROC_EVENT_COMM:
				err := processCommEvent(event, process_table, stats_table, config, logger)
				if err != nil {
					err_logger.Println(err.Error())
				}
			}
		case <-done:
			goto End
		}
	}

End:
	if config.Statistics {
		SpacePrinter(os.Stdout, stats_table.List())
	}

	return nil
}
