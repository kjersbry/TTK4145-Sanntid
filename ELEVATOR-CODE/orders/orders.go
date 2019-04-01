package orders
import (
    "../elevio"
    "../types"
    "../constants"
)

func ClearAtCurrentFloor(e types.Elevator) [constants.N_FLOORS][constants.N_BUTTONS]types.Order {
    for i := 0; i < constants.N_BUTTONS; i++ {
        e.Orders[e.Floor][i].State = types.OS_NoOrder
        e.Orders[e.Floor][i].Counter++
    }
    return e.Orders
}

func IsOrder(e types.Elevator, floor int, button elevio.ButtonType) bool {
    if(button == elevio.BT_Cab){
        return (e.Orders[floor][button].State != types.OS_NoOrder)
    } else {
        return (e.Orders[floor][button].State == types.OS_AcceptedOrder)
    }
}

func IsOrderCurrentFloor(e types.Elevator) bool {
    for i := 0; i < constants.N_BUTTONS; i++ {
        if(IsOrder(e, e.Floor, elevio.ButtonType(i))){
            return true
        }
    }
    return false
}

func isOrderAbove(e types.Elevator) bool {
	for floor := e.Floor + 1; floor < constants.N_FLOORS; floor++ {
        for button := 0; button < constants.N_BUTTONS; button++ {
            if IsOrder(e, floor, elevio.ButtonType(button)) {
                return true
            }
        }
    }
    return false
}

func isOrderBelow(e types.Elevator) bool {
	for floor := 0; floor < e.Floor; floor++ {
        for button := 0; button < constants.N_BUTTONS; button++ {
            if IsOrder(e, floor, elevio.ButtonType(button)) {
                return true
            }
        }
    }
    return false
}


func ShouldStop(e types.Elevator) bool {
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
    }
    return true
}


func ChooseDirection(e types.Elevator) elevio.MotorDirection {
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
        fallthrough
    case elevio.MD_Stop:
        if  isOrderBelow(e) {
            return elevio.MD_Down
        } else if isOrderAbove(e) {
            return elevio.MD_Up
        }
    }
    return elevio.MD_Stop
}
