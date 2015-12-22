package chromebus

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

const SiteBlocker PluginSpec = "SiteBlocker"

var trackerHost string

var goofingBlocker bool = false
var blockerEngaged bool = false

var siteBlockerPlugin *Plugin = &Plugin{
	Init: func(input chan ChromebusRecord, aggregator Aggregator) {
		log.Printf("Started SiteBlocker plugin")
		for r := range input {
			aggregator.aggregate(r)
			if r.action != string(Closed) && aggregator.getTabById(r.id).focused { // very important that this is short circuited
				focusedTab := aggregator.getTabById(r.id)
				goofingBefore := goofingBlocker
				goofingBlocker = GoofingOff(focusedTab.url)
				if goofingBlocker != goofingBefore {
					report(goofingBlocker)
				}
				//if r.action == string(UrlChanged) {
				tab := aggregator.getTabById(r.id)
				if blockerEngaged && GoofingOff(tab.url) {
					err := Navigate(aggregator.getIndexById(r.id), Redirect)
					notifier.SendMessage("Quit goofin'. Redirecting to " + Redirect)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
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

func report(g bool) {
	var method string
	if g {
		method = "goofing"
	} else {
		method = "notgoofing"
	}
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
