package fsm
import (
	"../orders"
	 "../states"
	"../elevio"
)
//doors and lights

var elevator states.Elevator

func InitElevator(drv_floors <-chan int){
	//somehow initialize ID uniquely
	//init rest of elev stuff!!
    select{
    case floor:= <- drv_floors:
        elevator.Floor = floor
    }
    if(elevator.Floor == -1){
        //init between floors:
        elevator.Direction = elevio.MD_Down
        elevio.SetMotorDirection(elevator.Direction)
        elevator.State = states.ES_Moving
    }
}

func UpdateElevator(update_ID <-chan int, update_state <-chan states.ElevatorState,
	update_floor <-chan int, update_direction <-chan elevio.MotorDirection, 
	/*Orders: */
	add_order <-chan elevio.ButtonEvent, clear_floor <-chan int, order_added chan<- bool){
        for {
            select{
            case new_id:= <- update_ID:
                elevator.Elevator_ID = new_id
            case new_state:= <- update_state:
                elevator.State = new_state
            case new_floor:= <- update_floor:
                elevator.Floor = new_floor
            case new_dir:= <- update_direction:
				elevator.Direction = new_dir
				
			case order:= <-add_order:
				orders.SetOrder(elevator, order)
				order_added <- true
			//case floor:= <- clear_floor:
			case <- clear_floor:
				orders.ClearAtCurrentFloor(elevator)
			}
        }
}

func ReadElevator() states.Elevator {
	return elevator
}


/*------------------------------ funcs for running FSM -----------------------------------------------------------------------------------*/
//FSM trenger kun å oppdatere elevator via channels dersom det er fare for at den oppdateres fra et annet sted samtidig.
//kun ved init etter krasj tror jeg?
//men mye kan oppdateres direkte til updateelev istedet for via fsm funk
func FSM(drv_floors <-chan int/*, clear_floor chan<- int*/, order_added <-chan bool, start_door_timer chan<- bool, door_timeout <-chan bool, 
	update_state chan<- states.ElevatorState, update_floor chan<- int, update_direction chan<- elevio.MotorDirection /*, ...chans*/){
	for{
		select {
		case floor:= <- drv_floors:
			update_floor <- floor //kan flyttes til elevio! Dette er update på last floor.
			if (onFloorArrival(floor)){ //stops on order
				//clear_floor <- floor //old
				orders.ClearAtCurrentFloor(elevator)
				start_door_timer <- true 
			}

		case <- order_added:
			state:= onListUpdate()
			update_state <- state

		case <- door_timeout: //=> doors should be closed
				state, dir := onDoorTimeout()
				update_state <- state
				update_direction <- dir	
		}
	}
}


func onFloorArrival(floor int) bool  {
	//SetFloorIndicator(floor)
	if orders.ShouldStop(elevator) {
		elevio.SetMotorDirection(elevio.MD_Stop)
		return true //does stop
	}
	return false //does not stop
}

func onDoorTimeout() (states.ElevatorState, elevio.MotorDirection) {
	var dir elevio.MotorDirection
	var state states.ElevatorState

	switch(elevator.State){
	case states.ES_DoorOpen:
		dir = orders.ChooseDirection(elevator)
		elevio.SetMotorDirection(dir)
		//SetFloorIndicator to zero
		if(dir == elevio.MD_Stop){
			state = states.ES_Idle
			} else {
			state = states.ES_Moving
		}
		break;
	default:
		break;
	}

	return state, dir
	
}

func onListUpdate() states.ElevatorState {
	state := elevator.State
	switch(state){
	case states.ES_Idle:
		dir := orders.ChooseDirection(elevator)
		elevio.SetMotorDirection(dir)
		state = states.ES_Moving 
		break;
	default:
		break;
	}
	return state
}

