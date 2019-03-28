package fsm

import (
	"../elevio"
	"../orders"
	"../states"
	"../types"
	"../constants"
	"fmt"
)

/*
Init between floors kjøres hver gang
init: den åpner døra selv om den bare skulle kjøre ned til og idle etasjen
 sjekk hva de gjør. har de spes tilfelle for dører når det ikke er orders above or below

initElev burde også bruke channels mot elevator
*/

/*------------------------------ funcs for running FSM -----------------------------------------------------------------------------------*/
//FSM trenger kun å oppdatere elevator via channels dersom det er fare for at den oppdateres fra et annet sted samtidig
//kun ved init etter krasj tror jeg?
func FSM(floor_reached <-chan bool, clear_floor chan<- int, order_added <-chan bool, start_door_timer chan<- bool, door_timeout <-chan bool,
	update_state chan<- types.ElevatorState, update_floor chan<- int, update_direction chan<- elevio.MotorDirection /*, ...chans*/) {
	for {
		select {
		case <-floor_reached:
			if onFloorArrival() { //stops on order
				clear_floor <- states.ReadLocalElevator().Floor
				start_door_timer <- true
				update_state <- types.ES_DoorOpen
			}

		case <-order_added:
			state, dir, start_timer := onListUpdate()
			if start_timer {
				clear_floor <- states.ReadLocalElevator().Floor
				fmt.Printf("\nstart_timer was true\n")
			}
			update_state <- state
			update_direction <- dir
			start_door_timer <- start_timer

		case <-door_timeout: //=> doors should be closed
			state, dir := onDoorTimeout()
			update_state <- state
			update_direction <- dir
		}
	}
}

func onFloorArrival() bool {
fmt.Printf("\nonFloorArrival\n")
	elevio.SetFloorIndicator(states.ReadLocalElevator().Floor)

	switch states.ReadLocalElevator().State {
	case types.ES_Moving:
		if orders.ShouldStop(states.ReadLocalElevator()) {
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

	switch states.ReadLocalElevator().State {
	case types.ES_DoorOpen:
		elevio.SetDoorOpenLamp(false)
		dir = orders.ChooseDirection(states.ReadLocalElevator())
		elevio.SetMotorDirection(dir)

		if dir == elevio.MD_Stop {
			state = types.ES_Idle
		} else {
			state = types.ES_Moving
		}
	}
	return state, dir
}

func onListUpdate() (types.ElevatorState, elevio.MotorDirection, bool) {
	fmt.Printf("\nonListUpdate\n")

	e := states.ReadLocalElevator()
	state := e.State
	fmt.Printf("Current state: %s\n", types.StateToString(state))
	dir := e.Direction
	start_timer := false

	switch state {
	case types.ES_DoorOpen:
		//if(states.ReadLocalElevator().Floor == floor){ previous version
		if orders.IsOrderCurrentFloor(e) {
			start_timer = true
		}

	case types.ES_Idle:
		//if(states.ReadLocalElevator().Floor == floor){
		if orders.IsOrderCurrentFloor(e) {
			elevio.SetDoorOpenLamp(true)
			start_timer = true
			state = types.ES_DoorOpen
		} else {
			fmt.Printf("\nExecuting ChooseDirection\n")
			dir = orders.ChooseDirection(e)
			fmt.Printf("\nFinished ChooseDirection\n")
			elevio.SetMotorDirection(dir)
			fmt.Printf("\nDir was set to: %s\n", types.DirToString(dir))
			state = types.ES_Moving
		}
	}
	/*this should not happen, but need to avoid getting 
		stuck if there's a bug in listupdatesignalling somewhere:*/
	if dir == elevio.MD_Stop && state == types.ES_Moving {
		if e.Floor == constants.N_FLOORS {
			dir = elevio.MD_Down
		} else {
			dir = elevio.MD_Up
		}
		elevio.SetMotorDirection(dir)
	}
	return state, dir, start_timer
}
