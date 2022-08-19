/*
 * Copyright 2021 Justin Reichardt
 *
 * This file is part of rhosts.
 *
 * rhosts is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * rhosts is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with rhosts.  If not, see <https://www.gnu.org/licenses/>.
 */

// rhosts - Program used to maintain a blocklist appended to a host file
package main

import (
	"flag"
	"fmt"
	"jbreich/rhosts/cfg"
	"jbreich/rhosts/hosts"
	"jbreich/rhosts/serve"
	sysos "jbreich/rhosts/sys"
	"log"
	"time"
)

var Exit chan bool

const GPL = `
    rhosts maintains a blocklist and appends it to the system hosts file

    Copyright (C) 2021  Justin Reichardt

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
`

func main() {
	tmpdir := ""
	hostsloc := ""
	cfgloc := ""
	var daemon bool = false
	var interval int = 1440
	var versionflag bool = false
	var removetimestamp bool = false

	// Parsing Flags
	flag.BoolVar(&daemon, "d", false, "Should this be run in daemon mode")
	flag.IntVar(&interval, "t", 1440, "Minutes until next run of daemon")
	flag.BoolVar(&versionflag, "version", false, "show version information")
	flag.BoolVar(&removetimestamp, "removetimestamp", false, "Removes the timestamp, used with logging programs to prevent a double timestamp")
	flag.Parse()

	// Display version information
	if versionflag {
		fmt.Print("Rhosts version: " + version)
		return
	}

	// Check if timestamp should be removed
	if removetimestamp {
		log.SetFlags(0)
	} else {
		// GPL information
		fmt.Println(GPL)
	}

	if daemon {
		log.Print("daemon:", daemon)
		log.Print("interval:", interval)
	}

	sysos.Detect(&tmpdir, &hostsloc, &cfgloc)

	// Read the config file
	config := cfg.Create(cfgloc)
	err, config := config.Update()
	if err != nil {
		log.Panic("Failed to parse config: " + cfgloc)
	}

	// Starting web server
	serve.Start("blank")

	// Update the hosts file
	if daemon == false {
		err := hosts.Update(config, tmpdir, hostsloc)
		if err != nil {
			log.Print(err)
		}
	} else {

		for true {
			err := hosts.Update(config, tmpdir, hostsloc)
			if err != nil {
				log.Print(err)
			}

			// Check if daemon
			if daemon == false {
				break
			}

			if err == nil {
				i := time.Now().Add(time.Duration(interval) * time.Minute).Format(time.Layout)
				log.Printf("Sleeping for %d minutes", interval)
				log.Print("Should restart at: " + i)
				time.Sleep(time.Duration(interval) * time.Minute)
			}
		}
	}
	<-Exit
}
