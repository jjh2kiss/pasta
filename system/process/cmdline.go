package process

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

//read from /proc/[pid]/cmdline
type Cmdline struct {
	slice []string
}

func NewCmdline(reader io.Reader) (*Cmdline, error) {
	cmdline, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	cmdline = bytes.TrimSpace(cmdline)

	if len(cmdline) == 0 {
		return &Cmdline{}, nil
	}

	if cmdline[len(cmdline)-1] == 0 {
		cmdline = cmdline[:len(cmdline)-1]
	}

	if len(cmdline) == 0 {
		return &Cmdline{}, nil
	}

	//issue #8 google chrome 는 args가 space를 이용해 구분되어 있음
	// 널문자를 포함하는 경우 null로 구분, 그렇지 않다면 sp로 구분한다.
	sep := []byte{32}
	if bytes.Contains(cmdline, []byte{0}) {
		sep = []byte{0}
	}

	parts := bytes.Split(cmdline, sep)
	strParts := []string{}
	for _, p := range parts {
		strParts = append(strParts, string(p))
	}

	return &Cmdline{
		slice: strParts,
	}, nil
}

func NewCmdlineByPid(pid int) (*Cmdline, error) {
	filepath := path.Join("/proc", strconv.Itoa(pid), "cmdline")
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewCmdline(f)
}

func (self *Cmdline) String() string {
	return strings.Join(self.slice, " ")
}

func (self *Cmdline) ShortString(n int) string {
	if n > len(self.slice) {
		n = len(self.slice)
	}

	return strings.Join(self.slice[:n], " ")
}

func (self *Cmdline) Slice() []string {
	return self.slice
}

func (self *Cmdline) CombinedString(kernelThread bool, short bool, dirstrip bool) string {
	if self == nil {
		return ""
	}

	if len(self.slice) == 0 {
		return ""
	}

	name := self.slice[0]

	if dirstrip == true && name != "" {
		name = path.Base(name)
	}

	if kernelThread == true {
		name = "[" + name + "]"
	}

	if short == false && len(self.slice) > 1 {
		args := strings.Join(self.slice[1:], " ")
		name = name + " " + args
	}

	return name

}
