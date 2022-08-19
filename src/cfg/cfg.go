package cfg

import (
	"bufio"
	"log"
	"os"
	"errors"
)

const CFG = `
# There are 3 types of entries: download, site, and whitelist. Downloads are
# downloaded and stripped of comments and bad entries if possible before being
# added to a list of sites. Whitelisted urls are removed from the list of sites.
# From there all the urls are added to the hosts file for both IPv4 and IPv6.
# You can also add comments by prepending with a '#'.

# This is a static entry
#site=www.site.xyz
# This is a download entry
#download=w3.site.xyz/location/to/config.txt
# This is a whitelist entry
#whitelist=www.site.xyz

# A suggested download is: https://github.com/StevenBlack/hosts
#download=https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts`

type Config struct {
	CfgLoc    string
	Sites     []string
	Downloads []string
	Whitelist []string
	System struct {
		OS string
		TmpDir string
		HostsLoc string
		CfgLoc string
	}
}

// Used to hold a list of functions to be run when building
var configFuncs []func(*Config)

// Create initialized a config to be used the entire session
func Create() (err error,cfg Config) {
	for _,fp := range(configFuncs){
		fp(&cfg)
	}
	if (cfg.System.OS == ""){return errors.New("Failed to detect the OS"), cfg}
	err, cfg = cfg.Update()
	return
}

// cfgparse recieves the location of the config file and returns a list of sites to add and content to download
func (cfg Config) Update() (error, Config) {
	l := (cfg.System.CfgLoc + "rhosts.cfg")
	var err error = nil
	log.Print("Opening: ", l)
	if _, err = os.Stat(cfg.System.CfgLoc); os.IsNotExist(err) {
		log.Print(cfg.System.CfgLoc + " Does not exist, attempting to create it")
		err = os.MkdirAll(cfg.System.CfgLoc, 0755)
		if err != nil {
			log.Fatal("Could not create " + cfg.System.CfgLoc)
		}
	}
	if _, err = os.Stat(l); os.IsNotExist(err) {
		log.Print(l + " does not exist, attempting to create a placeholder")
		err = os.WriteFile(l, []byte(CFG), 0644)
		if err != nil {
			log.Fatal("Unable to create file: " + l)
		}
	}
	file, err := os.Open(l)
	defer file.Close()
	if err != nil {
		log.Print(err)
		return err, cfg
	}
	filebuf := bufio.NewScanner(file)
	filebuf.Split(bufio.ScanLines)
	for res := filebuf.Scan(); res; res = filebuf.Scan() {
		state, body := cfgparseline(filebuf.Text())
		switch state {
		case 3:
			cfg.Sites = append(cfg.Sites, body)
		case 4:
			cfg.Downloads = append(cfg.Downloads, body)
		case 5:
			cfg.Whitelist = append(cfg.Whitelist, body)
		}
	}
	err = filebuf.Err()
	if err != nil {
		log.Print(err)
		return err, cfg
	}
	return err, cfg
}

// cfgparseline reads a single line of the config and returns the type and content of the line
func cfgparseline(buf string) (uint8, string) {
	// State options
	// 0 - Init
	// 1 - Error
	// 2 - Comment
	// 3 - Site
	// 4 - Download
	// 5 - Whitelist
	var state uint8 = 0
	body := buf[:]
	for i := 0; i < len(buf); i++ {
		//fmt.Printf("%c",buf[i])
		switch buf[i] {
		case ' ':
		case '#':
			state = 2
		case 'd':
			if len(buf) < i+10 {
				state = 1
				break
			}
			if buf[i:(i+9)] == "download=" {
				i += 9
				state = 4
				body = buf[i:]
			} else {
				state = 1
			}
		case 's':
			if len(buf) < i+6 {
				state = 1
				break
			}
			if buf[i:(i+5)] == "site=" {
				i += 5
				state = 3
				body = buf[i:]
			} else {
				state = 0
			}
			//compare buf[i:(i+3)] to "site"
		case 'w':
			if len(buf) < i+10 {
				state = 1
				break
			}
			if buf[i:(i+10)] == "whitelist=" {
				i += 10
				state = 5
				body = buf[i:]
			} else {
				state = 1
			}
		}
		if state != 0 {
			return state, body
		}
	}
	return state, body
}
