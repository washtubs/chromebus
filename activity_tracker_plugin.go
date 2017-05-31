package chromebus

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/robfig/cron"
)

const ActivityTracker PluginSpec = "ActivityTracker"

// HTTP response codes
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
var autoCloseCountdown = 0

var plugin *Plugin = &Plugin{
	Init: func(input chan ChromebusRecord, aggregator Aggregator) {
		log.Printf("Started Activity plugin")

		cron := cron.New()
		cron.AddFunc(resetCron, resetLeeway)
		//cron.AddFunc("0 */5 * * * *", resetLeeway) // 3pm every day
		cron.Start()

		log.Printf("cron initialized with [%s]", resetCron)
		initLeeway()
		go monitor()
	},
	Handle: func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/" + string(ActivityTracker) + "/suspend":
			if leewayExpired {
				log.Printf("Suspending...")
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
		case "/" + string(ActivityTracker) + "/open":
			leewayMutex.Lock()
			c := r.URL.Query().Get("count")
			log.Printf("opening for " + c + " seconds")
			var err error
			autoCloseCountdown, err = strconv.Atoi(c)
			if err != nil {
				log.Printf("ERROR: " + c + " is not an integer")
			}
			leewayMutex.Unlock()
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
		case "/" + string(ActivityTracker) + "/logstatus":
			leewayMutex.Lock()
			if autoCloseCountdown == 0 {
				log.Printf("Site blocker engaged because autoCloseCountdown has reached 0")
			}
			if leewayExpired {
				if suspendEnabled {
					log.Printf("Leeway has expired (%dm/%dm) but suspend is enabled", totalDuration/time.Minute, minutesBeforeBlock)
				} else {
					log.Printf("Leeway has expired (%dm/%dm)", totalDuration/time.Minute, minutesBeforeBlock)
				}
			}
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
	// TODO: we are now ignoreing manualBlock to reduce complexity. Either re-add if desired, or remove all logic that surrounds it.
	// using autoCloseCountdown for everything instead
	return autoCloseCountdown == 0 || (leewayExpired && !suspendEnabled)
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
	log.Print("Resetting...")
	initLeeway()
}

func init() {
	PluginRegistry[ActivityTracker] = plugin
}

func monitor() {
	ticker = time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		leewayMutex.Lock()
		if autoCloseCountdown > 0 {
			autoCloseCountdown--
			if autoCloseCountdown == 0 {
				goofing = false
			}
		}
		if goofing && lastStartedGoofing == nil {
			lastStartedGoofing = new(time.Time)
			*lastStartedGoofing = t
			log.Printf("Started goofing: time %d", t)
		} else if lastStartedGoofing != nil {
			latestDuration := time.Since(*lastStartedGoofing)
			currentTime := uint64(totalDuration) + uint64(latestDuration)
			if currentTime >= targetDuration() {
				leewayExpired = true
				suspendEnabled = false
			} else if currentTime > uint64(time.Duration(minutesBeforeWarn*time.Minute)) && !wasWarned {
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

func targetDuration() uint64 {
	return uint64(time.Duration(uint64(minutesBeforeBlock+suspendMinutes*suspendCount) * uint64(time.Minute)))
}

func pushDuration(duration time.Duration) {
	totalDuration = time.Duration(uint64(totalDuration) + uint64(duration))
	timesSpentGoofing[index] = duration
	index++
}
