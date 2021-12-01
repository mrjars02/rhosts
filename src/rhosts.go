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


// rhosts - Program used to maintain a blocklist within a hostfile
package main

import (
	"runtime"
	"fmt"
	"os"
)

func main() {
	// Detect which OS is running
	switch runtime.GOOS {
	case "windows":
		fmt.Println("Windows is currently unsupported")
		os.Exit(1)
		//tmpdir := "C:\\tmp"
		//hostsloc := "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux":
		tmpdir := "/tmp/"
		hostsloc := "/etc/hosts"
		cfgloc :="/etc/rhosts/rhosts.cfg"
		fmt.Println("linux:",tmpdir,hostsloc,cfgloc)
	case "ios":
		fmt.Println("ios")
	default:
		fmt.Println(runtime.GOOS," is not supported")
		os.Exit(1)
	}

	// Parse Config
}
