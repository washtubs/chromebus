package chromebus

import ()

type Aggregator struct {
	tabCache map[string]*ChromeTab
}

func (a *Aggregator) aggregate(record ChromebusRecord) {
	switch record.action {
	case string(New):
	case string(UrlChanged):
	case string(FocusChanged):
		a.tabCache[record.id] = record.newTab
	case string(Closed):
		delete(a.tabCache, record.id)
		// TODO: we need to update as well
	}
}

func (a *Aggregator) getIndexById(id string) (index int) {
	tab := a.tabCache[id]
	return tab.index
}
