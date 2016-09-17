package chromebus

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

const SiteBlocker PluginSpec = "SiteBlocker"

var trackerHost string

var goofingBlocker bool = false
var blockerEngaged bool = false

var siteBlockerPlugin *Plugin = &Plugin{
	Init: func(input chan ChromebusRecord, aggregator Aggregator) {
		log.Printf("Started SiteBlocker plugin")
		go trackCurrentTab(&aggregator)
		go reportRegularlyIfGoofing()
		for r := range input {
			aggregator.aggregate(r)
			//if r.action != string(Closed) && aggregator.getTabById(r.id).focused { // very important that this is short circuited
			//focusedTab := aggregator.getTabById(r.id)
			//goofingBefore := goofingBlocker
			//goofingBlocker = GoofingOff(focusedTab.url)
			//if goofingBlocker != goofingBefore {
			//report(goofingBlocker)
			//}
			////if r.action == string(UrlChanged) {
			//tab := aggregator.getTabById(r.id)
			//if blockerEngaged && GoofingOff(tab.url) {
			//err := Navigate(aggregator.getIndexById(r.id), Redirect)
			//notifier.SendMessage("Quit goofin'. Redirecting to " + Redirect)
			//if err != nil {
			//log.Fatal(err)
			//}
			//}
			//}
		}
	},
	Handle: func(w http.ResponseWriter, r *http.Request) {
		log.Printf(r.URL.Path)
		switch r.URL.Path {
		case "/" + string(SiteBlocker) + "/suspend":
			expiredR, expiredErr := http.Get(trackerUrlPrefix() + "/isleewayexpired")
			if expiredErr != nil {
				log.Fatal(expiredErr)
			}
			if expiredR.StatusCode == leewayIsExpired {
				if issueChallenge() {
					_, suspendErr := http.Get(trackerUrlPrefix() + "/suspend")
					if suspendErr != nil {
						log.Fatal(suspendErr)
					}
				}
			} else if expiredR.StatusCode == leewayIsNotExpired {
				log.Printf("leeway has not expired. this is unnecessary. aborting request")
			} else {
				log.Printf("ERROR: unexpected http status %s", expiredR.StatusCode)
			}
		}
	},
	Cleanup: func() {
	},
}

func trackerUrlPrefix() string {
	if trackerHost == "" {
		panic("trackerHost not specified!")
	}
	return "http://" + trackerHost + "/" + string(ActivityTracker)
}

func reportRegularlyIfGoofing() {
	ticker = time.NewTicker(reportIntervalSeconds * time.Second)
	for _ = range ticker.C {
		// NOTE: theoretically not thread safe, but not a big deal, just means
		// we may redundantly report that we are not goofing
		if goofingBlocker {
			log.Printf("Regular reporting goofing, will do so again in %d seconds", reportIntervalSeconds)
			reportAndBlock(goofingBlocker)
		}
	}
}

func trackCurrentTab(aggregator *Aggregator) {
	ticker = time.NewTicker(2 * time.Second)
	for _ = range ticker.C {
		focusedTab := aggregator.getFocusedTab()
		goofingBefore := goofingBlocker
		goofingBlocker = GoofingOff(focusedTab.url)
		if goofingBlocker != goofingBefore {
			reportAndBlock(goofingBlocker)
		}
		//if r.action == string(UrlChanged) {
		if blockerEngaged && goofingBlocker {
			err := Navigate(focusedTab.index, Redirect)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
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

func reportAndBlock(g bool) {
	report(goofingBlocker)
	if blockerEngaged && goofingBlocker {
		notifier.SendMessage("Quit goofin'. Redirecting to " + Redirect)
	}
}

func report(g bool) {
	var method string
	if g {
		method = "goofing"
	} else {
		method = "notgoofing"
	}
	log.Printf("Reporting %s", trackerUrlPrefix()+"/"+method)
	r, err := http.Get(trackerUrlPrefix() + "/" + method)
	if err != nil {
		log.Printf("ERROR: requesting tracker: %s", err)
	}
	switch r.StatusCode {
	case shouldBlockHttpStatus:
		blockerEngaged = true
	case shouldNotBlockHttpStatus:
		blockerEngaged = false
	default:
		log.Printf("ERROR: unexpected http status %s", r.StatusCode)
	}
}

func issueChallenge() (passed bool) {
	challengeText := "You cannot step twice into the same river; for other waters are continually flowing in."
	cmd := exec.Command("zenity", "--text", challengeText, "--entry")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("It seems zenity is not installed. so we cannot issue a proper challenge. Granting access.")
		passed = true
		return
	}
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

func init() {
	PluginRegistry[SiteBlocker] = siteBlockerPlugin
	flag.StringVar(&trackerHost, "tracker", "", "specify an activity tracker host")
}
