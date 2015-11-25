package chromebus

var AggregatorModel = Aggregator{
	tabCache: make(map[string]*ChromeTab),
}

var notifier Notifier = new(NotifySendNotifier)

func GetNotifier() Notifier {
	return notifier
}
