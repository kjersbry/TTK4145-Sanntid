package timer

import (
	"time"
	"../constants"
)

func DoorTimer(start <-chan bool, door_timeout chan<- bool){
	is_active := false
	timestamp:= time.Now()
	//var tick time.Ticker
	for{
		select{
		case should_start :=<- start:
			if(should_start){
				timestamp = time.Now()
				is_active = true
			}
		default:
			if (time.Now().Sub(timestamp) > time.Second*constants.DOOR_OPEN_SEC) && is_active {
				door_timeout <- true
				is_active = false
			}
			/*tick:= time.NewTicker(time.Millisecond*7)
		case <- tick.C:
			if (time.Now().Sub(timestamp) > time.Second*constants.DOOR_OPEN_SEC) && is_active {
				door_timeout <- true
				is_active = false
			}*/
		}
	}
}
