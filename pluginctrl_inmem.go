// TODO: remove shitty doc
// Call the plugin's functions directly. If this were over stdin, we would need to init with a channel we create. And a go routine would be used to pass to that channel from stdin

package chromebus

import "net/http"

type InMemPluginController struct {
	pluginChannels map[PluginSpec]chan ChromebusRecord
}

// TODO: dummy code this is the only way I know how to assert an interface is implemented
var test PluginController = &InMemPluginController{}

func initHttpServer(spec PluginSpec) {
	http.HandleFunc("/"+string(spec)+"/", PluginRegistry[spec].Handle)
}

func (i *InMemPluginController) Init(spec PluginSpec) {
	go func() {
		PluginRegistry[spec].Init(i.pluginChannels[spec], AggregatorModel)
	}()
	initHttpServer(spec)
}

func (i *InMemPluginController) Cleanup(spec PluginSpec) {
	PluginRegistry[spec].Cleanup()
}

func (i *InMemPluginController) Send(spec PluginSpec, record *ChromebusRecord) {
	i.pluginChannels[spec] <- *record
}

func init() {

}
