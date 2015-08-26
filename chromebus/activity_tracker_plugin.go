package chromebus

import (
	"log"
	"net/url"
	"strconv"
	"time"
)

const ActivityTracker PluginSpec = "ActivityTracker"

var goofing bool = false
var ticker *time.Ticker
var lastStartedGoofing *time.Time

// total time tracking
var timesSpentGoofing = new([1024]time.Duration)
var index = 0
var totalDuration time.Duration

var wasWarned = false

var plugin *Plugin = &Plugin{
	Init: func(input chan ChromebusRecord, aggregator Aggregator) {
		log.Printf("Started Activity plugin")
		go monitor()
		for r := range input {
			aggregator.aggregate(r)
			if r.action != string(Closed) && aggregator.getTabById(r.id).focused { // very important that this is short circuited
				focusedTab := aggregator.getTabById(r.id)
				goofing = GoofingOff(focusedTab.url)
			}
		}
	},
	Cleanup: func() {
		if ShouldBlockSites {
			notifier.SendMessage("You failed to stay under your max time.")
		}
	},
}

func init() {
	PluginRegistry[ActivityTracker] = plugin
}

func monitor() {
	ticker = time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		if goofing && lastStartedGoofing == nil {
			lastStartedGoofing = new(time.Time)
			*lastStartedGoofing = t
			//notifier.SendMessage("Started goofing")
			log.Printf("Started goofing: time %d", t)
		} else if lastStartedGoofing != nil {
			latestDuration := time.Since(*lastStartedGoofing)
			currentTime := int64(totalDuration) + int64(latestDuration)
			if currentTime >= int64(time.Duration(minutesBeforeBlock*time.Minute)) {
				ShouldBlockSites = true
			} else if currentTime > int64(time.Duration(minutesBeforeWarn*time.Minute)) && !wasWarned {
				notifier.SendMessage(strconv.Itoa(minutesBeforeBlock-minutesBeforeWarn) + " minutes till block")
				wasWarned = true
			}
			//log.Printf("Still goofing %d", latestDuration.String())
			if !goofing {
				lastStartedGoofing = nil
				pushDuration(latestDuration)
				log.Printf("Stopped goofing: goofing duration %fs", latestDuration.Seconds())
			}
		}
	}

}

func pushDuration(duration time.Duration) {
	totalDuration = time.Duration(int64(totalDuration) + int64(duration))
	timesSpentGoofing[index] = duration
	index++
}

func GoofingOff(urlRaw string) bool {
	url, err := url.Parse(urlRaw)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Host: " + url.Host)
	for _, candidate := range GoofHosts {
		if url.Host == candidate {
			return true
		}
	}
	return false
}
