package chromebus

import (
	"log"
)

const SiteBlocker PluginSpec = "SiteBlocker"

var ShouldBlockSites bool = false

var siteBlockerPlugin *Plugin = &Plugin{
	Init: func(input chan ChromebusRecord, aggregator Aggregator) {
		log.Printf("Started SiteBlocker plugin")
		for r := range input {
			aggregator.aggregate(r)
			if r.action != string(Closed) && aggregator.getTabById(r.id).focused { // very important that this is short circuited
				//if r.action == string(UrlChanged) {
				tab := aggregator.getTabById(r.id)
				if ShouldBlockSites && GoofingOff(tab.url) {
					err := Navigate(aggregator.getIndexById(r.id), Redirect)
					notifier.SendMessage("Quit goofin'. Redirecting to " + Redirect)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	},
	Cleanup: func() {
	},
}

func init() {
	PluginRegistry[SiteBlocker] = siteBlockerPlugin
}
