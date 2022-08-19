package hosts

import (
	"bufio"
	"io"
	"jbreich/rhosts/cfg"
	"log"
	"net/http"
	"os"
)

// siteList holds the location of all the sites along with a list of their location
type siteList struct {
	location  string
	siteEntry []siteEntry
}

// siteEntry holds a single entry and if it is a repeat
type siteEntry struct {
	repeat bool
	site   string
}

func Update() (err error) {
	err, config := cfg.Create()
	if (err != nil){return}
	var siteBuff []siteList

	err = error(nil)
	err = copystatichosts(config.System.TmpDir, config.System.HostsLoc)
	if err != nil {
		log.Print("Failed to copy static entries")
		return
	}
	defer os.Remove(config.System.TmpDir + "rhosts")
	err, siteBuff = downloadcontent(config.Downloads, config.System.TmpDir, config.System.HostsLoc)
	if err != nil {
		log.Print("Failed to download entries")
		return
	}
	err = writesites(config.Sites, config.System.TmpDir, &siteBuff)
	if err != nil {
		log.Print("Failed to failed to copy rhosts static entries")
		return
	}
	removeduplicates(&siteBuff, &config.Whitelist)
	err = write2tmp(config.System.TmpDir, &siteBuff)
	if err != nil {
		log.Print("Failed to write sites to tmpfile")
		return
	}
	err = writetmp2hosts(config.System.HostsLoc, config.System.TmpDir)
	if err != nil {
		log.Print("Failed to copy to hosts file")
		return
	}
	log.Print("Finished updating host")
	return
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
	for res := filebuf.Scan(); res; res = filebuf.Scan() {
		buff := filebuf.Text()
		if buff == "# rhosts begin" {
			break
		}
		_, err := file.WriteString(buff + "\n")
		if err != nil {
			log.Print(err)
			return err
		}
	}
	_, err = file.WriteString("# rhosts begin\n")
	err = filebuf.Err()
	return err
}

// downloadcontent attempts to download the provided url and create a siteList. If the file fails to download it attempts to find an old copy from the hosts file.
func downloadcontent(downloads []string, tmpdir string, hostsloc string) (err error, list []siteList) {
	for _, d := range downloads {
		var site siteList
		site.location = d
		log.Print("Downloading: ", d)
		response, err := http.Get(d)
		if err != nil {
			log.Print(err)
			log.Print("Looking for old record in hosts file")
			downloadoldlookup(hostsloc, d, &site)
		} else {
			defer response.Body.Close()
			scanner := bufio.NewScanner(response.Body)
			for scanner.Scan() {
				resp := checkDownloadLine(scanner.Text())
				if resp.site != "" {
					site.siteEntry = append(site.siteEntry, resp)
				}
			}
		}
		list = append(list, site)
	}
	return err, list
}

// checkDownloadLine parses the download line into just the address that needs to be blocked
func checkDownloadLine(line string) (address siteEntry) {
	var token []string
	address.repeat = false
	address.site = ""
	buff := ""
	lineLength := len(line) - 1
	for i, c := range line {
		if c != ' ' && i < lineLength {
			buff += string(c)
		} else if len(buff) > 0 {
			if i == lineLength {
				buff += string(c)
			}
			token = append(token, buff)
			buff = ""
		}
	}
	if len(token) == 0 {
		return
	}
	if token[0][0] == '#' {
		return
	}
	for _, t := range token {
		var period uint
		var failed bool
		period = 0
		failed = false
		for _, c := range t {
			switch c {
			case '.':
				period++
			case '#':
				return
			case ':':
				failed = true
				break
			}
		}
		if period <= 2 && failed == false {
			address.site = t
			return
		}
	}
	return
}

// downloadoldlookup attemps to find an entry in the hosts file and add it to the siteList
func downloadoldlookup(hostsloc, d string, site *siteList) error {
	var err error = nil
	var state uint8 = 0

	hostsf, err := os.Open(hostsloc)
	if err != nil {
		log.Print(err)
		return err
	}
	defer hostsf.Close()

	fbuff := bufio.NewScanner(hostsf)
	fbuff.Split(bufio.ScanLines)
	for res := fbuff.Scan(); res; res = fbuff.Scan() {
		buff := fbuff.Text()
		switch state {
		case 0:
			if buff == "# rhosts download - "+d {
				log.Print("Found old record in hosts file:" + buff)
				state = 1
			}
		case 1:
			if len(buff) >= 9 && buff[0:8] == "# rhosts" {
				state = 2
			} else {
				siteBuff := checkDownloadLine(buff)
				if siteBuff.site != "" {
					site.siteEntry = append(site.siteEntry, siteBuff)
				}
			}
		case 3:
			return nil
		}

	}

	return err
}

