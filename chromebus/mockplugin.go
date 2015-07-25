package chromebus

import (
	"log"
)

const Mock PluginSpec = "Mock"

var mockPlugin *Plugin = &Plugin{
	Init: func(input chan *ChromebusRecord) {
		log.Printf("Started MOCK plugin")
		//for r := range input {
		log.Printf("Mock %s", <-input)
		//}
	},
	Cleanup: func() {
		log.Printf("Cleaning up MOCK plugin")
	},
}

func init() {
	PluginRegistry[Mock] = mockPlugin
}
