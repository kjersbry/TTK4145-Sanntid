package timer

import (
	"time"
	"../constants"
)

func DoorTimer(start <-chan bool, door_timeout chan<- bool){
	is_active := false
	timestamp:= time.Now()
	for{
		select{
		case should_start :=<- start:
			if(should_start){
				timestamp = time.Now()
				is_active = true
				//fmt.Printf("\n\nDOOR OPEN\n")
			}
		default:
			if (time.Now().Sub(timestamp) > time.Second*constants.DOOR_OPEN_SEC) && is_active {
				door_timeout <- true
				is_active = false
				//fmt.Printf("\n\ntimer DOOR CLOSE\n")

			}
		}
	}
}
