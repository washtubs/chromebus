package chromebus

import (
	"sync"
	//"log"
)

type Aggregator struct {
	tabCache   map[string]*ChromeTab
	focusedKey string // shortcut to get the currently focused tab
	mutex      *sync.Mutex
}

func (a *Aggregator) aggregate(record ChromebusRecord) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if record.action == string(Closed) {
		delete(a.tabCache, record.id)
	} else {
		tab := record.newTab
		a.tabCache[record.id] = tab
		if tab.focused {
			a.focusedKey = record.id
		}
	}
}

func (a *Aggregator) getTabById(id string) *ChromeTab {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.tabCache[id]
}

func (a *Aggregator) getFocusedTab() *ChromeTab {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.tabCache[a.focusedKey]
}

func (a *Aggregator) getIndexById(id string) (index int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	tab := a.tabCache[id]
	return tab.index
}
