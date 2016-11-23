package stats

import (
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/jjh2kiss/pasta/system/process"
)

func TestStatsNewStat(t *testing.T) {
	self_cmdline, err := process.NewCmdlineByPid(os.Getpid())
	if err != nil {
		t.Errorf("fail to get cmdline of self")
		return
	}

	cmdline := self_cmdline.String()

	testcases := []struct {
		cmdline  string
		expected *Stats
	}{
		{
			cmdline: cmdline,
			expected: &Stats{
				Cmdline: cmdline,
			},
		},
	}

	for _, testcase := range testcases {
		actual := NewStats(testcase.cmdline)
		if reflect.DeepEqual(actual, testcase.expected) == false {
			t.Errorf("expected\n%v\nbut\n%v\n", testcase.expected, actual)
		}
	}
}

func TestStatsUpdate(t *testing.T) {
	self_cmdline, err := process.NewCmdlineByPid(os.Getpid())
	if err != nil {
		t.Errorf("fail to get cmdline of self")
		return
	}

	cmdline := self_cmdline.String()

	testcases := []struct {
		events   []int
		expected *Stats
	}{
		{
			events: []int{STAT_FORK},
			expected: &Stats{
				Cmdline: cmdline,
				Total:   1,
				Count: [STAT_LAST]uint64{
					STAT_FORK: 1,
				},
			},
		},
		{
			events: []int{STAT_EXEC},
			expected: &Stats{
				Cmdline: cmdline,
				Total:   1,
				Count: [STAT_LAST]uint64{
					STAT_EXEC: 1,
				},
			},
		},
		{
			events: []int{STAT_EXIT},
			expected: &Stats{
				Cmdline: cmdline,
				Total:   1,
				Count: [STAT_LAST]uint64{
					STAT_EXIT: 1,
				},
			},
		},
		{
			events: []int{STAT_CORE},
			expected: &Stats{
				Cmdline: cmdline,
				Total:   1,
				Count: [STAT_LAST]uint64{
					STAT_CORE: 1,
				},
			},
		},
		{
			events: []int{STAT_COMM},
			expected: &Stats{
				Cmdline: cmdline,
				Total:   1,
				Count: [STAT_LAST]uint64{
					STAT_COMM: 1,
				},
			},
		},
		{
			events: []int{99},
			expected: &Stats{
				Cmdline: cmdline,
			},
		},
		{
			events: []int{STAT_FORK, STAT_EXEC},
			expected: &Stats{
				Cmdline: cmdline,
				Total:   2,
				Count: [STAT_LAST]uint64{
					STAT_FORK: 1,
					STAT_EXEC: 1,
				},
			},
		},
		{
			events: []int{STAT_FORK, STAT_FORK, STAT_EXEC},
			expected: &Stats{
				Cmdline: cmdline,
				Total:   3,
				Count: [STAT_LAST]uint64{
					STAT_FORK: 2,
					STAT_EXEC: 1,
				},
			},
		},
	}

	for _, testcase := range testcases {
		stats := NewStats(cmdline)
		if stats == nil {
			t.Errorf("Fail to make Stats")
			break
		}

		for _, event := range testcase.events {
			stats.Update(event)
		}

		if reflect.DeepEqual(stats, testcase.expected) == false {
			t.Errorf("expected %v\nbut%v\n", testcase.expected, stats)
		}
	}

}

func TestStatsMapList(t *testing.T) {
	self_cmdline, err := process.NewCmdlineByPid(os.Getpid())
	if err != nil {
		t.Errorf("fail to get cmdline of self")
		return
	}

	cmdline := self_cmdline.String()

	stats := NewStats(cmdline)
	if stats == nil {
		t.Errorf("Fail to make Stats")
		return
	}

	stats_map := StatsMap{
		cmdline: stats,
	}

	stats_list := stats_map.List()

	if len(stats_list) != 1 {
		t.Errorf("expected len is 1, but %d", len(stats_list))
		return
	}

	if reflect.DeepEqual(stats, stats_list[0]) == false {
		t.Errorf("expected\n%vbut\n%v\n", stats, stats_list[0])
	}
}

func TestSortByTotal(t *testing.T) {
	input := StatsList{
		&Stats{Total: 1},
		&Stats{Total: 9},
		&Stats{Total: 2},
		&Stats{Total: 8},
		&Stats{Total: 3},
		&Stats{Total: 7},
	}
	expected := StatsList{
		&Stats{Total: 1},
		&Stats{Total: 2},
		&Stats{Total: 3},
		&Stats{Total: 7},
		&Stats{Total: 8},
		&Stats{Total: 9},
	}

	sort.Sort(ByTotal(input))

	if reflect.DeepEqual(input, expected) == false {
		t.Errorf("expected\n%v\nbut\n%v\n", expected, input)
	}

}
