package main

import(
	"fmt"
	"time"
)

func main(){
     stop := make(chan string)
     stopped := make(chan string)
     go listen(stop, stopped)
 
     time.Sleep(120 * time.Second)
     close(stop)
     <-stopped

     fmt.Println(GetAllUniques(time.Minute))
     fmt.Println(GetCurrentUniques(time.Minute))
     fmt.Println(GetReturningUniques(time.Now().Add(-300 * time.Second), time.Now()))
     fmt.Println(GetStrengthHistogram(time.Now().Add(-300 * time.Second), time.Now()))
}