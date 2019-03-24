package states
import (
	"../orders"
	 "../types"
	"../elevio"
	"../lamps"
	"../localip"
	"../constants"
	"fmt"
	"os"
	"time"
)
var all_elevs map[string]types.Elevator
var localelev_ID string

func InitSynchElevators(drv_floors <-chan int){
	all_elevs = make(map[string]types.Elevator)
	e := InitElevator()

	localelev_ID = e.Elevator_ID
	all_elevs[e.Elevator_ID] = e

	select{
	case floor:= <- drv_floors:
			setFloor(floor, e.Elevator_ID)
			fmt.Printf("\nhei\n")
	default:
			if(all_elevs[e.Elevator_ID].Floor == -1){
					//init between floors:
					setDir(elevio.MD_Down, e.Elevator_ID)
					elevio.SetMotorDirection(all_elevs[e.Elevator_ID].Direction)
					setState(types.ES_Moving, e.Elevator_ID)
					fmt.Printf("\nhei2\n")
    	}
	}
}

/*----------------------Single elevator-----------------------------*/
func getPeerID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISElevatorCONNECTED"
	}
	return fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
}

func InitElevator() types.Elevator {
	var elevator types.Elevator

	elevator.Elevator_ID = getPeerID()
	fmt.Printf("\nID: %s\n", elevator.Elevator_ID)
	elevator.State = types.ES_Idle
	elevator.Floor = -1
	elevator.Direction = elevio.MD_Stop
	var ord [constants.N_FLOORS][constants.N_BUTTONS](types.Order)
	elevator.Orders = ord

	return elevator
}
/*------------------------------------------------------------------*/
/*This is the only routine that should write to all_elevs*/
func UpdateElevator(
	/*FSM channels, single elevator:*/
	update_state <-chan types.ElevatorState, update_floor <-chan int, update_direction <-chan elevio.MotorDirection,
	clear_floor <-chan int, floor_reached chan<- bool, order_added chan<- int,
	/*Multiple elevator stuff:*/
	/*elev_tx is only used in this func to send other elevs when requested bc of backup*/
	add_order <-chan types.AssignedOrder, elev_rx <-chan types.Wrapped_Elevator, sendwrap_request <-chan string, elev_tx chan<- types.Wrapped_Elevator) {
		
		for {
    	select{
      case new_state:= <- update_state:
				setState(new_state, localelev_ID)
				//types.PrintStates(all_elevs[localelev_ID])

      case new_floor:= <- update_floor:
				setFloor(new_floor, localelev_ID)
				floor_reached <- true
				//fmt.Printf("\n elev floor set: %d\n", new_floor)

			case new_dir:= <- update_direction:
				setDir(new_dir, localelev_ID)

				//types.PrintStates(all_elevs[localelev_ID])

			case order:= <-add_order:
				//ord := types.Order{-1/*todo: ORDER-ID*/, types.OS_UnacceptedOrder}
				setOrdered(order.Order.Floor, int(order.Order.Button), order.Elevator_ID)

				if(order.Elevator_ID == localelev_ID){
					order_added <- order.Order.Floor /*
					TODO: viktig, dette signalet må flyttes til å gis når det har blitt lagt til noe i local
					sin kø, nå ville den kun fått varsel når den selv har lagt til i sin kø, ikke når andre har lagt til*/
					//fmt.Printf("\nAdded order fl. %d\n", order.Order.Floor)
					lamps.SetAllLamps(all_elevs[localelev_ID]) //todo
				}
			case <- clear_floor:
				setOrderList(orders.ClearAtCurrentFloor(all_elevs[localelev_ID]), localelev_ID)
				lamps.SetAllLamps(all_elevs[localelev_ID])

			case received := <- elev_rx:
				//fmt.Printf("\n\nID: %s\n", received.Elevator_ID)
				if(received.Elevator_ID == localelev_ID){
					break
				} 
				if !keyExists(received.Elevator_ID) {

					var newElev types.Elevator
					newElev.Elevator_ID = received.Elevator_ID
					all_elevs[received.Elevator_ID] = newElev
					setOrderList(received.Orders[received.Elevator_ID] , received.Elevator_ID)
					fmt.Printf("\nAdded new elevator!\n")

				}
				//update states
				setFields(received.State, received.Floor, received.Direction, received.Elevator_ID)
				
				//update orders: Uncomment next two lines when merge is ready
				//order_map := modulename.mergeOrders(getOrderMap(all_elevs), received.Orders)
				//setFromOrderMap(order_map)

			case requestedID := <- sendwrap_request: //when an elev needs its backup
				elev_tx <- wrapElevator(requestedID)
			
			case update:= <- connectionUpdate	//Untested case
				if update.Connected {
					all_elevs[update.Elevator_ID].Connected = true
				} else {
					all_elevs[update.Elevator_ID].Connected = false
					orderReassigner(update.Elevator_ID, false)
				}
			case update:= <- operationUpdate	//Untested case
				if update.isOperational {
					all_elevs[update.Elevator_ID].isOperational = true
				} else {
					all_elevs[update.Elevator_ID].isOperational = false
					orderReassigner(update.Elevator_ID, true)
				}
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
		for key, val := range all_elevs {
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
		temp := all_elevs[elevator_ID] // = the non-wrapped elev
		wrapped.Elevator_ID = temp.Elevator_ID
		wrapped.State = temp.State
		wrapped.Floor = temp.Floor
		wrapped.Direction = temp.Direction
		wrapped.Orders = getOrderMap(all_elevs)
	} else {
		wrapped.Elevator_ID = "NOTFOUND" //should probably have used nil but don't know how
	}
	return wrapped
}

func ReadLocalElevator() types.Elevator {
	return all_elevs[localelev_ID]
}

func keyExists(key string) bool {
	_, exists := all_elevs[key]
	return exists
}

//Workaround functions because go does not allow setting structs in maps directly
func setFields(s types.ElevatorState,f int, d elevio.MotorDirection, ID string){
	setState(s, ID)
	setFloor(f, ID)
	setDir(d, ID)
}

func setState(s types.ElevatorState, ID string) {
	temp, is := all_elevs[ID]
	if(is){
		temp.State = s
		all_elevs[ID] = temp
	} /*else {
		//FATAL //todo, log
	}*/
}

func setFloor(f int, ID string){
	temp, is := all_elevs[ID]
	if(is){
		temp.Floor = f
		all_elevs[ID] = temp
	}
}

func setDir(d elevio.MotorDirection, ID string){
	temp, is := all_elevs[ID]
	if(is){
		temp.Direction = d
		all_elevs[ID] = temp
	}
}


func setOrdered(floor int, button int, ID string){
	temp, is := all_elevs[ID]
	if(is){
		temp.Orders[floor][button].State = types.OS_UnacceptedOrder
		//leave counter unchanged
		all_elevs[ID] = temp
	}
}

func setOrderList(list [constants.N_FLOORS][constants.N_BUTTONS]types.Order, ID string) {
	temp, is := all_elevs[ID]
	if(is){
		temp.Orders = list
		all_elevs[ID] = temp
	}
}

//Untested function
func orderReassigner (faultyElevID string, operationError bool) {
	var e = all_elevs[faultyElevID]

	for i := 0 , i < N_FLOORS, i++ {
		for y := N_BUTTONS - 1, i++ {
			if e.Orders[i][y].State = OS_AcceptedOrder {
				all_elevs[localelev_ID].Orders[i][y].state = OS_AcceptedOrder
			}
		}
	} 

	if operationError {
		var dummyOrder = e.Floor += FloorPlusDir(e)
		all_elevs[faultyElevID].Orders[dummyOrder][0].State = OS_AcceptedOrder //Default: Set to hallUp, should this be changed?
	}
}
