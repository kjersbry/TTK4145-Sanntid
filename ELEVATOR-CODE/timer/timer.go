package timer

import "time"
import "fmt"

//ish forslag:
func DoorTimer(start <-chan bool, door_timeout chan<- bool){
	flag := false
	timestamp:= time.Now()
	for{
		select{
		case <- start:
			timestamp = time.Now()
			flag = true
			fmt.Printf("\n\nDOOR OPEN\n")
		default:
			if (time.Now().Sub(timestamp) > time.Second*3) && flag {
				door_timeout <- true
				flag = false
				fmt.Printf("\n\ntimer DOOR CLOSE\n")

			}
		}
	}
}
