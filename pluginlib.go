package chromebus

import (
	"net/http"
	"os/exec"
	"strconv"
)

type Plugin struct {
	Init    func(input chan ChromebusRecord, aggregator Aggregator)
	Handle  func(w http.ResponseWriter, r *http.Request)
	Cleanup func()
}

func Focus(id string) (e error) {
	e = exec.Command("node", Bin.activateTab, id).Run()
	return
}

func NewTab(url string) (e error) {
	e = exec.Command("node", Bin.newTab, url).Run()
	return
}

func CloseTab(id string) (e error) {
	e = exec.Command("node", Bin.closeTab, id).Run()
	return
}

func Navigate(index int, url string) (e error) {
	exec.Command("node", Bin.navigateTab, strconv.Itoa(index), url).Run()
	return
}
