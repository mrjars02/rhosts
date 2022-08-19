// +build linux

package cfg


func init(){
	fp := func(c *Config){
			c.System.OS  = "linux"
			c.System.TmpDir = "/tmp/"
			c.System.HostsLoc = "/etc/hosts"
			c.System.CfgLoc = "/etc/rhosts/"
	}

	configFuncs = append(configFuncs,fp)
}

