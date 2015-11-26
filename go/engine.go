package chromebus

import (
	"log"
	"net/http"
	"strconv"
)

type Engine struct {
	rcv        ChromebusRecordReceiver
	plugins    []PluginSpec
	pluginCtrl PluginController
	port       int
}

var primaryConfig map[Component]interface{}

func (e *Engine) Start() {
	e.run()
}

func welcome(w http.ResponseWriter, r *http.Request) {
	log.Printf("HIHIHIHI")
}

func (e *Engine) run() {
	for _, pluginSpec := range e.plugins {
		e.pluginCtrl.Init(pluginSpec)
	}
	http.HandleFunc("/", welcome)
	log.Printf("Listening on %s", ":"+strconv.Itoa(e.port))
	go http.ListenAndServe(":"+strconv.Itoa(e.port), nil)
	for record := new(ChromebusRecord); e.rcv.GetRecord(record); {
		e.broadcast(record)
	}
}

func (e *Engine) CleanUp() {
	for _, pluginSpec := range e.plugins {
		e.pluginCtrl.Cleanup(pluginSpec)
	}
}

func (e *Engine) broadcast(record *ChromebusRecord) {
	for _, pluginSpec := range e.plugins {
		e.pluginCtrl.Send(pluginSpec, record)
	}
}
