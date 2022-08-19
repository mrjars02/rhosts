// +build linux

package sys


func init(){

	system.os  = "linux"
	system.tmpdir = "/tmp/"
	system.hostsloc = "/etc/hosts"
	system.cfgloc = "/etc/rhosts/"
}

