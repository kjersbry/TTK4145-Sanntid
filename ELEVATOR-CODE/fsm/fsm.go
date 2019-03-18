package fsm
import (
	"../orders"
	 "../states"
	"../elevio"
	"../globalconstants"
	"fmt"
)
//lights
/*
PROBLEMS:
Stopper noen ganger når den ikke skal, andre ganger ikke når den skal
Grunn: dir = stop

Noen ganger venter den i 1 etg selv om den har bestilling i 0


Når alt funker:
Init between floors kjøres hver gang
init: den åpner døra selv om den bare skulle kjøre ned til og idle etasjen
 sjekk hva de gjør. har de spes tilfelle for dører når det ikke er orders above or below

Fiks det setalllampsgreiene
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
	update_floor <-chan int, update_direction <-chan elevio.MotorDirection, floor_reached chan<- bool,
	/*Orders: */
	add_order <-chan elevio.ButtonEvent, clear_floor <-chan int, order_added chan<- int){

		for {
    	select{
      case new_id:= <- update_ID:
        elevator.Elevator_ID = new_id
      case new_state:= <- update_state:
				elevator.State = new_state
				
				states.PrintStates(elevator)

      case new_floor:= <- update_floor:
				elevator.Floor = new_floor
				floor_reached <- true
				fmt.Printf("\n elev floor set: %d\n", new_floor)

			case new_dir:= <- update_direction:
				elevator.Direction = new_dir

				states.PrintStates(elevator)

			case order:= <-add_order:
				elevator.Orders = orders.SetOrder(elevator, order)
				order_added <- order.Floor
				//states.PrintOrders(elevator)
			case <- clear_floor:
				elevator.Orders = orders.ClearAtCurrentFloor(elevator)
			}
    }
}

func ReadElevator() states.Elevator {
	return elevator
}


/*------------------------------ funcs for running FSM -----------------------------------------------------------------------------------*/
//FSM trenger kun å oppdatere elevator via channels dersom det er fare for at den oppdateres fra et annet sted samtidig
//kun ved init etter krasj tror jeg?
func FSM(floor_reached <-chan bool, clear_floor chan<- int, order_added <-chan int, start_door_timer chan<- bool, door_timeout <-chan bool,
	update_state chan<- states.ElevatorState, update_floor chan<- int, update_direction chan<- elevio.MotorDirection /*, ...chans*/){
	for{
		select {
		case <-floor_reached:
			//update_floor <- floor //flyttet til elevio! Dette er update på last floor.
			//fmt.Printf("\nnew floor: %d\n", floor)
			if (onFloorArrival()){ //stops on order
				clear_floor <- elevator.Floor
				start_door_timer <- true
				update_state <- states.ES_DoorOpen
			}

		case floor := <- order_added:
			state, dir, start_timer := onListUpdate(floor)
			update_state <- state
			update_direction <- dir
			start_door_timer <- start_timer
			if(start_timer){
				clear_floor <- floor
			}

		case <- door_timeout: //=> doors should be closed
				state, dir := onDoorTimeout()
				update_state <- state
				update_direction <- dir
		}
	}
}


func onFloorArrival() bool  {
	fmt.Printf("\nonFloorArrival\n")
	elevio.SetFloorIndicator(elevator.Floor)
	if orders.ShouldStop(elevator) {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true)
		return true //does stop
	}
	return false //does not stop
}

func onDoorTimeout() (states.ElevatorState, elevio.MotorDirection) {
	fmt.Printf("\nonDoorTimeout\n")

	var dir elevio.MotorDirection
	var state states.ElevatorState

	switch(elevator.State){
	case states.ES_DoorOpen:
		elevio.SetDoorOpenLamp(false)
		dir = orders.ChooseDirection(elevator)
		elevio.SetMotorDirection(dir)
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

func onListUpdate(floor int) (states.ElevatorState, elevio.MotorDirection, bool) {
	fmt.Printf("\nonListUpdate\n")

	state := elevator.State
	dir := elevator.Direction
	start_timer := false


	switch(state){
	case states.ES_DoorOpen:
		if(elevator.Floor == floor){
			start_timer = true
		}
		break

	case states.ES_Idle:
		if(elevator.Floor == floor){
			elevio.SetDoorOpenLamp(true)
			start_timer = true
			state = states.ES_DoorOpen
		} else{
			dir = orders.ChooseDirection(elevator)
			elevio.SetMotorDirection(dir)
			state = states.ES_Moving
		}
		break;
	default:
		break;
	}
	//setalllights //todo


	return state, dir, start_timer
}
