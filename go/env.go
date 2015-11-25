package chromebus

import (
	"log"
	"os"
)

var ChromebusHome string // root install dir
//var Bin string           // all binaries
var Events string // events file
var Bin struct {
	activateTab string
	closeTab    string
	newTab      string
	navigateTab string
}

func EnvSetup() {
	ChromebusHome = os.Getenv("CHROMEBUS_HOME")
	if ChromebusHome == "" {
		log.Fatalf("environment variable CHROMEBUS_HOME not set")
	}
	if _, err := os.Stat(ChromebusHome); os.IsNotExist(err) {
		log.Fatalf("no such file or directory: %s", ChromebusHome)
	}

	bin := ChromebusHome + "/bin"
	Bin.activateTab = bin + "/activate-tab.js"
	Bin.closeTab = bin + "/close-tab.js"
	Bin.newTab = bin + "/new-tab.js"
	Bin.navigateTab = bin + "/navigate-tab.js"

	Events = ChromebusHome + "/events"
}
