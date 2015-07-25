package chromebus

import (
	"bufio"
	"github.com/ActiveState/tail"
	"log"
	"os"
)

type ChromebusRecordReceiver interface {
	GetRecord(record *ChromebusRecord) (more bool)
}

type ChromebusRecordTailReceiver struct {
	opened bool
	t      *tail.Tail
}

type ChromebusRecordStdinReceiver struct {
	opened  bool
	scanner *bufio.Scanner
}

func (recv ChromebusRecordTailReceiver) open() {
	var err error
	recv.t, err = tail.TailFile(Events, tail.Config{Follow: true})
	if err != nil {
		log.Fatal(err)
	}
	recv.opened = true
}

func (recv ChromebusRecordTailReceiver) close() {
	recv.t.Cleanup()
	recv.close()
}

func (recv ChromebusRecordTailReceiver) GetRecord() (*ChromebusRecord, bool) {
	if !recv.opened {
		recv.open()
	}
	line, more := <-recv.t.Lines

	if more {
		return RecordFromString(line.Text), true
	} else {
		recv.close()
		return nil, false
	}
}

func (recv *ChromebusRecordStdinReceiver) open() {
	recv.scanner = bufio.NewScanner(os.Stdin)
	recv.opened = true
}

func (recv *ChromebusRecordStdinReceiver) GetRecord(record *ChromebusRecord) bool {
	if !recv.opened {
		recv.open()
	}
	if recv.scanner.Scan() {
		*record = *RecordFromString(recv.scanner.Text())
		return true
	} else {
		return false
	}
}
