package chromebus

import (
	"log"
)

const Mock PluginSpec = "Mock"

var mockPlugin *Plugin = &Plugin{
	Init: func(input chan ChromebusRecord, aggregator Aggregator) {
		log.Printf("Started MOCK plugin")
		for r := range input {
			aggregator.aggregate(r)
			log.Printf("Mock %s", r)
		}
	},
	Cleanup: func() {
		log.Printf("Cleaning up MOCK plugin")
	},
}

func init() {
	PluginRegistry[Mock] = mockPlugin
}
