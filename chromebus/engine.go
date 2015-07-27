package chromebus

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
	}
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
