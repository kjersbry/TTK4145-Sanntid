package states
import (
	"../orders"
	 "../types"
	"../elevio"
	"../lamps"
	"../constants"
	"fmt"
	"time"
)
var all_elevators map[string]types.Elevator
var localelev_ID string

func InitElevators(local_ID string, drv_floors <-chan int){
	localelev_ID = local_ID
	all_elevators = make(map[string]types.Elevator)

	//Initialize local elevator
	var empty_elevator types.Elevator
	empty_elevator.Elevator_ID = localelev_ID
	all_elevators[localelev_ID] = empty_elevator

	var ord [constants.N_FLOORS][constants.N_BUTTONS](types.Order)
	setFields(types.ES_Idle, -1, elevio.MD_Stop, localelev_ID)
	setOrderList(ord, localelev_ID)
	
	//wait to allow floor signal to arrive if we start on a floor
	time.Sleep(time.Millisecond*50)
	select{
	case floor:= <- drv_floors:
			setFloor(floor, localelev_ID)
			fmt.Printf("\nhei\n")
	default:
			if(all_elevators[localelev_ID].Floor == -1){
					//init between floors:
					setDir(elevio.MD_Down, localelev_ID)
					elevio.SetMotorDirection(all_elevators[localelev_ID].Direction)
					setState(types.ES_Moving, localelev_ID)
					fmt.Printf("\nhei2\n")
    	}
	}
}

/*------------------------------------------------------------------*/
/*This is the only routine that should write to all_elevators*/
func UpdateElevator(
	/*FSM channels, single elevator:*/
	update_state <-chan types.ElevatorState, update_floor <-chan int, update_direction <-chan elevio.MotorDirection,
	clear_floor <-chan int, floor_reached chan<- bool, order_added chan<- bool,
	/*Multiple elevator stuff:*/
	/*elev_tx is only used in this func to send other elevs when requested bc of backup*/
	add_order <-chan types.AssignedOrder, elev_rx <-chan types.Wrapped_Elevator, sendwrap_request <-chan string, elev_tx chan<- types.Wrapped_Elevator) {
		
		for {
    	select{
      case new_state:= <- update_state:
				setState(new_state, localelev_ID)
				//types.PrintStates(all_elevators[localelev_ID])

      case new_floor:= <- update_floor:
				setFloor(new_floor, localelev_ID)
				floor_reached <- true
				//fmt.Printf("\n elev floor set: %d\n", new_floor)

			case new_dir:= <- update_direction:
				setDir(new_dir, localelev_ID)

				//types.PrintStates(all_elevators[localelev_ID])

			case order:= <-add_order:
				//ord := types.Order{-1/*todo: ORDER-ID*/, types.OS_UnacceptedOrder}
				setOrdered(order.Order.Floor, int(order.Order.Button), order.Elevator_ID)

				if(order.Elevator_ID == localelev_ID){
					order_added <- true 
					/*TODO: viktig, dette signalet må i tillegg gis når det har blitt lagt til noe i local
					sin kø fra merge, nå ville den kun fått varsel når den selv har lagt til i sin kø, ikke når andre har lagt til*/
					//fmt.Printf("\nAdded order fl. %d\n", order.Order.Floor)
					lamps.SetAllLamps(all_elevators[localelev_ID]) //todo
				}
			case <- clear_floor:
				setOrderList(orders.ClearAtCurrentFloor(all_elevators[localelev_ID]), localelev_ID)
				lamps.SetAllLamps(all_elevators[localelev_ID])

			case received := <- elev_rx:
				//fmt.Printf("\n\nID: %s\n", received.Elevator_ID)
				if(received.Elevator_ID == localelev_ID || received.Elevator_ID == "NOTFOUND"){
					break
				} 
				if !keyExists(received.Elevator_ID) {

					var newElev types.Elevator
					newElev.Elevator_ID = received.Elevator_ID
					all_elevators[received.Elevator_ID] = newElev
					setOrderList(received.Orders[received.Elevator_ID] , received.Elevator_ID)
					fmt.Printf("\nAdded new elevator!\n")

				}
				//update states
				setFields(received.State, received.Floor, received.Direction, received.Elevator_ID)
				
				//update orders: Uncomment next two lines when merge is ready
				//order_map, is_new_local_order := modulename.mergeOrders(getOrderMap(all_elevators), received.Orders)
				//setFromOrderMap(order_map)
				//if(is_new_local_order){ order_added <- true } 

			case /*requestedID := */ <- sendwrap_request: //when an elev needs its backup
				//elev_tx <- wrapElevator(requestedID)
			}
    }
}

func TransmitElev(elev_tx chan<- types.Wrapped_Elevator){
	for {
			elev_tx <- wrapElevator(localelev_ID)
			time.Sleep(time.Millisecond*constants.TRANSMIT_MS)
	}
}

func TestPrintAllElevators(){
	for{
		fmt.Printf("\n\n")
		for key, val := range all_elevators {
			/*TEST:*/
			fmt.Printf("\nElev: %s", key)
			types.PrintStates(val)
		}
		time.Sleep(time.Second*5)
	}
}

func getOrderMap(elevator_map map[string]types.Elevator) map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order {
	order_map := make(map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order)
	for key, val := range elevator_map {
		order_map[key] = val.Orders
	}
	return order_map
}

func setFromOrderMap(order_map map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order) {
	for key, val := range order_map {
		setOrderList(val, key)
	}
}

/*By using ID we can easily use this function for wrapping any elevator*/
func wrapElevator(elevator_ID string) types.Wrapped_Elevator {
	var wrapped types.Wrapped_Elevator
	if(keyExists(elevator_ID)){
		temp := all_elevators[elevator_ID] // = the non-wrapped elev
		wrapped.Elevator_ID = temp.Elevator_ID
		wrapped.State = temp.State
		wrapped.Floor = temp.Floor
		wrapped.Direction = temp.Direction
		wrapped.Orders = getOrderMap(all_elevators)
	} else {
		wrapped.Elevator_ID = "NOTFOUND" //should probably have used nil but don't know how
	}
	return wrapped
}

func ReadLocalElevator() types.Elevator {
	return all_elevators[localelev_ID]
}

func keyExists(key string) bool {
	_, exists := all_elevators[key]
	return exists
}

//Workaround functions because go does not allow setting structs in maps directly
func setFields(s types.ElevatorState, f int, d elevio.MotorDirection, ID string){
	setState(s, ID)
	setFloor(f, ID)
	setDir(d, ID)
}

func setState(s types.ElevatorState, ID string) {
	temp, is := all_elevators[ID]
	if(is){
		temp.State = s
		all_elevators[ID] = temp
	} /*else {
		//FATAL //todo, log
	}*/
}

func setFloor(f int, ID string){
	temp, is := all_elevators[ID]
	if(is){
		temp.Floor = f
		all_elevators[ID] = temp
	}
}

func setDir(d elevio.MotorDirection, ID string){
	temp, is := all_elevators[ID]
	if(is){
		temp.Direction = d
		all_elevators[ID] = temp
	}
}


func setOrdered(floor int, button int, ID string){
	temp, is := all_elevators[ID]
	if(is){
		temp.Orders[floor][button].State = types.OS_UnacceptedOrder
		//leave counter unchanged
		all_elevators[ID] = temp
	}
}

func setOrderList(list [constants.N_FLOORS][constants.N_BUTTONS]types.Order, ID string) {
	temp, is := all_elevators[ID]
	if(is){
		temp.Orders = list
		all_elevators[ID] = temp
	}
}