package main

import (
	"github.com/nlittlepoole/observatory/rover"
)

func listen(stop <-chan string, stopped chan<- string) {
	defer close(stopped)
	found := make(chan rover.Probe, 1000)
	go rover.Scan(
		found,
		ACTIVE_SETTINGS.Driver,
		ACTIVE_SETTINGS.Session(),
		ACTIVE_SETTINGS.Location,
		ACTIVE_SETTINGS.Threshold,
		ACTIVE_SETTINGS.SampleRate,
	)
	for {
		select {
		default:
			p := <-found
			log.Debug(p)
			if err := logEvent(p); err != nil {
				log.Warn(err)
			}
		case <-stop:
			// stop
			return
		}
	}
}
