package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var _BootTime uint64

func bootTime() (uint64, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "btime") {
			fields := strings.Fields(line)
			if len(fields) != 2 {
				return 0, fmt.Errorf("wroing btime format")
			}

			btime, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				return 0, err
			}
			return uint64(btime), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("could not found btime")
}

func BootTime() uint64 {
	if _BootTime == 0 {
		boottime, err := bootTime()
		if err != nil {
			boottime = 0
		}
		_BootTime = boottime
	}
	return _BootTime
}
