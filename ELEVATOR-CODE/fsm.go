package fsm
import "orders"
//doors and lights

const door_open_duration_s int = 3

func FSM(drv_floors <-chan int, order_added <-chan bool /*, ...chans*/){
	for{
		select {
		case floor:= <- drv_floors:
			fsm.onFloorArrival(floor, current_direction)

		case <- order_added:
			fsm.onListUpdate()

		/*case <- door_timeout: //= closing doors
			FSM.onDoorTimeout()
		*/
		
		}
	}
}


func onFloorArrival(/*floor, direction*/){
	//update last_registered_floor
	//check if orders.Should_stop()
		//set motor dir and state to stop
}

func onDoorTimeout() {
	/*if(state = door_open)
		direction = orders.choose_direction()
		//set direction

	if(elevator.dirn == D_Stop){
	elevator.behaviour = EB_Idle;
	} else {
	elevator.behaviour = EB_Moving;
	*/
}

func onListUpdate() {
	//switch(state), case == IDLE:
		//direction = orders.choose_direction()
		//set direction
		//state = moving 
	//default:
		//do nothing
}

