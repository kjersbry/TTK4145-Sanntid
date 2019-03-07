package fsm
import "orders"
import "elevatorstates"

//doors and lights

//const door_open_duration_s int = 3

func FSM(drv_floors <-chan int, clear_floor chan<- int, order_added <-chan bool/*, ...chans*/){
	for{
		select {
		case floor:= <- drv_floors:
			if (fsm.onFloorArrival(floor)){
				clear_floor <- floor
			}

		case <- order_added:
			fsm.onListUpdate()

		/*case <- door_timeout: //= closing doors
			state = FSM.onDoorTimeout()
			channel_for_updating_state <- state //forslag
		*/
		
		}
	}
}


func onFloorArrival(floor) bool {
	elevator.floor = floor
	if orders.Should_stop() {
		elevio.SetMotorDirection(MD_Stop)
		return true //does stop
	}
	return false //does not stop
}

func onDoorTimeout() {
	elev := elevatorstates.Elevator()

	switch(elev.state){
	case ES_DoorOpen:
		dir := orders.ChooseDirection(elev.floor, elev.direction)
		elevio.SetMotorDirection(dir)
		//Sets state:
		//NB: do through channels, can't do directly  like this.
		/* elev.direction = dir
		if(dir == MD_Stop){
			elev.state = ES_Idle
			} else {
			elev.state = EB_Moving
		}*/
		break;
	default:
		break;
	}

	
}

func onListUpdate() {
	elev := elevatorstates.Elevator()
	switch(elev.state){
	case ES_Idle:
		dir := orders.ChooseDirection(elev.floor, elev.direction)
		elevio.SetMotorDirection(dir)
		/*set through write channel:
		elev.state = ES_Moving */
		break;
	default:
		break;
	}
}

