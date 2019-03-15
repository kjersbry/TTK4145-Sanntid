package timer

import "time"

//ish forslag:
func DoorTimer(start <-chan bool, door_timeout chan<- bool){
	flag := false
	timestamp:= time.Now()
	for{
		select{
		case <- start:
			timestamp = time.Now()
			flag = true
		default:
			if (time.Now().Sub(timestamp) > time.Second*3) && flag {
				door_timeout <- true
				flag = false
			}
		}
	}
}
