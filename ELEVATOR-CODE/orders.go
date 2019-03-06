package orders
import "elevio"
import "elevatorstates"


//init orders??:
var orders[N_FLOORS][N_BUTTONS] bool = {false}

func OrderServer(/*channels for write requests, */order_added chan<- bool){
    //handle wishes from other modules to write to orders
    /*for{
        select{
        case a:= <-add_order_request:
            setOrder(order(some type))
            order_added <- true
        case a:= <- clearcurrfloororder_request:
            clearOrderCurrentFloor(elevator)
        }
    }*/
}

//MUST BE PRIVATE METHOD - only used by orders_server
func clearOrderCurrentFloor(e Elevator) {
    //clear order
}

//MUST BE PRIVATE METHOD - only used by orders_server
func setOrder(/*order(some type)*/){
    //orders[x][x] = order

}

func isOrderAbove(floor int) bool {
	//for floor = m to floor > current, floor --
		//if orders.get_is_order(elevator,floor)
			//return true
	//return false
}

func isOrderBelow(floor int) bool {
	//for floor = 0  to floor < current, floor ++
		//if orders.get_is_order(elevator,floor)
			//return true
	//return false
}


func ShouldStop() bool{
	/*
	switch(e.dirn){
    case D_Down:
        return
            e.requests[e.floor][B_HallDown] ||
            e.requests[e.floor][B_Cab]      ||
            !requests_below(e);
    case D_Up:
        return
            e.requests[e.floor][B_HallUp]   ||
            e.requests[e.floor][B_Cab]      ||
            !requests_above(e);
    case D_Stop:
    default:
        return 1;
    }
	
	
	*/
}


func ChooseDirection() /*direction*/{
	/* translate to go:
	//men husk på å fikse sånn at heisen ikke kan kjøres fast

	switch(e.dirn){
    case D_Up:
        return  is_order_above(e) ? D_Up    :
                is_order_below(e) ? D_Down  :
                                    D_Stop  ;
    case D_Down:
    case D_Stop: // there should only be one request in this case. Checking up or down first is arbitrary.
        return  is_order_below(e) ? D_Down  :
                is_order_above(e) ? D_Up    :
                                    D_Stop  ;
    default:
        return D_Stop;
    }*/
}