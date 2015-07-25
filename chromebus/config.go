package chromebus

// use specs to lookup Plugin objects
type PluginSpec string
type PluginCouplingMode int
type Component int

var PluginRegistry = make(map[PluginSpec]*Plugin)

const (
	inMemory PluginCouplingMode = iota
	osLaunch
)

// TODO: probably not needed, as I dont think i need to name the components
const (
//BusReceiver      Component = iota
//PluginController           // for now I'll use a channel, and do everything in memory
)

// A plugin spec is not the plugin object itself. As obviously we do not want
// the engine to have to depend on those. Instead, the PluginSpec is an identifier
// which serves the purpose of allowing the engine to place a channel where the plugin
// can find it, or invoke the plugin's executable

const USE_MOCK_PLUGIN = true

func getDefaultEnabledPlugins() []PluginSpec {
	if USE_MOCK_PLUGIN {
		return []PluginSpec{
			Mock,
		}
	} else {
		return []PluginSpec{
		// define REAL static plugins here
		}
	}
}

func CreateEngine() *Engine {
	return &Engine{
		&ChromebusRecordStdinReceiver{},
		getDefaultEnabledPlugins(),
		&InMemPluginController{
			pluginChannels: map[PluginSpec]chan *ChromebusRecord{
				Mock: make(chan *ChromebusRecord, 100),
			},
		},
	}
}
