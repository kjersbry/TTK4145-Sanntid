package fsm
import "../orders"
import "../elevatorstates"
import "../elevio"

//doors and lights
//const door_open_duration_s int = 3

func FSM(drv_floors <-chan int, clear_floor chan<- int, order_added <-chan bool/*, ...chans*/){
	for{
		select {
		case floor:= <- drv_floors:
			if (onFloorArrival(floor)){
				clear_floor <- floor
			}

		case <- order_added:
			onListUpdate()

		/*case <- door_timeout: //= closing doors
			state = onDoorTimeout()
			channel_for_updating_state <- state //forslag, bruke stateserver
		*/
		
		}
	}
}


func onFloorArrival(floor int) bool {
	//elevator.floor = floor //use channel, stateserver
	if orders.ShouldStop(floor, elevatorstates.ReadElevator().Direction) {
		elevio.SetMotorDirection(elevio.MD_Stop)
		//and start door timer!
		return true //does stop
	}
	return false //does not stop
}

func onDoorTimeout() {
	elev := elevatorstates.ReadElevator()

	switch(elev.State){
	case elevatorstates.ES_DoorOpen:
		dir := orders.ChooseDirection(elev.Floor, elev.Direction)
		elevio.SetMotorDirection(dir)
		//Sets state:
		//NB: do through stateserver
		/* elev.Direction = dir
		if(dir == MD_Stop){
			elev.State = elevatorstates.ES_Idle
			} else {
			elev.State = EB_Moving
		}*/
		break;
	default:
		break;
	}

	
}

func onListUpdate() {
	elev := elevatorstates.ReadElevator()
	switch(elev.State){
	case elevatorstates.ES_Idle:
		dir := orders.ChooseDirection(elev.Floor, elev.Direction)
		elevio.SetMotorDirection(dir)
		/*set through write channel: stateserver
		elev.State = elevatorstates.ES_Moving */
		break;
	default:
		break;
	}
}

