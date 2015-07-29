package chromebus

import (
	"log"
	"os/exec"
)

type Notifier interface {
	Init() (e error)
	SendMessage(message string)
}

type PushBulletNotifier struct{}

func (PushBulletNotifier) Init() (e error) {
	return nil
}

func (PushBulletNotifier) SendMessage(message string) {

}

var n Notifier = new(PushBulletNotifier)

type NotifySendNotifier struct{}

func (NotifySendNotifier) Init() (e error) {
	if _, err := exec.LookPath("notify-send"); err != nil {
		// TODO: customize error?
		return err
	} else {
		return nil
	}
}

func (NotifySendNotifier) SendMessage(message string) {
	err := exec.Command("notify-send", message).Run()
	if err != nil {
		log.Fatal(err)
	}
}

var n2 Notifier = new(NotifySendNotifier)
