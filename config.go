package chromebus

import (
	"errors"
)

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
const defaultPort = 8081

func getDefaultEnabledPlugins() []PluginSpec {
	if USE_MOCK_PLUGIN {
		return []PluginSpec{
			//Mock,
			ActivityTracker,
			SiteBlocker,
		}
	} else {
		return []PluginSpec{
		// define REAL static plugins here
		}
	}
}

func StringToPluginSpec(s string) (PluginSpec, error) {
	switch s {
	case "ActivityTracker":
		return ActivityTracker, nil
	case "SiteBlocker":
		return SiteBlocker, nil
	}
	return "", errors.New(s + " is not a valid plugin")
}

func CreateEngine(specifiedPlugins []PluginSpec, port int) *Engine {
	var enabledPlugins []PluginSpec
	if specifiedPlugins != nil {
		enabledPlugins = specifiedPlugins
	} else {
		enabledPlugins = getDefaultEnabledPlugins()
	}
	if port == 0 {
		port = defaultPort
	}
	return &Engine{
		&ChromebusRecordStdinReceiver{},
		enabledPlugins,
		&InMemPluginController{
			pluginChannels: map[PluginSpec]chan ChromebusRecord{
				// just allocate all the channels unnecessarily. It's not like it's expensive.
				//Mock:            make(chan ChromebusRecord, 0),
				ActivityTracker: make(chan ChromebusRecord, 0),
				SiteBlocker:     make(chan ChromebusRecord, 0),
			},
		},
		port,
	}
}
