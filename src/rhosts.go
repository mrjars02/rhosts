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
	"runtime"
	"os"
	"io"
	"bufio"
	"log"
	"net/http"
	"flag"
	"time"
)

func main() {
	tmpdir := ""
	hostsloc := ""
	cfgloc := ""
	var daemon bool=false
	var interval int=1440

	// Parsing Flags
	flag.BoolVar(&daemon, "d", false, "Should this be run in daemon mode")
	flag.IntVar(&interval, "t", 1440, "Minutes until next run of daemon")
	flag.Parse()
	log.Print("daemon:" , daemon)
	log.Print("interval:",interval)

	sysdetect (&tmpdir, &hostsloc, &cfgloc)

	for true {
		err := error(nil)
		sites, downloads, err := cfgparse(cfgloc)
		if (err != nil){
			log.Print("Failed to parse config file")
			continue
		}
		err = copystatichosts(tmpdir, hostsloc)
		if (err != nil){
			log.Print("Failed to copy static entries")
			continue
		}
		defer os.Remove(tmpdir + "rhosts")
		err = downloadcontent(downloads, tmpdir, hostsloc)
		if (err != nil){
			log.Print("Failed to download entries")
			continue
		}
		err = writesites(sites, tmpdir)
		if (err != nil){
			log.Print("Failed to failed to copy rhosts static entries")
			continue
		}
		err = writetmp2hosts(hostsloc, tmpdir)
		if (err != nil){
			log.Print("Failed to copy to hosts file")
			continue
		}
		log.Print("Finished updating host")
		if (daemon == true){
			time.Sleep(time.Duration(interval) * time.Minute)
		}else{
			break
		}
	}
}

// sysdetect determines which OS it is running on and set the default locations
func sysdetect (tmpdir, hostsloc, cfgloc *string) {
	// Detect OS and set params
	switch runtime.GOOS {
	case "windows":
		log.Fatal("Windows is not supported")
		*tmpdir = "C:\\tmp"
		*hostsloc = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux":
		*tmpdir = "/tmp/"
		*hostsloc = "/etc/hosts"
		*cfgloc ="/etc/rhosts/rhosts.cfg"
	case "ios":
		log.Fatal("IOS is not supported")
	default:
		log.Fatal(runtime.GOOS," is not supported")
	}
}

// cfgparse recieves the location of the config file and returns a list of sites to add and content to download
func cfgparse (cfgloc string) ([]string, []string, error){
	var err error=nil
	var downloads []string
	var sites []string
	log.Print("Opening: ", cfgloc)
	file, err := os.Open(cfgloc)
	defer file.Close()
	if err != nil {
		log.Print(err)
		return nil, nil,err
	}
	filebuf := bufio.NewScanner(file)
	filebuf.Split(bufio.ScanLines)
	for res := filebuf.Scan();res;res = filebuf.Scan() {
		state, body := cfgparseline(filebuf.Text())
		switch state {
		case 3:
			sites =append(sites,body)
		case 4:
			downloads = append(downloads,body)
	}
	}
	err = filebuf.Err()
	if err != nil {
		log.Print(err)
		return nil, nil, err
	}
	return sites, downloads, err
}

// cfgparseline reads a single line of the config and returns the type and content of the line
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

// copystatichosts copies the hosts not managed by rhosts from the hosts file to the start of the temporary file
func copystatichosts(tmpdir, hostsloc string) error {
	fileloc := tmpdir + "rhosts"
	file, err := os.Create(fileloc)
	defer file.Close()
	if err != nil {
		log.Print(err)
		return err
	}
	filer, err := os.Open(hostsloc)
	defer filer.Close()
	if err != nil {
		log.Print(err)
		return err
	}
	filebuf := bufio.NewScanner(filer)
	filebuf.Split(bufio.ScanLines)
	for res := filebuf.Scan();res;res = filebuf.Scan() {
		buff := filebuf.Text()
		if (buff == "# rhosts begin"){
			break
		}
		_,err := file.WriteString(buff + "\n")
		if (err != nil) {
			log.Print(err)
			return err
		}
	}
	_,err = file.WriteString("# rhosts begin\n")
	err = filebuf.Err()
	return err
}

// downloadcontent attempts to download the provided url to the temp file. If the file fails to download it attempts to find an old copy from the hosts file.
func downloadcontent(downloads []string, tmpdir string, hostsloc string) error{
	var err error=nil
	fileloc := tmpdir + "rhosts"
	log.Print("Opening: ", fileloc)
	file,err := os.OpenFile(fileloc, os.O_APPEND|os.O_WRONLY, 0644)
	if (err != nil) {
		log.Print(err)
		return err
	}
	defer file.Close()

	for _, d := range downloads {
		log.Print("Downloading: ",d)
		file.WriteString("# rhosts download - " + d + "\n")
		if (err != nil) {
		}
		response, err := http.Get(d)
		if (err !=nil) {
			log.Print(err)
			log.Print("Looking for old record in hosts file")
			downloadoldlookup(file, hostsloc, d)
			continue
		}else{
			_,err := io.Copy(file,response.Body)
			if (err != nil){
				log.Print(err)
				return err
			}
		}
		defer response.Body.Close()
		file.WriteString("\n")
		if (err != nil) {
			log.Print(err)
			return err
		}
	}
	return err
}

// downloadoldlookup attemps to find an entry in the hosts file and write it to the temp file.
func downloadoldlookup(file *os.File, hostsloc, d string) error {
	var err error = nil
	var state uint8 = 0
	
	hostsf, err := os.Open(hostsloc)
	if (err != nil){
		log.Print(err)
		return err
	}
	defer hostsf.Close()

	fbuff := bufio.NewScanner(hostsf)
	fbuff.Split(bufio.ScanLines)
	for res := fbuff.Scan();res;res = fbuff.Scan() {
		buff := fbuff.Text()
		switch state {
		case 0:
			if (buff == "# rhosts download - " + d){
				log.Print("Found old record in hosts file:" + buff)
				state =1
			}
		case 1:
			if (len(buff) >=9 && buff[0:8] == "# rhosts"){
				state = 2
			}else{
				_,err := file.WriteString(buff + "\n")
				if (err != nil) {
					log.Print(err)
					return err
				}
			}
		case 3:
			return nil
		}
			
	}

	return err
}

// writesites writes the list of sites from the config file to the temp file
func writesites(sites []string, tmpdir string) error {
	var err error = nil
	fileloc := tmpdir + "rhosts"
	log.Print("Opening: " + fileloc)
	file,err := os.OpenFile(fileloc, os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()
	if (err != nil) {
		log.Print(err)
		return err
	}
	_,err = file.WriteString("# rhosts sites\n")
	if (err != nil){
		log.Print(err)
		return err
	}
	for _,s := range sites {
		_,err = file.WriteString(s)
		if (err != nil){
			log.Print(err)
			break
		}
	}
	return err
}

// writetmp2hosts overwrites the hostsfile with the tmp file
func writetmp2hosts(hostsloc, tmpdir string) error {
	var err error = nil
	tmploc := tmpdir + "rhosts"

	hosts, err := os.Create(hostsloc)
	if (err != nil){
		log.Print(err)
		return err
	}
	tmp, err := os.Open(tmploc)
	if (err != nil){
		log.Print(err)
		return err
	}
	_,err = io.Copy(hosts,tmp)
	if (err != nil){
		log.Print(err)
	}


	return err
}
