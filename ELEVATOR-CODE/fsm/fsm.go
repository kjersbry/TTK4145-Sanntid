package fsm
import (
	"../orders"
	 "../states"
	"../elevio"
	"../globalconstants"
	"fmt"
)
//doors and lights
/*
PREVIOUS PROBLEM:
bestillinger ble ikke lagt til i køen og heller ikke cleared
Grunn: orders tar inn kopi av elevator, så det blir ikke endret direkte på den
løsning: set og clear orders må returnere elevator!
-
NESTE PROBLEM:
heisen kjører ikke når den har bestillinger
Finner ut neste gang
-
*/



var elevator states.Elevator

func InitElevator(drv_floors <-chan int){
	//somehow initialize ID uniquely
	//init rest of elev stuff
	elevator.Elevator_ID = -1
	elevator.State = states.ES_Idle
	elevator.Floor = -1
	elevator.Direction = elevio.MD_Stop
	var ord [globalconstants.N_FLOORS][globalconstants.N_BUTTONS](states.Order)
	elevator.Orders = ord

  select{
  case floor:= <- drv_floors:
      elevator.Floor = floor
			fmt.Printf("\nhei\n")
	default:
			if(elevator.Floor == -1){
		      //init between floors:
		      elevator.Direction = elevio.MD_Down
		      elevio.SetMotorDirection(elevator.Direction)
		      elevator.State = states.ES_Moving
					fmt.Printf("\nhei2\n")
    	}
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
				
				states.PrintStates(elevator)
				
      case new_floor:= <- update_floor:
        elevator.Floor = new_floor
					fmt.Printf("\nelev floor is: %d\n", elevator.Floor)
      case new_dir:= <- update_direction:
				elevator.Direction = new_dir

				states.PrintStates(elevator)

			case order:= <-add_order:
				elevator = orders.SetOrder(elevator, order)
				order_added <- true
				states.PrintOrders(elevator)
			//case floor:= <- clear_floor:
			case <- clear_floor:
				elevator = orders.ClearAtCurrentFloor(elevator)
			}
    }
}

func ReadElevator() states.Elevator {
	return elevator
}


/*------------------------------ funcs for running FSM -----------------------------------------------------------------------------------*/
//FSM trenger kun å oppdatere elevator via channels dersom det er fare for at den oppdateres fra et annet sted samtidig
//kun ved init etter krasj tror jeg?
func FSM(drv_floors <-chan int, clear_floor chan<- int, order_added <-chan bool, start_door_timer chan<- bool, door_timeout <-chan bool,
	update_state chan<- states.ElevatorState, update_floor chan<- int, update_direction chan<- elevio.MotorDirection /*, ...chans*/){
	for{
		select {
		case floor:= <- drv_floors:
			update_floor <- floor //kan flyttes til elevio! Dette er update på last floor.
			fmt.Printf("\nnew floor: %d\n", floor)

			if (onFloorArrival(floor)){ //stops on order
				clear_floor <- floor
				start_door_timer <- true
				update_state <- states.ES_DoorOpen
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
		//set door indicator
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
		//SetDoorIndicator to zero
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
