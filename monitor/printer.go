package monitor

import (
	"fmt"
	"io"
	"sort"

	"github.com/jjh2kiss/pasta/stats"
)

func SpacePrinter(output io.Writer, l stats.StatsList) {
	sort.Sort(sort.Reverse(stats.ByTotal(l)))

	fmt.Fprintln(output)
	fmt.Fprintf(output, "%8s %8s %8s %8s %8s %8s %8s\n", "Fork", "Exec", "Exit", "Coredump", "Comm", "Total", "Process")

	for _, s := range l {
		fmt.Fprintf(output, "%8d %8d %8d %8d %8d %8d %s\n",
			s.Count[stats.STAT_FORK],
			s.Count[stats.STAT_EXEC],
			s.Count[stats.STAT_EXIT],
			s.Count[stats.STAT_CORE],
			s.Count[stats.STAT_COMM],
			s.Total,
			s.Cmdline,
		)
	}
	fmt.Fprintln(output)

}
