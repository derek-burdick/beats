package crawler

import (
	"fmt"

	"github.com/elastic/beats/filebeat/config"
	"github.com/elastic/beats/filebeat/input"
	"github.com/elastic/beats/libbeat/logp"
)

/*
 The hierarchy for the crawler objects is explained as following

 Crawler: Filebeat has one crawler. The crawler is the single point of control
 	and stores the state. The state is written through the registrar
 Prospector: For every FileConfig the crawler starts a prospector
 Harvester: For every file found inside the FileConfig, the Prospector starts a Harvester
 		The harvester send their events to the spooler
 		The spooler sends the event to the publisher
 		The publisher writes the state down with the registrar
*/

type Crawler struct {
	// Registrar object to persist the state
	Registrar *Registrar
	running   bool
}

func (crawler *Crawler) Start(prospectorConfigs []config.ProspectorConfig, eventChan chan *input.FileEvent) error {

	crawler.running = true

	if len(prospectorConfigs) == 0 {
		return fmt.Errorf("No prospectors defined. You must have at least one prospector defined in the config file.")
	}

	var prospectors []*Prospector

	logp.Info("Loading Prospectors: %v", len(prospectorConfigs))

	// Prospect the globs/paths given on the command line and launch harvesters
	for _, prospectorConfig := range prospectorConfigs {

		logp.Debug("prospector", "File Configs: %v", prospectorConfig.Paths)

		prospector, err := NewProspector(prospectorConfig, crawler.Registrar, eventChan)
		prospectors = append(prospectors, prospector)

		if err != nil {
			return fmt.Errorf("Error in initing prospector: %s", err)
		}
	}
	logp.Info("Loading Prospectors completed")

	logp.Info("Running Prospectors")
	for _, prospector := range prospectors {
		go prospector.Run()
	}
	logp.Info("All prospectors are running")

	logp.Info("All prospectors initialised with %d states to persist", len(crawler.Registrar.State))

	return nil
}

func (crawler *Crawler) Stop() {
	// TODO: Properly stop prospectors and harvesters
}
