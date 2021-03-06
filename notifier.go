package chromebus

import (
	"fmt"
	"os/exec"
)

type Notifier interface {
	Init() error
	SendMessage(string)
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
		fmt.Errorf("notify-send not found.", err)
		fmt.Printf("Message was %s:", message)
	}
}

var n2 Notifier = new(NotifySendNotifier)
