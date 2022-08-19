//go:build windows
// +build windows

package sys

func init() {

}

func init() {
	fp := func(c *Config) {
		c.System.os = "windows"
		c.System.tmpdir = "/tmp/"
		c.System.hostsloc = "/Windows/System32/drivers/etc/hosts"
		c.System.cfgloc = "/ProgramData/rhosts/"
		c.System.Var = "/ProgramData/rhosts/"
	}

	configFuncs = append(configFuncs, fp)
}
