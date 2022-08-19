//go:build linux
// +build linux

package cfg

func init() {
	fp := func(c *Config) {
		c.System.OS = "linux"
		c.System.TmpDir = "/tmp/"
		c.System.HostsLoc = "/etc/hosts"
		c.System.CfgLoc = "/etc/rhosts/"
		c.System.Var = "/var/lib/rhosts/"
	}

	configFuncs = append(configFuncs, fp)
}
