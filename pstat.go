/*
 * Copyright (C) 2016 Jehyun.Jeon
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.
 *
 * Written by Jehyun.jeon <jjh2kiss@gmail.com>
 *
 * Some of this code originally derived from forkstat
 * http://kernel.ubuntu.com/~cking/forkstat/
 *
 */

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"

	"github.com/jjh2kiss/pstat/config"
	"github.com/jjh2kiss/pstat/monitor"
	"github.com/jjh2kiss/pstat/utils"
	"golang.org/x/net/context"
)

func main() {
	app := cli.NewApp()
	app.Name = "pstat"
	app.Usage = "monitoring process event, such as fork(), exec(), exit()"
	app.Version = "0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "jerry",
			Email: "jjh2kiss@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "d, dirstrip",
			Usage: "strip off the directory path from the process name",
		},
		cli.Int64Flag{
			Name:  "D, duration",
			Usage: "specify run duration in seconds",
		},
		cli.StringSliceFlag{
			Name:  "e, event",
			Usage: "select which events to monitor",
		},
		cli.BoolFlag{
			Name:  "s, shortname",
			Usage: "show short process name information",
		},
		cli.BoolFlag{
			Name:  "S, statistics",
			Usage: "show event statistics at end of the run",
		},
		cli.BoolFlag{
			Name:  "q, quiet",
			Usage: "run quietly and enable -S option",
		},
	}

	config := config.Config{}

	app.Before = func(c *cli.Context) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("Need to run with root privilege.\n")
		}

		config.Dirstrip = c.Bool("dirstrip")
		config.Duration = c.Int64("duration")
		config.Shortname = c.Bool("shortname")
		config.Statistics = c.Bool("statistics")
		config.Quiet = c.Bool("quiet")

		names := c.StringSlice("event")
		if len(names) > 0 {
			for _, name := range names {
				event := utils.ProcEventUint32(name)
				config.Events |= event
			}
		} else {
			config.Events = 0xffffffff
		}

		return nil
	}

	app.Action = func(c *cli.Context) {
		ctx, cancel := context.WithCancel(context.Background())

		//signal handler
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			cancel()
		}()

		//timeout
		if config.Duration > 0 {
			go func() {
				select {
				case <-time.After(time.Duration(config.Duration) * time.Second):
					cancel()
				}
			}()
		}

		//run monitor
		err := monitor.Monitor(&config, ctx.Done())
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	app.After = func(c *cli.Context) error {
		return nil
	}

	app.RunAndExitOnError()
}
