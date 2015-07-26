package chromebus

import (
	"os/exec"
)

type Plugin struct {
	Init    func(input chan ChromebusRecord)
	Cleanup func()
}

func Focus(id string) {
	exec.Command(Bin.activateTab, id)
}

func NewTab(url string) {
	exec.Command(Bin.newTab, url)
}

func CloseTab(id string) {
	exec.Command(Bin.closeTab, id)
}

func Navigate(id string, url string) {
	// TODO: this won't work. we need to use the index instead of the id
	// A prereq for this is to set up an aggregator

	//exec.Command(Bin.navigateTab, id, url)
}
