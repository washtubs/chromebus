package chromebus

// must be a valid url. include the scheme, i.e. http://
const Redirect = "https://isitchristmas.com/"

// sample site hosts
var GoofHosts = [...]string{
	"www.homestarrunner.com",
	"www.facebook.com",
	"www.reddit.com",
}

const minutesBeforeBlock = 2
const minutesBeforeWarn = 1
const suspendMinutes = 5
const resetCron = "0 0 23 * * *" // 3pm every day
const reportIntervalSeconds = 5
