package chromebus

import (
	"log"
	"os"
)

var ChromebusHome string // root install dir
var Bin string           // all binaries
var Events string        // events file

func EnvSetup() {
	ChromebusHome = os.Getenv("CHROMEBUS_HOME")
	if ChromebusHome == "" {
		log.Fatalf("environment variable CHROMEBUS_HOME not set")
	}
	if _, err := os.Stat(ChromebusHome); os.IsNotExist(err) {
		log.Fatalf("no such file or directory: %s", ChromebusHome)
	}
	Bin = ChromebusHome + "/bin"
	Events = ChromebusHome + "/events"
}
