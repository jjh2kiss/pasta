package system

/*
#include <unistd.h>
*/
import "C"

var _ClockTick uint64

func ClockTick() uint64 {
	if _ClockTick == 0 {
		var tick C.long
		tick = C.sysconf(C._SC_CLK_TCK)
		_ClockTick = uint64(tick)
	}
	return _ClockTick
}
