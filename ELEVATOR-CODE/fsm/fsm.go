package fsm

import (
	"../elevio"
	"../orders"
	"../states"
	"../types"
)

func FSM(floorReached <-chan bool, clearFloor chan<- int, orderAdded <-chan bool, startDoorTimer chan<- bool, doorTimeout <-chan bool,
	updateState chan<- types.ElevatorState, updateFloor chan<- int, updateDirection chan<- elevio.MotorDirection /*, ...chans*/) {
	for {
		select {
		case <-floorReached:
			if onFloorArrival() {
				clearFloor <- states.ReadLocalElevator().Floor
				startDoorTimer <- true
				updateState <- types.ES_DoorOpen
			}

		case <-orderAdded:
			state, dir, startTimer := onListUpdate()
			if startTimer {
				clearFloor <- states.ReadLocalElevator().Floor
			}
			updateState <- state
			updateDirection <- dir
			startDoorTimer <- startTimer

		case <-doorTimeout:
			state, dir := onDoorTimeout()
			updateState <- state
			updateDirection <- dir
		}
	}
}

func onFloorArrival() bool {
	elevio.SetFloorIndicator(states.ReadLocalElevator().Floor)

	switch states.ReadLocalElevator().State {
	case types.ES_Moving:
		if orders.ShouldStop(states.ReadLocalElevator()) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)

			return true
		}
		return false
	default:
	}
	return false
}

func onDoorTimeout() (types.ElevatorState, elevio.MotorDirection) {
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
	e := states.ReadLocalElevator()
	state := e.State
	dir := e.Direction
	startTimer := false

	switch state {
	case types.ES_DoorOpen:
		if orders.IsOrderCurrentFloor(e) {
			startTimer = true
		}

	case types.ES_Idle:
		if orders.IsOrderCurrentFloor(e) {
			elevio.SetDoorOpenLamp(true)
			startTimer = true
			state = types.ES_DoorOpen
		} else {
			dir = orders.ChooseDirection(e)
			elevio.SetMotorDirection(dir)
			state = types.ES_Moving
		}
	}
	return state, dir, startTimer
}
