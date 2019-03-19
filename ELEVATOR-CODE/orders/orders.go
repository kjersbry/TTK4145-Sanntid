package orders
import (
    "../elevio"
    "../globalconstants"
    "../states"
    "fmt"
)

func ClearAtCurrentFloor(e states.Elevator) [globalconstants.N_FLOORS][globalconstants.N_BUTTONS]states.Order {
    if(e.Floor < 0 || e.Floor > 3){
      fmt.Printf("\nclear: out of range %d \n", e.Floor)
      return e.Orders
    } //todo: ta vekk! litt for quickfix

    for i := 0; i < globalconstants.N_BUTTONS; i++ {
        e.Orders[e.Floor][i].State = states.OS_NoOrder
        //elevio.SetButtonLamp(elevio.ButtonType(i), e.Floor, false) //todo move when several elevs
    }

    return e.Orders
}

func SetOrder(e states.Elevator, order elevio.ButtonEvent) [globalconstants.N_FLOORS][globalconstants.N_BUTTONS]states.Order {
  if(order.Floor < 0 || order.Button > 3){
    fmt.Printf("\nSet: out of range %d \n", e.Floor)
    return e.Orders
  } //todo:ta vekk! litt for quickfix

    e.Orders[order.Floor][order.Button].State = states.OS_UnacceptedOrder
    //elevio.SetButtonLamp(order.Button, order.Floor, true) //todo move when several elevs

    return e.Orders
}

func IsOrder(e states.Elevator, floor int, button elevio.ButtonType) bool {
    //todo: vurder denne. Når skal den si at det er bestilling (= når skal den stoppe), skal den stoppe på unaccepted også?
    if(floor < 0 || floor > 3){
      fmt.Printf("\nIs: out of range %d \n", floor)
      return false
    } //todo: ta vekk! litt for quickfix


    return (e.Orders[floor][button].State != states.OS_NoOrder)
}

func isOrderAbove(e states.Elevator) bool {
	for floor := e.Floor + 1; floor < globalconstants.N_FLOORS; floor++ {
        for button := 0; button < globalconstants.N_BUTTONS; button++ {
            if IsOrder(e, floor, elevio.ButtonType(button)) {
                return true
            }
        }
    }
    return false
}

func isOrderBelow(e states.Elevator) bool {
	for floor := 0; floor < e.Floor; floor++ {
        for button := 0; button < globalconstants.N_BUTTONS; button++ {
            if IsOrder(e, floor, elevio.ButtonType(button)) {
                return true
            }
        }
    }
    return false
}


func ShouldStop(e states.Elevator) bool {
	switch(e.Direction){
    case elevio.MD_Down:
        return (IsOrder(e, e.Floor, elevio.BT_HallDown) ||
            IsOrder(e, e.Floor, elevio.BT_Cab)      ||
            !isOrderBelow(e))
    case elevio.MD_Up:
        return (IsOrder(e, e.Floor, elevio.BT_HallUp)   ||
            IsOrder(e, e.Floor, elevio.BT_Cab)      ||
            !isOrderAbove(e))
    case elevio.MD_Stop:
    default:
    }
    return true
}


func ChooseDirection(e states.Elevator) elevio.MotorDirection {
	// husk på å teste at heisen ikke kan kjøres fast hvis noen vil være kjipe

	switch(e.Direction){
    case elevio.MD_Up:
        if isOrderAbove(e) {
            return elevio.MD_Up
        } else if isOrderBelow(e){
            return elevio.MD_Down
        } else {
            return elevio.MD_Stop
        }
    case elevio.MD_Down:
        if  isOrderBelow(e) {
            return elevio.MD_Down
        } else if isOrderAbove(e) {
            return elevio.MD_Up
        }
    case elevio.MD_Stop: // there should only be one request in this case. Checking up or down first is arbitrary.
        if  isOrderBelow(e) {
            return elevio.MD_Down
        } else if isOrderAbove(e) {
            return elevio.MD_Up
        }
    }
    return elevio.MD_Stop
}
