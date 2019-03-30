package timer

import (
	"../constants"
	"time"
)

func DoorTimer(start <-chan bool, doorTimeout chan<- bool){
	isActive := false
	timestamp:= time.Now()
	tick := time.NewTicker(time.Millisecond*7)
	for{
		select{
		case shouldStart :=<- start:
			if(shouldStart){
				timestamp = time.Now()
				isActive = true
			}
		case <- tick.C:
			if isActive && (time.Now().Sub(timestamp) > time.Second*constants.DOOR_OPEN_SEC)  {
				doorTimeout <- true
				isActive = false
			}
		}
	}
}
