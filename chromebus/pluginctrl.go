package chromebus

// All calls should be asynchronous
type PluginController interface {
	Init(spec PluginSpec)
	Send(spec PluginSpec, record *ChromebusRecord)
	Cleanup(spec PluginSpec)
}
