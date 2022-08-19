package cfg

import (
	"bufio"
	"log"
	"os"
	"strconv"
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
		Var string
	}
}

// Used to hold a list of functions to be run when building
var configFuncs []func(*Config)

// Create initialized a config to be used the entire session
func Create() (cfg Config) {
	for _,fp := range(configFuncs){
		fp(&cfg)
	}
	if (cfg.System.OS == ""){log.Fatal("Failed to detect the OS")}
	err, cfg := cfg.Update()

	if (err != nil){log.Fatal("Failed to handle the config file: " + err.Error())}
	return
}

// cfgparse recieves the location of the config file and returns a list of sites to add and content to download
func (cfg Config) Update() (error, Config) {
	// Opening the config file
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

	// Create one if it doesn't exist
	// This is done after so that it doesn't read the default one
	// in the event one doesn't exist
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
	for i , res := 0,filebuf.Scan(); res; res = filebuf.Scan() {
		i++
		buf := filebuf.Text()
		if (cfgparseline(buf, &cfg) == true){

			log.Fatal("Failed to read line: " + strconv.Itoa(i) + ": " + buf)
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
func cfgparseline(buf string, cfg *Config) (fail bool) {
	if len(buf) == 0 {
		return
	}
		switch buf[0] {
		case ' ':
		case '#':
			return
		case 'd':
			if (len(buf) > 10 && buf[0:9] == "download=") {
				cfg.Downloads = append(cfg.Downloads, buf[9:])
			} else {
				fail = true
			}
		case 's':
			if (len(buf) > 6 && buf[0:5] == "site=") {
				cfg.Sites = append(cfg.Sites, buf[5:])
			} else {
				fail = true
			}
		case 'w':
			if (len(buf) > 10 && buf[0:10] == "whitelist=") {
				cfg.Whitelist = append(cfg.Whitelist, buf[9:])
			} else {
				fail = true
			}
	}
	return
}
