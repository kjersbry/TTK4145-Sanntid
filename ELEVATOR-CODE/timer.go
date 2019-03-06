package timer


/* ish forslag:
func DoorTimer(start chan<- bool, door_timeout <-chan bool){
	for{
		select{
		case <- start:
			start timer
			while sometimerof3sec() not ended do nothing
			stopped <- true
		}
	}
}
*/