package chromebus

import (
	//"log"
	"os/exec"
	"strconv"
)

type Plugin struct {
	Init    func(input chan ChromebusRecord, aggregator Aggregator)
	Cleanup func()
}

func Focus(id string) (e error) {
	e = exec.Command(Bin.activateTab, id).Run()
	return
}

func NewTab(url string) (e error) {
	e = exec.Command(Bin.newTab, url).Run()
	return
}

func CloseTab(id string) (e error) {
	e = exec.Command(Bin.closeTab, id).Run()
	return
}

func Navigate(index int, url string) (e error) {
	exec.Command(Bin.navigateTab, strconv.Itoa(index), url).Run()
	return
}
