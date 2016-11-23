package stats

import (
	"os"
	"reflect"
	"testing"

	"github.com/jjh2kiss/pasta/config"
	"github.com/jjh2kiss/pasta/system/process"
)

func TestNewStatsTable(t *testing.T) {
	processes, err := process.NewProcessTable()
	if err != nil {
		t.Errorf("Fail to get Process table")
		return
	}

	table := NewStatsTable(processes, &config.Config{})

	if table == nil {
		t.Errorf("Fail to make Stats Table")
		return
	}
}

//table만 존재하는 경우
func TestNewStatsTableUpdate1(t *testing.T) {
	self_cmdline, err := process.NewCmdlineByPid(os.Getpid())
	if err != nil {
		t.Errorf("fail to get cmdline of self")
		return
	}

	cmdline := self_cmdline.String()

	expected := &Stats{
		Cmdline: cmdline,
		Total:   1,
		Count: [STAT_LAST]uint64{
			STAT_FORK: 1,
		},
	}

	processes, err := process.NewProcessTable()
	if err != nil {
		t.Errorf("Fail to get Process table")
		return
	}
	processes.Add(os.Getpid())

	table := NewStatsTable(processes, &config.Config{})

	if table == nil {
		t.Errorf("Fail to make Stats Table")
		return
	}

	table.Update(os.Getpid(), STAT_FORK)

	stats_list := table.List()

	if len(stats_list) != 1 {
		t.Errorf("stats_list should equal to 1, but %d", len(stats_list))
		return
	}

	if reflect.DeepEqual(expected, stats_list[0]) == false {
		t.Errorf("expected\n%v\nbut\n%v\n", expected, stats_list[0])
	}
}

//동일한 프로세스의 stats이 이미 등록되어 있는 경우
func TestNewStatsTableUpdate2(t *testing.T) {
	self_cmdline, err := process.NewCmdlineByPid(os.Getpid())
	if err != nil {
		t.Errorf("fail to get cmdline of self")
		return
	}

	cmdline := self_cmdline.String()

	expected := &Stats{
		Cmdline: cmdline,
		Total:   2,
		Count: [STAT_LAST]uint64{
			STAT_FORK: 2,
		},
	}

	processes, err := process.NewProcessTable()
	if err != nil {
		t.Errorf("Fail to get Process table")
		return
	}
	processes.Add(os.Getpid())

	table := NewStatsTable(processes, &config.Config{})

	if table == nil {
		t.Errorf("Fail to make Stats Table")
		return
	}

	table.Update(os.Getpid(), STAT_FORK)
	table.Update(os.Getpid(), STAT_FORK)

	stats_list := table.List()

	if len(stats_list) != 1 {
		t.Errorf("stats_list should equal to 1, but %d", len(stats_list))
		return
	}

	if reflect.DeepEqual(expected, stats_list[0]) == false {
		t.Errorf("expected\n%v\nbut\n%v\n", expected, stats_list[0])
	}
}
