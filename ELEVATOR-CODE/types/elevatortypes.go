package types

import (
	"fmt"
	"strconv"
	"../constants"
	"../elevio"
)

type ElevatorState int

const (
	ES_Idle     ElevatorState = 0
	ES_DoorOpen               = 1
	ES_Moving                 = 2
)

type OrderState int

const (
	OS_NoOrder         OrderState = 0
	OS_UnacceptedOrder            = 1
	OS_AcceptedOrder              = 2
)

type Order struct {
	Counter int64
	State   OrderState
}

type AssignedOrder struct { 
	ElevatorID string
	Order       elevio.ButtonEvent
}

type Elevator struct {
	ElevatorID    string
	State          ElevatorState
	Floor          int
	Direction      elevio.MotorDirection
	IsOperational bool
	IsConnected      bool
	Orders         [constants.N_FLOORS][constants.N_BUTTONS]Order
}

type WrappedElevator struct {
	ElevatorID string
	State       ElevatorState
	Floor       int
	Direction   elevio.MotorDirection
	Orders      map[string][constants.N_FLOORS][constants.N_BUTTONS]Order
}

type OperationEvent struct {
	ElevatorID    string
	IsOperational bool
}

type ConnectionEvent struct {
	ElevatorID string
	IsConnected   bool
}

/*The below functions are used for debugging*/

func OrderToString(o Order) string {
	switch o.State {
	case OS_NoOrder:
		return "No"
	case OS_AcceptedOrder:
		return "Acc"
	case OS_UnacceptedOrder:
		return "Un"
	}
	return "-"
}

func PrintOrders(e Elevator) {
	fmt.Printf("\n\n-----Queue------\n")
	for j := 0; j < constants.N_BUTTONS; j++ {
		for i := 0; i < constants.N_FLOORS; i++ {
			if(j == 0 && i == 3) || (j == 1 && i == 0){
				fmt.Printf("-  ")

			} else {
				fmt.Printf("%s ", OrderToString(e.Orders[i][j]))
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n----------------\n\n")
}

func DirToString(dir elevio.MotorDirection) string {
	var result string
	switch dir {
	case elevio.MD_Up:
		result = "up"
	case elevio.MD_Down:
		result = "down"
	case elevio.MD_Stop:
		result = "stop"
	}
	return result
}

func StateToString(st ElevatorState) string {
	var result string
	switch st {
	case ES_Idle:
		result = "Idle"
	case ES_DoorOpen:
		result = "Door open"
	case ES_Moving:
		result = "Moving"
	}
	return result
}

func elevToString(e Elevator) string {
	result := "State: " + StateToString(e.State) + "\n"
	result += "Floor: "
	result += strconv.Itoa(e.Floor)
	result += "\nDirection: " + DirToString(e.Direction) + "\n"

	return result
}

func PrintStates(e Elevator) {
	fmt.Printf("\n-----States------\n")
	fmt.Printf("%s", elevToString(e))
	fmt.Printf("\n----------------\n\n")
}
