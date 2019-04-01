package fsm

import (
	"../elevio"
	"../orders"
	"../types"
	"fmt"
)

func FSM(floorReached <-chan bool, clearFloor chan<- int, orderAdded <-chan bool, startDoorTimer chan<- bool, doorTimeout <-chan bool,
	updateState chan<- types.ElevatorState, updateDirection chan<- elevio.MotorDirection, allElevsUpdate <-chan types.Elevator) {
	
	var localElevator types.Elevator

	for {
		select {
		case localElevator = <- allElevsUpdate:
			fmt.Printf("\nupdated\n")
			continue
		default:
			select{
			case localElevator = <- allElevsUpdate:
				fmt.Printf("\nupdated\n")
				continue				
			case <-floorReached:
				fmt.Printf("\n fl reached\n")

				if onFloorArrival(localElevator) {
					clearFloor <- localElevator.Floor
					startDoorTimer <- true
					updateState <- types.ES_DoorOpen
				}

		case <-orderAdded:
			state, dir, startTimer := onListUpdate(localElevator)
			fmt.Printf("\nlist updated\n")

			if startTimer {
				clearFloor <- localElevator.Floor
			}
			updateState <- state
			updateDirection <- dir
			startDoorTimer <- startTimer

		case <-doorTimeout:
			state, dir := onDoorTimeout(localElevator)
			fmt.Printf("\ndoor timeout\n")
			updateState <- state
			updateDirection <- dir
		}
		break
	}
	}
}

func onFloorArrival(localElev types.Elevator) bool {
	elevio.SetFloorIndicator(localElev.Floor)

	switch localElev.State {
	case types.ES_Moving:
		if orders.ShouldStop(localElev) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)

			return true
		}
		return false
	}
	return false
}

func onDoorTimeout(localElev types.Elevator) (types.ElevatorState, elevio.MotorDirection) {
	var dir elevio.MotorDirection
	var state types.ElevatorState

	switch localElev.State {
	case types.ES_DoorOpen:
		elevio.SetDoorOpenLamp(false)
		dir = orders.ChooseDirection(localElev)
		elevio.SetMotorDirection(dir)

		if dir == elevio.MD_Stop {
			state = types.ES_Idle
		} else {
			state = types.ES_Moving
		}
	}
	return state, dir
}

func onListUpdate(localElev types.Elevator) (types.ElevatorState, elevio.MotorDirection, bool) {
	state := localElev.State
	dir := localElev.Direction
	startTimer := false

	switch state {
	case types.ES_DoorOpen:
		if orders.IsOrderCurrentFloor(localElev) {
			startTimer = true
		}

	case types.ES_Idle:
		if orders.IsOrderCurrentFloor(localElev) {
			elevio.SetDoorOpenLamp(true)
			startTimer = true
			state = types.ES_DoorOpen
		} else {
			dir = orders.ChooseDirection(localElev)
			elevio.SetMotorDirection(dir)
			state = types.ES_Moving
		}
	}
	return state, dir, startTimer
}
