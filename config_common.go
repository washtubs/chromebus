package chromebus

import (
	"sync"
)

var AggregatorModel = Aggregator{
	tabCache:   make(map[string]*ChromeTab),
	focusedKey: "",
	mutex:      new(sync.Mutex),
}

var notifier Notifier = new(NotifySendNotifier)

func GetNotifier() Notifier {
	return notifier
}
