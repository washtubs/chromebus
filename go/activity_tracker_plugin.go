package chromebus

import (
	"github.com/robfig/cron"
	"log"
	"net/http"
	"sync"
	"time"
)

const ActivityTracker PluginSpec = "ActivityTracker"

const shouldBlockHttpStatus = 202
const shouldNotBlockHttpStatus = 203
const leewayIsExpired = 204
const leewayIsNotExpired = 205

var goofing bool = false
var ticker *time.Ticker
var lastStartedGoofing *time.Time

// total time tracking
var timesSpentGoofing *[1024]time.Duration
var index int
var totalDuration time.Duration
var wasWarned = false
var leewayExpired = false
var leewayMutex = new(sync.Mutex)
var suspendCount = 0
var suspendEnabled = false
var manualBlocked = false

var plugin *Plugin = &Plugin{
	Init: func(input chan ChromebusRecord, aggregator Aggregator) {
		log.Printf("Started Activity plugin")

		cron := cron.New()
		cron.AddFunc(resetCron, resetLeeway)
		//cron.AddFunc("0 */5 * * * *", resetLeeway) // 3pm every day
		cron.Start()
		initLeeway()
		go monitor()
	},
	Handle: func(w http.ResponseWriter, r *http.Request) {
		log.Printf(r.URL.Path)
		switch r.URL.Path {
		case "/" + string(ActivityTracker) + "/suspend":
			if leewayExpired {
				log.Printf("suspending...")
				leewayMutex.Lock()
				suspendEnabled = true
				suspendCount++
				leewayMutex.Unlock()
			}
		case "/" + string(ActivityTracker) + "/isleewayexpired":
			if leewayExpired {
				w.WriteHeader(leewayIsExpired)
			} else {
				w.WriteHeader(leewayIsNotExpired)
			}
		case "/" + string(ActivityTracker) + "/block":
			leewayMutex.Lock()
			manualBlocked = true
			leewayMutex.Unlock()
		case "/" + string(ActivityTracker) + "/unblock":
			leewayMutex.Lock()
			manualBlocked = false
			leewayMutex.Unlock()
		case "/" + string(ActivityTracker) + "/goofing":
			leewayMutex.Lock()
			goofing = true
			if shouldBlockSites() {
				w.WriteHeader(shouldBlockHttpStatus)
			} else {
				w.WriteHeader(shouldNotBlockHttpStatus)
			}
			leewayMutex.Unlock()
		case "/" + string(ActivityTracker) + "/notgoofing":
			leewayMutex.Lock()
			goofing = false
			w.WriteHeader(202)
			leewayMutex.Unlock()
		default:
			log.Printf("not recognized... %s" + r.URL.Path)
		}
	},
	Cleanup: func() {
		if leewayExpired {
			// TODO: pushbullet. this cant work on the server otherwise
			//notifier.SendMessage("You failed to stay under your max time.")
		}
	},
}

func shouldBlockSites() bool {
	return manualBlocked || (leewayExpired && !suspendEnabled)
}

func initLeeway() {
	leewayMutex.Lock()
	defer leewayMutex.Unlock()
	timesSpentGoofing = new([1024]time.Duration)
	totalDuration = time.Duration(0)
	lastStartedGoofing = nil
	index = 0
	leewayExpired = false
	wasWarned = false
}

func resetLeeway() {
	initLeeway()
}

func init() {
	PluginRegistry[ActivityTracker] = plugin
}

func monitor() {
	ticker = time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		leewayMutex.Lock()
		if goofing && lastStartedGoofing == nil {
			lastStartedGoofing = new(time.Time)
			*lastStartedGoofing = t
			log.Printf("Started goofing: time %d", t)
		} else if lastStartedGoofing != nil {
			latestDuration := time.Since(*lastStartedGoofing)
			currentTime := int64(totalDuration) + int64(latestDuration)
			if currentTime >= targetDuration() {
				leewayExpired = true
				suspendEnabled = false
			} else if currentTime > int64(time.Duration(minutesBeforeWarn*time.Minute)) && !wasWarned {
				// TODO: pushbullet. this cant work on the server otherwise
				//notifier.SendMessage(strconv.Itoa(minutesBeforeBlock-minutesBeforeWarn) + " minutes till block")
				wasWarned = true
			}
			//log.Printf("Still goofing %d", latestDuration.String())
			if !goofing {
				lastStartedGoofing = nil
				pushDuration(latestDuration)
				log.Printf("Stopped goofing: goofing duration %fs", latestDuration.Seconds())
			}
		}
		leewayMutex.Unlock()
	}
}

func targetDuration() int64 {
	return int64(time.Duration((minutesBeforeBlock + suspendMinutes*suspendCount) * int(time.Minute)))
}

func pushDuration(duration time.Duration) {
	totalDuration = time.Duration(int64(totalDuration) + int64(duration))
	timesSpentGoofing[index] = duration
	index++
}
