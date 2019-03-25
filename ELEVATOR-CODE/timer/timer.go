package timer

import (
	"time"
	"../constants"
)

func DoorTimer(start <-chan bool, door_timeout chan<- bool){
	is_active := false
	timestamp:= time.Now()
	tick := time.NewTicker(time.Millisecond*7)
	for{
		select{
		case should_start :=<- start:
			if(should_start){
				timestamp = time.Now()
				is_active = true
			}
		case <- tick.C:
			if is_active && (time.Now().Sub(timestamp) > time.Second*constants.DOOR_OPEN_SEC)  {
				door_timeout <- true
				is_active = false
			}
		}
	}
}
