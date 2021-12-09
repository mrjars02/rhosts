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
	"bufio"
	"log"
)

func main() {
	// Detect which OS is running
		//tmpdir := ""
		//hostsloc := ""
		cfgloc := ""
	switch runtime.GOOS {
	case "windows":
		fmt.Println("Windows is currently unsupported")
		os.Exit(1)
		//tmpdir := "C:\\tmp"
		//hostsloc := "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux":
		//tmpdir = "/tmp/"
		//hostsloc = "/etc/hosts"
		cfgloc ="/etc/rhosts/rhosts.cfg"
	case "ios":
		fmt.Println("ios")
		os.Exit(1)
	default:
		log.Print(runtime.GOOS," is not supported")
		os.Exit(1)
	}

	// Parse Config
	cfgparse(cfgloc)
	
}

func cfgparse (cfgloc string){
	log.Print("Opening: ", cfgloc)
	file, err := os.Open(cfgloc)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	filebuf := bufio.NewScanner(file)
	filebuf.Split(bufio.ScanLines)
	for res := filebuf.Scan();res;res = filebuf.Scan() {
		state, body := cfgparseline(filebuf.Text())
		log.Print(state, body)
	}
	err = filebuf.Err()
	if err != nil {
		log.Fatal(err)
	}
}
func cfgparseline(buf string) (uint8, string){
	// State options
	// 0 - Init
	// 1 - Error
	// 2 - Comment
	// 3 - Site
	// 4 - Download
	var state uint8= 0
	body :=buf[:]
	for i:=0; i<len(buf);i++ {
		//fmt.Printf("%c",buf[i])
		switch buf[i] {
		case ' ':
		case '#':
			state = 2
		case 'd':
			if (len(buf) < i+10) {
				state = 1
				break
			}
			if (buf[i:(i+9)] == "download=") {
				i +=9
				state = 4
				body = buf[i:]
			} else{
				state = 1
			}
		case 's':
			if (len(buf) < i+6) {
				state = 1
				break
			}
			if (buf[i:(i+5)] == "site=") {
				i +=5
				state = 3
				body = buf[i:]
			} else{
				state = 0
			}
			//compare buf[i:(i+3)] to "site"
		}
		if (state !=0){ 
			return state,body
		}
	}
	return state, body
}
