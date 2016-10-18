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
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pstat"
	app.Usage = "monitoring fork(), exec(), exit()"
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

	app.Before = func(c *cli.Context) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("Need to run with root privilege.\n")
		}
		return nil
	}

	app.Action = func(c *cli.Context) {
	}

	app.After = func(c *cli.Context) error {
		return nil
	}

	app.RunAndExitOnError()
}
