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

	parts := bytes.Split(cmdline, []byte{0})
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

	name := ""

	if short == true {
		name = self.ShortString(1)
	} else {
		name = self.String()
	}

	if dirstrip == true {
		name = path.Base(name)
	}

	if kernelThread == true {
		name = "[" + name + "]"
	}

	return name

}
