package chromebus

import (
	"log"
	"strconv"
	"strings"
)

type ChromeTab struct {
	url     string
	tabType string
	index   int
	focused string
}

type ChromebusRecord struct {
	action string
	id     string
	oldTab *ChromeTab
	newTab *ChromeTab
}

func TabFromString(s string) *ChromeTab {
	if s == "nil" {
		return nil
	}
	fields := strings.Split(s, tabDelimiter)
	url := fields[0]
	tabType := fields[1]
	index, err := strconv.Atoi(fields[2])
	if err != nil {
		log.Fatal("Failed to convert tab index %s to integer", fields[2])
	}
	focused := fields[3]
	return &ChromeTab{
		url:     url,
		tabType: tabType,
		index:   index,
		focused: focused,
	}
}

func RecordFromString(s string) *ChromebusRecord {
	fields := strings.Split(s, delimiter)
	action := fields[0]
	id := fields[1]
	rawOldTab := fields[2]
	rawNewTab := fields[3]
	return &ChromebusRecord{
		action: action,
		id:     id,
		oldTab: TabFromString(rawOldTab),
		newTab: TabFromString(rawNewTab),
	}
}

type Action string

const (
	New          Action = "new"
	Closed       Action = "closed"
	UrlChanged   Action = "urlchanged"
	FocusChanged Action = "focuschanged"
)

const delimiter string = "||"
const tabDelimiter string = "|+|"
