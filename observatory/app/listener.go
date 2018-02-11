package main

import(
	"github.com/nlittlepoole/constellation/observatory/rover"
	"fmt"
)

func listen(stop <-chan string, stopped chan<- string){
    defer close(stopped)
    found := make(chan rover.Probe, 1000)
    go rover.Scan("mon0", found, 5 * 60)
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
