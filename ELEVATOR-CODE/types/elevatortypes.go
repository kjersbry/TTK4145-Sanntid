package types

import (
	"fmt"
	"strconv"

	"../constants"
	"../elevio"
	//"time"
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
	//litt midlertidig oppsett, hvordan skal denne være?
	Counter int64
	State   OrderState
	//Timestamp time.Time
}

type AssignedOrder struct { //Used by orderassigner.AssignOrder
	Elevator_ID string
	Order       elevio.ButtonEvent
}

type Elevator struct {
	Elevator_ID    string
	State          ElevatorState
	Floor          int
	Direction      elevio.MotorDirection //does only change to stop when IDLE, not when stopping for order
	Is_Operational bool
	Connected      bool //todo --> is
	Orders         [constants.N_FLOORS][constants.N_BUTTONS]Order
}

type Wrapped_Elevator struct {
	Elevator_ID string
	State       ElevatorState
	Floor       int
	Direction   elevio.MotorDirection //does only change to stop when IDLE, not when stopping for order
	Orders      map[string][constants.N_FLOORS][constants.N_BUTTONS]Order
	//Is_Operational bool
}

type Operation_Event struct {
	Elevator_ID    string
	Is_Operational bool
}

type Connection_Event struct {
	Elevator_ID string
	Connected   bool
}

/*The below functions are used for debugging*/

func OrderToString(o Order) string {
	switch o.State {
	case OS_NoOrder:
		return "No"
	case OS_AcceptedOrder:
	case OS_UnacceptedOrder:
	}
	return "Yes"
}

func PrintOrders(e Elevator) {
	fmt.Printf("\n\n-----Queue------\n")
	for i := 0; i < constants.N_FLOORS; i++ {
		for j := 0; j < constants.N_BUTTONS; j++ {
			fmt.Printf("%s ", OrderToString(e.Orders[i][j]))
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

func elevToString(e Elevator) string {
	result := "State: "
	switch e.State {
	case ES_Idle:
		result += "Idle\n"
	case ES_DoorOpen:
		result += "Door open\n"
	case ES_Moving:
		result += "Moving\n"
	}
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
