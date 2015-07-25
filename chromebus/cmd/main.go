package main

import (
	"chromebus"
	"log"
	//"os/exec"
)

func main() {
	chromebus.EnvSetup()
	log.SetFlags(log.Flags() | log.Llongfile)

	//recv := &chromebus.ChromebusRecordStdinReceiver{}
	//str, _ := recv.GetRecord()
	//log.Printf("%s", str)
	//str, _ = recv.GetRecord()
	//log.Printf("%s", str)
	chromebus.CreateEngine().Start()
	//rcv := ChromebusRecordStdinReceiver{}
	//exec.Command("date")
}
