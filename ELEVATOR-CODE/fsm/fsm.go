package fsm
import "../orders"
import "../elevatorstates"
import "../elevio"

//doors and lights
//const door_open_duration_s int = 3

func FSM(drv_floors <-chan int, clear_floor chan<- int, order_added <-chan bool, door_timeout <-chan bool, 
	update_state chan<- elevatorstates.ElevatorState, update_floor chan<- int, update_direction chan<- elevio.MotorDirection /*, ...chans*/){
	for{
		select {
		case floor:= <- drv_floors:
			update_floor <- floor //kan flyttes til elevio! Dette er update på last floor.
			if (onFloorArrival(floor)){
				clear_floor <- floor
			}

		case <- order_added:
			state:= onListUpdate()
			update_state <- state

		case is_timeout:=<- door_timeout: //=> doors should be closed
			if(is_timeout){
				state, dir := onDoorTimeout()
				update_state <- state
				update_direction <- dir
			}		
		}
	}
}


func onFloorArrival(floor int) bool  {
	if orders.ShouldStop(floor, elevatorstates.ReadElevator().Direction) {
		elevio.SetMotorDirection(elevio.MD_Stop)
		//SetFloorIndicator
		//and start door timer!
		return true //does stop
	}
	return false //does not stop
}

func onDoorTimeout() (elevatorstates.ElevatorState, elevio.MotorDirection) {
	elev := elevatorstates.ReadElevator()
	var dir elevio.MotorDirection
	var state elevatorstates.ElevatorState

	switch(elev.State){
	case elevatorstates.ES_DoorOpen:
		dir = orders.ChooseDirection(elev.Floor, elev.Direction)
		elevio.SetMotorDirection(dir)
		//SetFloorIndicator to zero
		if(dir == elevio.MD_Stop){
			state = elevatorstates.ES_Idle
			} else {
			state = elevatorstates.ES_Moving
		}
		break;
	default:
		break;
	}

	return state, dir
	
}

func onListUpdate() elevatorstates.ElevatorState {
	elev := elevatorstates.ReadElevator()
	state := elev.State
	switch(state){
	case elevatorstates.ES_Idle:
		dir := orders.ChooseDirection(elev.Floor, elev.Direction)
		elevio.SetMotorDirection(dir)
		state = elevatorstates.ES_Moving 
		break;
	default:
		break;
	}
	return state
}

