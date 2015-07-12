package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type ChromebusRecordReceiver interface {
	GetRecord() (record *ChromebusRecord)
}

type ChromebusRecordStdinReceiver struct {
}

type ChromeTab struct {
	url     string
	tabType string
	focused string
}

type ChromebusRecord struct {
	action string
	id     string
	oldTab *ChromeTab
	newTab *ChromeTab
}

func TabFromString(s string) *ChromeTab {
	fields := strings.Split(s, tabDelimiter)
	url := fields[0]
	tabType := fields[1]
	focused := fields[2]
	return &ChromeTab{
		url:     url,
		tabType: tabType,
		focused: focused,
	}
}

func RecordFromString(s string) *ChromebusRecord {
	fields := strings.Split(s, delimiter)
	log.Printf("fields: %s", fields)
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
	return nil
}

const delimiter string = "||"
const tabDelimiter string = "|+|"

func (*ChromebusRecordStdinReceiver) GetRecord() (record *ChromebusRecord) {
	stdinReader := bufio.NewReader(os.Stdin)
	recordBuf, _, err := stdinReader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	rawRecord := string(recordBuf)
	log.Printf("rawRecord: %s", rawRecord)
	return RecordFromString(rawRecord)
	return nil
}

func main() {
	rcv := ChromebusRecordStdinReceiver{}
	record := rcv.GetRecord()
	log.Printf("record: %s", record)
}
