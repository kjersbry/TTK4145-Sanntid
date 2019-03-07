package orders
import "../elevio"
import "../globalconstants"
/*orders should not have direct access to elevatorstates, I think?
should rather be passed only as args floor and direction from fsm

But, it may be prettier that the funcs that need floor and
dir takes elevator as argument from fsm. 
*/

//trengs det noe init orders??:
var orders = [globalconstants.N_FLOORS][globalconstants.N_BUTTONS]bool{{false}}

func OrdersServer(add_order <-chan elevio.ButtonEvent, clear_floor <-chan int, order_added chan<- bool){
    //handle wishes from other modules to write to orders
    for{
        select{
        case order:= <-add_order:
            setOrder(order)
            order_added <- true
        case floor:= <- clear_floor:
            clearOrdersAtFloor(floor)
       // case a:= <- clearspecificorder
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

func isOrder(floor int, button elevio.ButtonType) bool {
    return orders[floor][button]
}

func isOrderAbove(current_floor int) bool {
	for floor := globalconstants.N_FLOORS; floor > current_floor; floor-- {
        for button := 0; button < globalconstants.N_BUTTONS; button++ {
            if isOrder(floor, elevio.ButtonType(button)) {
                return true
            }
        }
    }
    return false
}

func isOrderBelow(current_floor int) bool {
	for floor := 0; floor < current_floor; floor++ {
        for button := 0; button < globalconstants.N_BUTTONS; button++ {
            if isOrder(floor, elevio.ButtonType(button)) {
                return true
            }
        }
    }
    return false
}


func ShouldStop(current_floor int, direction elevio.MotorDirection) bool{
	switch(direction){
    case elevio.MD_Down:
        return (isOrder(current_floor, elevio.BT_HallDown) ||
            isOrder(current_floor, elevio.BT_Cab)      ||
            !isOrderBelow(current_floor))
    case elevio.MD_Up:
        return (isOrder(current_floor, elevio.BT_HallUp)   ||
            isOrder(current_floor, elevio.BT_Cab)      ||
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
        /*
    case D_Up:
        return  isOrderAbove(current_floor) ? elevio.MD_Up    :
                isOrderBelow(current_floor) ? elevio.MD_Down  :
                                    elevio.MD_Stop  ;
    case D_Down:
    case D_Stop: // there should only be one request in this case. Checking up or down first is arbitrary.
        return  isOrderBelow(current_floor) ? elevio.MD_Down  :
                isOrderAbove(current_floor) ? elevio.MD_Up    :
                                    elevio.MD_Stop  ;
    default:
        return elevio.MD_Stop;*/
    }
    return elevio.MD_Down //temporary for build
}