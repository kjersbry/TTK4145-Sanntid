package fsm
import (
	"../orders"
	 "../types"
	"../elevio"
	"../lamps"
	"../localip"
	"../constants"
	"fmt"
	"os"
)
/*
Init between floors kjøres hver gang
init: den åpner døra selv om den bare skulle kjøre ned til og idle etasjen
 sjekk hva de gjør. har de spes tilfelle for dører når det ikke er orders above or below

initElev burde også bruke channels mot elevator
*/


var elevator types.Elevator

func getPeerID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	return fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
}

func InitElevator(drv_floors <-chan int){
	elevator.Elevator_ID = getPeerID()
	fmt.Printf("\nID: %s\n", elevator.Elevator_ID)
	elevator.State = types.ES_Idle
	elevator.Floor = -1
	elevator.Direction = elevio.MD_Stop
	var ord [constants.N_FLOORS][constants.N_BUTTONS](types.Order)
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
		      elevator.State = types.ES_Moving
					fmt.Printf("\nhei2\n")
    	}
  }
}

func UpdateElevator(update_ID <-chan string, update_state <-chan types.ElevatorState,
	update_floor <-chan int, update_direction <-chan elevio.MotorDirection, floor_reached chan<- bool,
	/*Orders: */
	add_order <-chan elevio.ButtonEvent, clear_floor <-chan int, order_added chan<- int){

		for {
    	select{
      case new_id:= <- update_ID:
        elevator.Elevator_ID = new_id
      case new_state:= <- update_state:
				elevator.State = new_state
				
				types.PrintStates(elevator)

      case new_floor:= <- update_floor:
				elevator.Floor = new_floor
				floor_reached <- true
				fmt.Printf("\n elev floor set: %d\n", new_floor)

			case new_dir:= <- update_direction:
				elevator.Direction = new_dir

				types.PrintStates(elevator)

			case order:= <-add_order:
				elevator.Orders = orders.SetOrder(elevator, order)
				order_added <- order.Floor
				fmt.Printf("\nAdded order fl. %d\n", order.Floor)
				lamps.SetAllLamps(elevator)
			case <- clear_floor:
				elevator.Orders = orders.ClearAtCurrentFloor(elevator)
				lamps.SetAllLamps(elevator)
			}
    }
}

func ReadElevator() types.Elevator {
	return elevator
}


/*------------------------------ funcs for running FSM -----------------------------------------------------------------------------------*/
//FSM trenger kun å oppdatere elevator via channels dersom det er fare for at den oppdateres fra et annet sted samtidig
//kun ved init etter krasj tror jeg?
func FSM(floor_reached <-chan bool, clear_floor chan<- int, order_added <-chan int, start_door_timer chan<- bool, door_timeout <-chan bool,
	update_state chan<- types.ElevatorState, update_floor chan<- int, update_direction chan<- elevio.MotorDirection /*, ...chans*/){
	for{
		select {
		case <- floor_reached:
			if (onFloorArrival()){ //stops on order
				clear_floor <- elevator.Floor
				start_door_timer <- true
				update_state <- types.ES_DoorOpen
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

	switch(elevator.State){
	case types.ES_Moving:
		if orders.ShouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			
			return true //does stop
		}
		return false //does not stop
	default:
	}
	return false
}

func onDoorTimeout() (types.ElevatorState, elevio.MotorDirection) {
	fmt.Printf("\nonDoorTimeout start\n")

	var dir elevio.MotorDirection
	var state types.ElevatorState

	switch(elevator.State){
	case types.ES_DoorOpen:
		elevio.SetDoorOpenLamp(false)
		dir = orders.ChooseDirection(elevator)
		elevio.SetMotorDirection(dir)
		fmt.Printf("dir: %s", types.DirToString(dir))

		if(dir == elevio.MD_Stop){
			fmt.Printf("\n chosen dir was stop")
			state = types.ES_Idle
			} else {
			state = types.ES_Moving
		}
	}
	fmt.Printf("\nonDoorTimeout end\n")
	return state, dir
}

func onListUpdate(floor int) (types.ElevatorState, elevio.MotorDirection, bool) {
	fmt.Printf("\nonListUpdate\n")

	state := elevator.State
	dir := elevator.Direction
	start_timer := false


	switch(state){
	case types.ES_DoorOpen:
		if(elevator.Floor == floor){
			start_timer = true
		}

	case types.ES_Idle:
		if(elevator.Floor == floor){
			elevio.SetDoorOpenLamp(true)
			start_timer = true
			state = types.ES_DoorOpen
		} else {
			dir = orders.ChooseDirection(elevator)
			elevio.SetMotorDirection(dir)
			state = types.ES_Moving
		}
	default:
	}
	return state, dir, start_timer
}