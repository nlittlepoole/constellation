package main

import(
	"github.com/nlittlepoole/constellation/observatory/rover"
	"fmt"
)

func listen(stop <-chan string, stopped chan<- string){
    defer close(stopped)
    found := make(chan rover.Probe, 1000)
    go rover.Scan(
       found,
       ACTIVE_SETTINGS.Driver,
       ACTIVE_SETTINGS.Window(),
       ACTIVE_SETTINGS.Location,
       ACTIVE_SETTINGS.Threshold,
       ACTIVE_SETTINGS.SampleRate,
    )
    for {
    	select{
		default:
			p := <-found
			fmt.Println(p)
			if err := logEvent(p); err != nil {
	   		   fmt.Println(err)
			}
		case <-stop:
		     // stop
		     return
	}
    }
}
