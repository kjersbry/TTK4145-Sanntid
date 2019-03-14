package orders
import "../elevio"
import "../globalconstants"
/*orders should not have direct access to elevatorstates, I think?
should rather be passed only as args floor and direction from fsm

But, it may be prettier that the funcs that need floor and
dir takes elevator as argument from fsm.
*/

type AssignedOrder struct {			//Used by orderassigner.AssignOrder --> Untested
	Order elevator_io.ButtonEvent
	Elevator_ID int
}

//trengs det noe init orders??:
var orders = [globalconstants.N_FLOORS][globalconstants.N_BUTTONS]bool{{false}}

func UpdateOrders(add_order <-chan elevio.ButtonEvent, clear_floor <-chan int, order_added chan<- bool){
    //handle wishes from other modules to write to orders
    for {
        select{
        case order:= <-add_order:
            setOrder(order)
            order_added <- true
        case floor:= <- clear_floor:
            clearOrdersAtFloor(floor)
       // case a:= <- clearspecificorder //eventuelt
        }
    }
}

//MUST BE PRIVATE METHOD - only used by orders_server
func clearOrdersAtFloor(floor int) {
    for i := 0; i < globalconstants.N_BUTTONS; i++ {
        orders[floor][i] = false
    }
}

//MUST BE PRIVATE METHOD - only used by orders_server
func setOrder(order elevio.ButtonEvent){
    orders[order.Floor][order.Button] = true
}

func IsOrder(floor int, button elevio.ButtonType) bool {
    return orders[floor][button]
}

func isOrderAbove(current_floor int) bool {
	for floor := globalconstants.N_FLOORS; floor > current_floor; floor-- {
        for button := 0; button < globalconstants.N_BUTTONS; button++ {
            if IsOrder(floor, elevio.ButtonType(button)) {
                return true
            }
        }
    }
    return false
}

func isOrderBelow(current_floor int) bool {
	for floor := 0; floor < current_floor; floor++ {
        for button := 0; button < globalconstants.N_BUTTONS; button++ {
            if IsOrder(floor, elevio.ButtonType(button)) {
                return true
            }
        }
    }
    return false
}


func ShouldStop(current_floor int, direction elevio.MotorDirection) bool{
	switch(direction){
    case elevio.MD_Down:
        return (IsOrder(current_floor, elevio.BT_HallDown) ||
            IsOrder(current_floor, elevio.BT_Cab)      ||
            !isOrderBelow(current_floor))
    case elevio.MD_Up:
        return (IsOrder(current_floor, elevio.BT_HallUp)   ||
            IsOrder(current_floor, elevio.BT_Cab)      ||
            !isOrderAbove(current_floor))
    case elevio.MD_Stop:
    default:
        break
    }
    return true
}


func ChooseDirection(current_floor int, direction elevio.MotorDirection) elevio.MotorDirection {
	// husk på å teste at heisen ikke kan kjøres fast hvis noen vil være kjipe

	switch(direction){
        //must use if else, go does not support " ? : "

    case elevio.MD_Up:
        if isOrderAbove(current_floor) {
            return elevio.MD_Up
        } else if isOrderBelow(current_floor){
            return elevio.MD_Down
        }
        return elevio.MD_Stop

    case elevio.MD_Down: //Any reason why this one is left empty? If intentional, add in a comment #codeQuality
    case elevio.MD_Stop: // there should only be one request in this case. Checking up or down first is arbitrary.

        if  isOrderBelow(current_floor) {
            return elevio.MD_Down
        } else if isOrderAbove(current_floor) {
            return elevio.MD_Up
        }
    }
    return elevio.MD_Stop
}
