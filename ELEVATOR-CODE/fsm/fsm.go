package fsm

import (
	"../elevio"
	"../orders"
	"../states"
	"../types"
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
func FSM(floorReached <-chan bool, clearFloor chan<- int, orderAdded <-chan bool, startDoorTimer chan<- bool, doorTimeout <-chan bool,
	updateState chan<- types.ElevatorState, updateFloor chan<- int, updateDirection chan<- elevio.MotorDirection /*, ...chans*/) {
	for {
		select {
		case <-floorReached:
			if onFloorArrival() { //stops on order
				clearFloor <- states.ReadLocalElevator().Floor
				startDoor_timer <- true
				updateState <- types.ES_DoorOpen
			}

		case <-orderAdded:
			state, dir, start_timer := onListUpdate()
			if startTimer {
				clearFloor <- states.ReadLocalElevator().Floor
			}
			updateState <- state
			updateDirection <- dir
			startDoorTimer <- startTimer

		case <-doorTimeout: //=> doors should be closed
			state, dir := onDoorTimeout()
			updateState <- state
			updateDirection <- dir
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

	state := states.ReadLocalElevator().State
	dir := states.ReadLocalElevator().Direction
	start_timer := false

	switch state {
	case types.ES_DoorOpen:
		//if(states.ReadLocalElevator().Floor == floor){ previous version
		if orders.IsOrderCurrentFloor(states.ReadLocalElevator()) {
			startTimer = true
		}

	case types.ES_Idle:
		//if(states.ReadLocalElevator().Floor == floor){
		if orders.IsOrderCurrentFloor(states.ReadLocalElevator()) {
			elevio.SetDoorOpenLamp(true)
			startTimer = true
			state = types.ES_DoorOpen
		} else {
			dir = orders.ChooseDirection(states.ReadLocalElevator())
			elevio.SetMotorDirection(dir)
			state = types.ES_Moving
		}
	default:
	}
	return state, dir, startTimer
}
