// This manages all the system specific settings
package sys

import (
	"log"
)

var system struct {
	os string
	tmpdir string
	hostsloc string
	cfgloc string
}

func Detect (tmpdirp, hostslocp, cfglocp *string) {
	if (system.os == "") {
		log.Panic("This OS does not seem to be supported")
	}
	*tmpdirp = system.tmpdir
	*hostslocp = system.hostsloc
	*cfglocp = system.cfgloc
}