// writesites writes the list of sites from the config file to the local siteList
func writesites(sites []string, tmpdir string, siteBuff *[]siteList) (err error) {
	var localList siteList
	localList.location = "local"
	err = nil
	fileloc := tmpdir + "rhosts"
	log.Print("Opening: " + fileloc)
	if len(sites) == 0 {
		return
	}
	for _, s := range sites {
		var site siteEntry
		site.repeat = false
		site.site = s
		localList.siteEntry = append(localList.siteEntry, site)
	}
	*siteBuff = append(*siteBuff, localList)
	return
}

// removeduplicates removes any duplicate or uneeded/unwanted addresses
func removeduplicates(siteBuff *[]siteList, whitelist *[]string) {
	var safewords = []string{"localhost", "localhost.localdomain", "broadcasthost", "ip6-loopback", "ip6-localhost", "ip6-localnet", "ip6-mcastprefix", "ip6-allnodes", "ip6-allrouters", "ip6-allhosts", "local"}
	var c struct {
		d uint
		s uint
		w uint
	}
	c.d = 0
	c.s = 0
	c.w = 0
	log.Print("Checking for duplicates")
	var entry []struct {
		r *bool
		s *string
	}
	var entryBuff struct {
		r *bool
		s *string
	}
	for i := len((*siteBuff)) - 1; i > -1; i-- {
		for j := len((*siteBuff)[i].siteEntry) - 1; j > -1; j-- {
			entryBuff.r = &((*siteBuff)[i].siteEntry[j].repeat)
			entryBuff.s = &((*siteBuff)[i].siteEntry[j].site)
			entry = append(entry, entryBuff)
		}
	}
	lenEntry := len(entry)
	for i, e := range entry {
		for _, w := range safewords {
			if *e.s == w {
				*(entry[i].r) = true
				c.s++
				break
			}
		}
		if *(entry[i].r) == true {
			continue
		}
		for _, w := range *whitelist {
			if *e.s == w {
				*(entry[i].r) = true
				c.w++
				break
			}
		}
		if *(entry[i].r) == true {
			continue
		}
		if i == lenEntry {
			break
		}
		for j, n := range entry[i+1:] {
			if *e.s == *n.s {
				*(entry[i+j].r) = true
				c.d++
			}
		}

	}
	log.Printf("Total: %d\tDuplicates: %d\tSafeWords: %d\tWhitelisted: %d\n", lenEntry, c.d, c.s, c.w)
}

// write2tmp write the siteBuff to the tempfile
func write2tmp(tmpdir string, siteBuff *[]siteList) (err error) {
	err = nil
	tmploc := tmpdir + "rhosts"
	tmpf, err := os.OpenFile(tmploc, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer tmpf.Close()
	if err != nil {
		log.Print(err)
		return err
	}
	for _, location := range *siteBuff {
		if len(location.siteEntry) == 0 {
			continue
		}
		_, err := tmpf.WriteString("# rhosts download - " + location.location + "\n")
		if err != nil {
			return err
		}
		for _, site := range location.siteEntry {
			if site.repeat == false {
				_, err = tmpf.WriteString("0.0.0.0 " + site.site + "\n")
				if err != nil {
					return err
				}
				_, err = tmpf.WriteString(":: " + site.site + "\n")
				if err != nil {
					return err
				}
			}
		}
	}
	return
}

// writetmp2hosts overwrites the hostsfile with the tmp file
func writetmp2hosts(hostsloc, tmpdir string) error {
	var err error = nil
	tmploc := tmpdir + "rhosts"

	hosts, err := os.Create(hostsloc)
	if err != nil {
		log.Print(err)
		return err
	}
	tmp, err := os.Open(tmploc)
	if err != nil {
		log.Print(err)
		return err
	}
	_, err = io.Copy(hosts, tmp)
	if err != nil {
		log.Print(err)
	}

	return err
}
