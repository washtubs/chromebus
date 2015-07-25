package chromebus

import ()

var tabCache map[string]*ChromeTab

func aggregate(record *ChromebusRecord) {

	switch record.action {
	case string(New):
		tabCache[record.id] = record.newTab
	case string(Closed):
		delete(tabCache, record.id)
		// TODO: we need to update as well
	}
}
