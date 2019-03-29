package states

import (
	"fmt"
	"time"

	"../constants"
	"../elevio"
	"../lamps"
	"../orders"
	"../types"
	"../merger"
)

var all_elevators map[string]types.Elevator
var localelev_ID string

/*Måå fikses, fikk concurrent map read and map write*/

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
	lamps.SetAllLamps(localelev_ID, all_elevators)
	elevio.SetDoorOpenLamp(false)
	//wait to allow floor signal to arrive if we start on a floor
	time.Sleep(time.Millisecond * 50)
	select {
	case floor := <-drv_floors:
		setFloor(floor, localelev_ID)
		elevio.SetFloorIndicator(floor)
		fmt.Printf("\nhei\n")
	default:
		if all_elevators[localelev_ID].Floor == -1 {
			//init between floors:
			setDir(elevio.MD_Down, localelev_ID)
			elevio.SetMotorDirection(all_elevators[localelev_ID].Direction)
			setState(types.ES_Moving, localelev_ID)
			fmt.Printf("\nhei2\n")
			/*//forslag for å fikse hvis den kjører ned etter init:
			send <- start timer signal 
			Ta i mot i et event
			Der skal heisen snus.*/
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
	add_order <-chan types.AssignedOrder, elev_rx <-chan types.Wrapped_Elevator, elev_tx chan<- types.Wrapped_Elevator, connectionUpdate <-chan types.Connection_Event, operationUpdate <-chan types.Operation_Event) {
	
	tick := time.NewTicker(time.Millisecond*constants.TRANSMIT_MS)
	for {
		select {
		case new_state := <- update_state:
			setState(new_state, localelev_ID)
			fmt.Printf("\nState was set to: %s\n", types.StateToString(new_state))
			//types.PrintStates(all_elevators[localelev_ID])

		case new_floor := <-update_floor:
			setFloor(new_floor, localelev_ID)
			floor_reached <- true
			//fmt.Printf("\n elev floor set: %d\n", new_floor)

		case new_dir := <-update_direction:
			setDir(new_dir, localelev_ID)

			//types.PrintStates(all_elevators[localelev_ID])

		case order := <- add_order:
			setOrdered(order.Order.Floor, int(order.Order.Button), order.Elevator_ID, false)
			//lamps.SetAllLamps(localelev_ID, all_elevators)
			if order.Elevator_ID == localelev_ID && order.Order.Button == elevio.BT_Cab {
				order_added <- true
				lamps.SetAllLamps(localelev_ID, all_elevators)
			}
		case <- clear_floor:
			setOrderList(orders.ClearAtCurrentFloor(all_elevators[localelev_ID]), localelev_ID)
			lamps.SetAllLamps(localelev_ID, all_elevators)

		case received := <- elev_rx:

			if received.Elevator_ID == localelev_ID {
				break
			}
			if received.Elevator_ID == "NOTFOUND" || !isValidID(received.Elevator_ID) {
				//log/print: something weird has happened
				break
			}

			if !keyExists(received.Elevator_ID) {
				all_elevators[received.Elevator_ID] = unwrapElevator(received)

				setOperational(true, received.Elevator_ID)
				setConnected(true, received.Elevator_ID)

				fmt.Printf("\nAdded new elevator!\n")

			} else {
				//update states
				setFields(received.State, received.Floor, received.Direction, received.Elevator_ID)
			}

			//update orders
			order_map, is_new_local_order, should_light := merger.MergeOrders(localelev_ID, getOrderMap(all_elevators), received.Orders)
			setFromOrderMap(order_map)
			if(is_new_local_order){
				order_added <- true 
				lamps.SetAllLamps(localelev_ID, all_elevators)
			}
			if(should_light){
				lamps.SetAllLamps(localelev_ID, all_elevators)
			}

		case update := <-connectionUpdate: //Untested case
				setConnected(update.Connected, update.Elevator_ID)
			if !update.Connected {
				
				fmt.Printf("\n%v has been disconnected\n", update.Elevator_ID)

				is_new_local_order := orderReassigner(update.Elevator_ID, false)
				if is_new_local_order {
					order_added <- true
				}
				lamps.SetAllLamps(localelev_ID, all_elevators)

				//test
				fmt.Printf("\nThe elevator is connected: %v\n", all_elevators[update.Elevator_ID].Connected)
				types.PrintOrders(all_elevators[localelev_ID])
			}
			//TestPeersPrint()
		case update := <-operationUpdate:
			setOperational(update.Is_Operational, update.Elevator_ID)

			if !update.Is_Operational && update.Elevator_ID != localelev_ID { 
				fmt.Printf("\nID: %v has been marked as !operational\n", update.Elevator_ID)
				is_new_local_order := orderReassigner(update.Elevator_ID, true)
				if is_new_local_order {
					order_added <- true
				}
				lamps.SetAllLamps(localelev_ID, all_elevators)
			}
		case <- tick.C:
			elev_tx <- wrapElevator(localelev_ID)
		}
	}
}

/*
func TransmitElev() {
		
		time.Sleep(time.Millisecond * constants.TRANSMIT_MS)
}*/

func TestPrintAllElevators() {
	for {
		fmt.Printf("\n\n")
		for key, val := range all_elevators {
			/*fmt.Printf("\nElev: %s\n", localelev_ID)
			types.PrintStates(all_elevators[localelev_ID])
			types.PrintOrders(all_elevators[localelev_ID])*/
			fmt.Printf("\nElev: %s\n", key)
			//types.PrintStates(val)
			types.PrintOrders(val)
		}
		time.Sleep(time.Second * 2)
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
func orderReassigner(faultyElevID string, operationError bool) bool {
	fmt.Printf("\nOrderReassigner was called\n")
	var e = all_elevators[faultyElevID]
	is_new_local_order := false
	for i := 0; i < constants.N_FLOORS; i++ {
		for j := 0; j < constants.N_BUTTONS-1; j++ {
			if e.Orders[i][j].State == types.OS_AcceptedOrder {
				fmt.Printf("\nsetOrdered will run\n")
				setOrdered(i, j, localelev_ID, true)
				is_new_local_order = true
			}
		}
	}

	if operationError {
		if e.Floor == constants.N_FLOORS && e.Direction == elevio.MD_Up {
			setOrdered(constants.N_FLOORS - 1, elevio.BT_Cab , faultyElevID, true) //Todo use cab cal
		} else if e.Floor == 0 && e.Direction == elevio.MD_Down {
			setOrdered(1, elevio.BT_Cab, faultyElevID, true) //Todo use cab cal
		} else {
			if e.Direction == elevio.MD_Stop {
				if e.Floor == 0 {
					e.Direction = elevio.MD_Up
				} else {
					e.Direction = elevio.MD_Down
				}
			}
			var dummyOrder = UpcomingFloor(e) //Forenkle denne?
			setOrdered(dummyOrder, elevio.BT_Cab, faultyElevID, true) //Todo use cab cal
		}
	}
	return is_new_local_order
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

func PrintCabs() {
	for {
		for key, val:= range all_elevators {
			fmt.Printf("\nCab call of %s, at floor 1 is %s\n", key, types.OrderToString(val.Orders[0][2]))
		}
		time.Sleep(time.Second*5)
	}

}