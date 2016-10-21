package stats

import (
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/jjh2kiss/pstat/system/process"
)

func TestStatsNewStat(t *testing.T) {
	self_cmdline, err := process.NewCmdlineByPid(os.Getpid())
	if err != nil {
		t.Errorf("fail to get cmdline of self")
		return
	}

	testcases := []struct {
		cmdline  *process.Cmdline
		expected *Stats
	}{
		{
			cmdline:  nil,
			expected: nil,
		},
		{
			cmdline: self_cmdline,
			expected: &Stats{
				Cmdline: self_cmdline,
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

	testcases := []struct {
		events   []int
		expected *Stats
	}{
		{
			events: []int{STAT_FORK},
			expected: &Stats{
				Cmdline: self_cmdline,
				Total:   1,
				Count: [STAT_LAST]uint64{
					STAT_FORK: 1,
				},
			},
		},
		{
			events: []int{STAT_EXEC},
			expected: &Stats{
				Cmdline: self_cmdline,
				Total:   1,
				Count: [STAT_LAST]uint64{
					STAT_EXEC: 1,
				},
			},
		},
		{
			events: []int{STAT_EXIT},
			expected: &Stats{
				Cmdline: self_cmdline,
				Total:   1,
				Count: [STAT_LAST]uint64{
					STAT_EXIT: 1,
				},
			},
		},
		{
			events: []int{STAT_CORE},
			expected: &Stats{
				Cmdline: self_cmdline,
				Total:   1,
				Count: [STAT_LAST]uint64{
					STAT_CORE: 1,
				},
			},
		},
		{
			events: []int{STAT_COMM},
			expected: &Stats{
				Cmdline: self_cmdline,
				Total:   1,
				Count: [STAT_LAST]uint64{
					STAT_COMM: 1,
				},
			},
		},
		{
			events: []int{99},
			expected: &Stats{
				Cmdline: self_cmdline,
			},
		},
		{
			events: []int{STAT_FORK, STAT_EXEC},
			expected: &Stats{
				Cmdline: self_cmdline,
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
				Cmdline: self_cmdline,
				Total:   3,
				Count: [STAT_LAST]uint64{
					STAT_FORK: 2,
					STAT_EXEC: 1,
				},
			},
		},
	}

	for _, testcase := range testcases {
		stats := NewStats(self_cmdline)
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

	stats := NewStats(self_cmdline)
	if stats == nil {
		t.Errorf("Fail to make Stats")
		return
	}

	name := self_cmdline.String()

	stats_map := StatsMap{
		name: stats,
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
