package chromebus

import (
//"log"
)

type Aggregator struct {
	tabCache map[string]*ChromeTab
}

func (a *Aggregator) aggregate(record ChromebusRecord) {
	if record.action == string(Closed) {
		delete(a.tabCache, record.id)
	} else {
		a.tabCache[record.id] = record.newTab
	}
}

func (a *Aggregator) getTabById(id string) *ChromeTab {
	return a.tabCache[id]
}

func (a *Aggregator) getIndexById(id string) (index int) {
	tab := a.tabCache[id]
	return tab.index
}
