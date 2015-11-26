package main

import (
	"flag"
	"github.com/washtubs/chromebus"
	"log"
	"os"
	"os/signal"
	//"os/exec"
)

// URGGHHHHH FUCK.
// OK, even with VisitAll, it still only takes the last instance of the same flag.
// to have multiple values for the same flag I need to follow *this* example:
// http://play.golang.org/p/Ig5sm7jA14
//
// HOWEVER: I dont really care. right now, there is no use case for multiple flags.
// TODO

var plugins []chromebus.PluginSpec = make([]chromebus.PluginSpec, 0, 32)
var pluginsDefined bool = false
var port int = 0

func defFlags(flag *flag.Flag) {
	switch flag.Name {
	case "p": // plugin
		//var err error
		//var p chromebus.PluginSpec
		p, err := chromebus.StringToPluginSpec(flag.Value.String())
		if err == nil {
			plugins = append(plugins, p)
			pluginsDefined = true
		} else {
			log.Fatal(err)
		}
	}
}

func init() {
	flag.String("p", "default", "specify a plugin")
	flag.IntVar(&port, "P", 0, "specify a port number to listen on")
}

func main() {
	chromebus.EnvSetup()
	log.SetFlags(log.Flags() | log.Llongfile)

	//recv := &chromebus.ChromebusRecordStdinReceiver{}
	//str, _ := recv.GetRecord()
	//log.Printf("%s", str)
	//str, _ = recv.GetRecord()
	//log.Printf("%s", str)
	flag.Parse()
	flag.VisitAll(defFlags)
	var e *chromebus.Engine
	if pluginsDefined {
		e = chromebus.CreateEngine(plugins, port)
	} else {
		e = chromebus.CreateEngine(nil, port)
	}
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
