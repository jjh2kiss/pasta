# pstat
[![Build Status](https://travis-ci.org/jjh2kiss/pstat.png?branch=master)](https://travis-ci.org/jjh2kiss/pstat)  

process fork/exec/exit monitoring tool(go implementation of forkstat)

pstat is a programm that logs process fork, exec, exit, crashdump, comm activity
It is very useful for monitoring process behaviour and to track down processes

pstat uses the CN_PROC of Linux Netlink Connector to gather process activity
pstat may miss events if the system is overly busy
Netlink Connector requires root privilege.
pstat same as forkstat(http://kernel.ubuntu.com/~cking/forkstat/)

## Install
### compiler
```
$ cd $GOPATH
$ go get github.com/jjh2kiss/pstat
$ cd ./src/github/jjh2kiss/pstat
$ go build
$ sudo ./pstat
```

### binary
```
$ cd $GOPATH
$ git clone git@github.com:jjh2kiss/pstat.git
$ cd ./src/github/jjh2kiss/pstat/bin
$ sudo ./pstat
```

## pstat command line options:
  * -d, --dirstrip              strip off the directory path from the process name
  * -D value, --duration value  specify run duration in seconds (default: 0)
  * -e value, --event value     select which events to monitor(default: all)
  * -s, --shortname             show short process name information
  * -S, --statistics            show event statistics at end of the run
  * -q, --quiet                 run quietly and enable -S option
  * --help, -h                  show help
  * --version, -v               print the version

## Examples:

### monitoring all process event
```
sudo ./pstat -S
Time                Event  PID   Info Duration Process
2016/10/27 14:13:24 fork  1366 parent          sudo ./pstat -S
2016/10/27 14:13:24 fork  1373 child           ./pstat -S
2016/10/27 14:13:24 fork  1366 parent          sudo ./pstat -S
2016/10/27 14:13:24 fork  1374 child           ./pstat -S
2016/10/27 14:13:34 fork  2030 parent          /usr/lib/unity-settings-daemon/unity-settings-daemon
2016/10/27 14:13:34 fork  1375 child           /usr/lib/unity-settings-daemon/unity-settings-daemon
2016/10/27 14:13:34 fork  1375 parent          /usr/lib/unity-settings-daemon/unity-settings-daemon
2016/10/27 14:13:34 fork  1376 child           /usr/lib/unity-settings-daemon/unity-settings-daemon
2016/10/27 14:13:34 exit  1375      0    0.005 /usr/lib/unity-settings-daemon/unity-settings-daemon
2016/10/27 14:13:34 exec  1376                 /usr/bin/perl -w /usr/bin/x-terminal-emulator
2016/10/27 14:13:34 exec  1376                 /usr/bin/python3 /usr/bin/gnome-terminal
2016/10/27 14:13:34 fork  1813 parent          /sbin/upstart --user
2016/10/27 14:13:34 fork  1377 child           /usr/bin/python3 /usr/bin/gnome-terminal
2016/10/27 14:13:34 comm  1377                 /usr/bin/python3 /usr/bin/gnome-terminal -> gmain
2016/10/27 14:13:34 fork  1376 parent          /usr/bin/python3 /usr/bin/gnome-terminal
2016/10/27 14:13:34 fork  1378 child           /usr/bin/python3 /usr/bin/gnome-terminal
2016/10/27 14:13:34 exec  1378                 /usr/bin/gnome-terminal.real
2016/10/27 14:13:34 fork  1376 parent          /usr/bin/python3 /usr/bin/gnome-terminal
2016/10/27 14:13:34 fork  1379 child           /usr/bin/gnome-terminal.real
2016/10/27 14:13:34 comm  1379                 /usr/bin/gnome-terminal.real -> dconf worker
2016/10/27 14:13:34 fork  1376 parent          /usr/bin/python3 /usr/bin/gnome-terminal
2016/10/27 14:13:34 fork  1380 child           /usr/bin/gnome-terminal.real
2016/10/27 14:13:34 fork  1376 parent          /usr/bin/python3 /usr/bin/gnome-terminal
2016/10/27 14:13:34 fork  1381 child           /usr/bin/gnome-terminal.real
2016/10/27 14:13:34 comm  1381                 /usr/bin/gnome-terminal.real -> gdbus
2016/10/27 14:13:34 comm  1380                 /usr/bin/gnome-terminal.real -> gmain
2016/10/27 14:13:34 fork 24790 parent          /usr/lib/gnome-terminal/gnome-terminal-server
2016/10/27 14:13:34 fork  1382 child           /usr/lib/gnome-terminal/gnome-terminal-server
2016/10/27 14:13:34 exec  1382                 bash
2016/10/27 14:13:34 fork  1382 parent          bash
2016/10/27 14:13:34 fork  1383 child           bash
```

### monitoring coredump event

```
sudo ./pstat -e coredump -e exec -e exit | grep segfault
2016/10/27 14:49:41 exec  2707                 ./segfault
2016/10/27 14:49:43 core  2707                 ./segfault
2016/10/27 14:49:43 exit  2707    139    2.089 ./segfault
```
