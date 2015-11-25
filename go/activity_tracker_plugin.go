package chromebus

import (
	"github.com/robfig/cron"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

const ActivityTracker PluginSpec = "ActivityTracker"

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

func issueChallenge() (passed bool) {
	challengeText := "You cannot step twice into the same river; for other waters are continually flowing in."
	cmd := exec.Command("zenity", "--text", challengeText, "--entry")
	out, err := cmd.Output()
	answer := strings.TrimSpace(string(out))
	passed = false
	if err == nil {
		if answer == challengeText {
			passed = true
		} else {
			log.Printf("strings did not match: %s vs %s", challengeText, answer)
		}
	} else {
		log.Printf("cancelled")
	}
	return
}

var plugin *Plugin = &Plugin{
	Init: func(input chan ChromebusRecord, aggregator Aggregator) {
		log.Printf("Started Activity plugin")

		cron := cron.New()
		cron.AddFunc(resetCron, resetLeeway)
		//cron.AddFunc("0 */5 * * * *", resetLeeway) // 3pm every day
		cron.Start()
		initLeeway()
		go monitor()
		for r := range input {
			aggregator.aggregate(r)
			if r.action != string(Closed) && aggregator.getTabById(r.id).focused { // very important that this is short circuited
				focusedTab := aggregator.getTabById(r.id)
				goofing = GoofingOff(focusedTab.url)
			}
		}
	},
	Handle: func(w http.ResponseWriter, r *http.Request) {
		log.Printf(r.URL.Path)
		switch r.URL.Path {
		case "/" + string(ActivityTracker) + "/suspend":
			log.Printf("suspending...")
			if leewayExpired {
				if issueChallenge() {
					leewayMutex.Lock()
					suspendEnabled = true
					suspendCount++
					leewayMutex.Unlock()
				}
			} else {
				log.Printf("leeway has not expired. this is unnecessary. aborting request")
			}
		default:
			log.Printf("not recognized... %s" + r.URL.Path)
		}
	},
	Cleanup: func() {
		if leewayExpired {
			notifier.SendMessage("You failed to stay under your max time.")
		}
	},
}

func shouldBlockSites() bool {
	return leewayExpired && !suspendEnabled
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
	notifier.SendMessage("I'm annoying")
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
			//notifier.SendMessage("Started goofing")
			log.Printf("Started goofing: time %d", t)
		} else if lastStartedGoofing != nil {
			latestDuration := time.Since(*lastStartedGoofing)
			currentTime := int64(totalDuration) + int64(latestDuration)
			if currentTime >= targetDuration() {
				leewayExpired = true
				suspendEnabled = false
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
