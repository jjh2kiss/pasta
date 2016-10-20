package process

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

//read from /proc/[pid]/stat
type ProcessStat struct {
	Pid                 int    //0
	Comm                string //1
	State               string //2
	PPid                int    //3
	UserTime            uint64 //13, device by sysconf(_SC_CLK_TCK)
	SystemTime          uint64 //14, device by sysconf(_SC_CLK_TCK)
	ChildUserTime       uint64 //15, device by sysconf(_SC_CLK_TCK)
	ChildSystemTime     uint64 //16, device by sysconf(_SC_CLK_TCK)
	Threads             int    //19
	StartTime           uint64 //21, device by sysconf(_SC_CLK_TCK)
	VmSize              uint64 //22
	VmRSS               uint64 //23
	DelayAcctBlkioTicks uint64 //41, delayacct_blkio_ticks, device by sysconf(_SC_CLK_TCK)
}

func NewProcessStat(reader io.Reader) (*ProcessStat, error) {
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if len(contents) == 0 {
		return nil, fmt.Errorf("stat is empty")
	}

	fields := strings.Fields(string(contents))
	if len(fields) == 0 {
		return nil, fmt.Errorf("fail to split")
	}

	stat := &ProcessStat{}

	stat.Pid, err = strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}

	i := 1
	for !strings.HasSuffix(fields[i], ")") {
		i++
	}

	stat.Comm = strings.Join(fields[1:i+1], " ")
	if strings.HasPrefix(stat.Comm, "(") {
		stat.Comm = strings.Replace(stat.Comm, "(", "", 1)
	}

	if strings.HasSuffix(stat.Comm, ")") {
		stat.Comm = strings.Replace(stat.Comm, ")", "", 1)
	}
	stat.Comm = strings.TrimSpace(stat.Comm)

	stat.State = fields[i+1]
	stat.PPid, err = strconv.Atoi(fields[i+2])
	if err != nil {
		return nil, err
	}

	stat.UserTime, err = strconv.ParseUint(fields[i+12], 10, 64)
	if err != nil {
		return nil, err
	}

	stat.SystemTime, err = strconv.ParseUint(fields[i+13], 10, 64)
	if err != nil {
		return nil, err
	}

	stat.ChildUserTime, err = strconv.ParseUint(fields[i+14], 10, 64)
	if err != nil {
		return nil, err
	}

	stat.ChildSystemTime, err = strconv.ParseUint(fields[i+15], 10, 64)
	if err != nil {
		return nil, err
	}

	stat.Threads, err = strconv.Atoi(fields[i+18])
	if err != nil {
		return nil, err
	}

	stat.StartTime, err = strconv.ParseUint(fields[i+20], 10, 64)
	if err != nil {
		return nil, err
	}

	stat.VmSize, err = strconv.ParseUint(fields[i+21], 10, 64)
	if err != nil {
		return nil, err
	}

	stat.VmRSS, err = strconv.ParseUint(fields[i+22], 10, 64)
	if err != nil {
		return nil, err
	}

	stat.DelayAcctBlkioTicks, err = strconv.ParseUint(fields[i+40], 10, 64)
	if err != nil {
		return nil, err
	}

	return stat, nil
}

func NewProcessStatByPid(pid int) (*ProcessStat, error) {
	filepath := path.Join("/proc", strconv.Itoa(pid), "stat")
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewProcessStat(f)
}

func (self *ProcessStat) Equal(other *ProcessStat) bool {
	return reflect.DeepEqual(self, other)
}

func (self *ProcessStat) StateString() string {
	return StateString(self.State)
}

func (self *ProcessStat) TotalCputime() uint64 {
	return self.UserTime + self.SystemTime
}
