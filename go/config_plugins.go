package chromebus

// must be a valid url. include the scheme, i.e. http://
const Redirect = "https://isitchristmas.com/"

// sample site hosts
var GoofHosts = [...]string{
	"www.homestarrunner.com",
	"www.facebook.com",
	"www.reddit.com",
}

const minutesBeforeBlock = 60
const minutesBeforeWarn = 5
const suspendMinutes = 5
const resetCron = "0 0 15 * * *"
const reportIntervalSeconds = 5
