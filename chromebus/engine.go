package chromebus

//import "log"

type Engine struct {
	rcv        ChromebusRecordReceiver
	plugins    []PluginSpec
	pluginCtrl PluginController
}

var primaryConfig map[Component]interface{}

func (e *Engine) Start() {
	e.run()
}

func (e *Engine) run() {
	for _, pluginSpec := range e.plugins {
		e.pluginCtrl.Init(pluginSpec)
		defer e.pluginCtrl.Cleanup(pluginSpec)
	}
	for record := new(ChromebusRecord); e.rcv.GetRecord(record); {
		//log.Printf("record: %s", record)
		e.broadcast(record)
	}
}

func (e *Engine) broadcast(record *ChromebusRecord) {
	for _, pluginSpec := range e.plugins {
		e.pluginCtrl.Send(pluginSpec, record)
	}
}
