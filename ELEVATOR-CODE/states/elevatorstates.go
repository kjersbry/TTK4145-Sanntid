package states

import (
	"fmt"
	"time"

	"../constants"
	"../elevio"
	"../lamps"
	"../orders"
	"../types"
)

var all_elevators map[string]types.Elevator
var localelev_ID string

func InitElevators(local_ID string, drv_floors <-chan int) {
	localelev_ID = local_ID
	all_elevators = make(map[string]types.Elevator)

	//Initialize local elevator
	var empty_elevator types.Elevator
	empty_elevator.Elevator_ID = localelev_ID
	all_elevators[localelev_ID] = empty_elevator

	var ord [constants.N_FLOORS][constants.N_BUTTONS](types.Order)
	setFields(types.ES_Idle, -1, elevio.MD_Stop, localelev_ID)
	setOrderList(ord, localelev_ID)
	setOperational(true, localelev_ID)
	setConnected(true, localelev_ID)

	//wait to allow floor signal to arrive if we start on a floor
	time.Sleep(time.Millisecond * 50)
	select {
	case floor := <-drv_floors:
		setFloor(floor, localelev_ID)
		fmt.Printf("\nhei\n")
	default:
		if all_elevators[localelev_ID].Floor == -1 {
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
	add_order <-chan types.AssignedOrder, elev_rx <-chan types.Wrapped_Elevator, connectionUpdate <-chan types.Connection_Event/*, operationUpdate <-chan types.Operation_Event*/) {

	for {
		select {
		case new_state := <-update_state:
			setState(new_state, localelev_ID)
			//types.PrintStates(all_elevators[localelev_ID])

		case new_floor := <-update_floor:
			setFloor(new_floor, localelev_ID)
			floor_reached <- true
			//fmt.Printf("\n elev floor set: %d\n", new_floor)

		case new_dir := <-update_direction:
			setDir(new_dir, localelev_ID)

			//types.PrintStates(all_elevators[localelev_ID])

		case order := <-add_order:
			//ord := types.Order{-1/*todo: ORDER-ID*/, types.OS_UnacceptedOrder}
			setOrdered(order.Order.Floor, int(order.Order.Button), order.Elevator_ID, false)

			if order.Elevator_ID == localelev_ID {
				order_added <- true
				/*TODO: viktig, dette signalet må i tillegg gis når det har blitt lagt til noe i local
				sin kø fra merge, nå ville den kun fått varsel når den selv har lagt til i sin kø, ikke når andre har lagt til*/
				//fmt.Printf("\nAdded order fl. %d\n", order.Order.Floor)
				lamps.SetAllLamps(all_elevators[localelev_ID])
			}
		case <-clear_floor:
			setOrderList(orders.ClearAtCurrentFloor(all_elevators[localelev_ID]), localelev_ID)
			lamps.SetAllLamps(all_elevators[localelev_ID])

		case received := <-elev_rx:
			//fmt.Printf("\n\nrec: ID: %s\n", received.Elevator_ID)

			if received.Elevator_ID == localelev_ID {
				break
			}
			if received.Elevator_ID == "NOTFOUND" || !isValidID(received.Elevator_ID) {
				//log/print: something weird has happened
				break
			}

			if !keyExists(received.Elevator_ID) {
				//fmt.Printf("\n\nID: %s\n", received.Elevator_ID)
				test := unwrapElevator(received)
				types.PrintStates(test)
				all_elevators[test.Elevator_ID] = test //unwrapElevator(received)

				setOperational(true, received.Elevator_ID)
				setConnected(true, received.Elevator_ID)

				//Test
				types.PrintOrders(all_elevators[test.Elevator_ID])

				fmt.Printf("\nAdded new elevator!\n")

			} else {
				//update states
				setFields(received.State, received.Floor, received.Direction, received.Elevator_ID)
				//fmt.Printf("\n\nrec2: ID: %s\n", received.Elevator_ID)
			}

			//update orders: Uncomment next two lines when merge is ready
			order_map, is_new_local_order := merger.MergeOrders(localelev_ID, getOrderMap(all_elevators), received.Orders)
			setFromOrderMap(order_map)
			if(is_new_local_order){ order_added <- true }

		case update := <-connectionUpdate: //Untested case
			if update.Connected {
				setConnected(true, update.Elevator_ID)
				fmt.Printf("The elevator is reconnected: %v\n", all_elevators[update.Elevator_ID].Connected)
			} else {
				
				setConnected(false, update.Elevator_ID)
				fmt.Printf("\n%v has been disconnected\n", update.Elevator_ID)

				orderReassigner(update.Elevator_ID, false)
				order_added <- true
				lamps.SetAllLamps(all_elevators[localelev_ID])

				//test
				fmt.Printf("\nThe elevator is connected: %v\n", all_elevators[update.Elevator_ID].Connected)
				types.PrintOrders(all_elevators[localelev_ID])
			}
			//TestPeersPrint()
		/*case update := <-operationUpdate: //Untested case
			if update.Is_Operational {
				setOperational(true, update.Elevator_ID)
			} else {
				fmt.Printf("\nThis was true\n")
				setOperational(false, update.Elevator_ID)
				orderReassigner(update.Elevator_ID, true)
				order_added <- true
				lamps.SetAllLamps(all_elevators[localelev_ID])
			}*/
		}
	}
}

func TransmitElev(elev_tx chan<- types.Wrapped_Elevator) {
	for {
		elev_tx <- wrapElevator(localelev_ID)
		//elev_tx <- all_elevators[localelev_ID]
		time.Sleep(time.Millisecond * constants.TRANSMIT_MS)
	}
}

func TestPrintAllElevators() {
	for {
		fmt.Printf("\n\n")
		//for key, val := range all_elevators {
			fmt.Printf("\nElev: %s\n", localelev_ID)
			types.PrintStates(all_elevators[localelev_ID])
			types.PrintOrders(all_elevators[localelev_ID])
		//}
		time.Sleep(time.Second * 10)
	}
}

func TestPeersPrint() {
	for {	
		//for key, val := range all_elevators {
			fmt.Printf("\nElev: %s\n", localelev_ID)
			/*if(val.Connected){
				fmt.Printf("Is connected")
			} else {
				fmt.Printf("Not connected")
			}*/
			types.PrintOrders(all_elevators[localelev_ID])
		//}
		time.Sleep(time.Second * 5)
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
	if keyExists(elevator_ID) {
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

func unwrapElevator(wrapped types.Wrapped_Elevator) types.Elevator {
	var unwrapped types.Elevator
	unwrapped.Elevator_ID = wrapped.Elevator_ID
	unwrapped.State = wrapped.State
	unwrapped.Floor = wrapped.Floor
	unwrapped.Direction = wrapped.Direction
	unwrapped.Orders = wrapped.Orders[wrapped.Elevator_ID]
	return unwrapped
}

func ReadLocalElevator() types.Elevator {
	return all_elevators[localelev_ID]
}

func ReadAllElevators() map[string]types.Elevator {
	return all_elevators
}

func ReadLocalElevatorID() string {
	return localelev_ID
}

func keyExists(key string) bool {
	_, exists := all_elevators[key]
	return exists
}

func isValidID(ID string) bool {
	//check substring stuff TODO
	return true
}

//Workaround functions because go does not allow setting structs in maps directly
func setID(ID string) {
	temp, is := all_elevators[ID]
	if is {
		temp.Elevator_ID = ID
		all_elevators[ID] = temp
	}
}

func setFields(s types.ElevatorState, f int, d elevio.MotorDirection, ID string) {
	setState(s, ID)
	setFloor(f, ID)
	setDir(d, ID)
}

func setState(s types.ElevatorState, ID string) {
	temp, is := all_elevators[ID]
	if is {
		temp.State = s
		all_elevators[ID] = temp
	} /*else {
		//FATAL //todo, log
	}*/
}

func setFloor(f int, ID string) {
	temp, is := all_elevators[ID]
	if is {
		temp.Floor = f
		all_elevators[ID] = temp
	}
}

func setDir(d elevio.MotorDirection, ID string) {
	temp, is := all_elevators[ID]
	if is {
		temp.Direction = d
		all_elevators[ID] = temp
	}
}

func setOperational(val bool, ID string) {
	temp, is := all_elevators[ID]
	if is {
		temp.Is_Operational = val
		all_elevators[ID] = temp
	}
}

func setConnected(val bool, ID string) {
	temp, is := all_elevators[ID]
	if is {
		temp.Connected = val
		all_elevators[ID] = temp
	}
}

func setOrdered(floor int, button int, ID string, accepted bool) {
	temp, is := all_elevators[ID]
	if is {
		if accepted {
			temp.Orders[floor][button].State = types.OS_AcceptedOrder
		} else {
			temp.Orders[floor][button].State = types.OS_UnacceptedOrder
		}

		//leave counter unchanged
		all_elevators[ID] = temp
	}
}

func setOrderList(list [constants.N_FLOORS][constants.N_BUTTONS]types.Order, ID string) {
	temp, is := all_elevators[ID]
	if is {
		temp.Orders = list
		all_elevators[ID] = temp
	}
}

//Untested function
//TODO: Move
func orderReassigner(faultyElevID string, operationError bool) {
	var e = all_elevators[faultyElevID]
	fmt.Printf("\norderReassigner started\n")
	for i := 0; i < constants.N_FLOORS; i++ {
		/*if i == 0 {   //Remove this after test
			fmt.Printf("\nOuter For loop started\n")
		} else if i == constants.N_FLOORS - 1 {
			fmt.Printf("\nOuter For loop ended\n")
		}*/
		for j := 0; j < constants.N_BUTTONS-1; j++ {
			/*if j == 0 {   //Remove this after test
				fmt.Printf("\nInner For loop started\n")
			} else if j == constants.N_BUTTONS - 2 {
				fmt.Printf("\nInner For loop ended\n")
			}*/
			if e.Orders[i][j].State == types.OS_AcceptedOrder {
				fmt.Printf("\nsetOrdered will run\n")
				setOrdered(i, j, localelev_ID, true)
				
			}
		}
	}

	if operationError {
		if e.Floor == constants.N_FLOORS && e.Direction == elevio.MD_Up {
			setOrdered(constants.N_FLOORS - 1, elevio.BT_Cab , faultyElevID, true) //Todo use cab cal
		} else if e.Floor == 0 && e.Direction == elevio.MD_Down {
			setOrdered(1, elevio.BT_Cab, faultyElevID, true) //Todo use cab cal
		} else {
			var dummyOrder = UpcomingFloor(e)
			setOrdered(dummyOrder, elevio.BT_Cab, faultyElevID, true) //Todo use cab cal
		}
	}
}

func UpcomingFloor(e types.Elevator) int {
	if e.Direction == elevio.MD_Up {
		return e.Floor + 1
	} else if e.Direction == elevio.MD_Down {
		return e.Floor - 1
	} else {
		return e.Floor
	}
}

//Returns a slice of the working elevators UNTESTED
func WorkingElevs( /*elevs map[string]types.Elevator  <- upgrade*/ ) []types.Elevator{
	var workingElevs []types.Elevator
	for _, v := range all_elevators { //change to elevs when you move variables
		if v.Is_Operational && v.Connected {
			workingElevs = append(workingElevs, v)
		}
	}
	return workingElevs
}

func orderReassignerTest(lostID string) {
	for i:=0; i < constants.N_BUTTONS; i++{
		for j:=0; j < constants.N_FLOORS; j++{
			setOrdered(j, i, lostID, true)
		}
	}
}