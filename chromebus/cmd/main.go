package main

import (
	"chromebus"
	"log"
	"os"
	"os/signal"
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
	e := chromebus.CreateEngine()
	go func() {
		e.Start()
	}()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)
	<-signalChannel
	e.CleanUp()

	//rcv := ChromebusRecordStdinReceiver{}
	//exec.Command("date")
}
