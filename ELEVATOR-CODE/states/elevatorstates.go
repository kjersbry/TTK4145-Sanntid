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

var all_elevators map[string]types.Elevator //todo
var localelev_ID string


func InitElevators(localID string, drvFloors <-chan int) {
	localelev_ID = localID
	all_elevators = make(map[string]types.Elevator)

	//Initialize local elevator
	var emptyElevator types.Elevator
	emptyElevator.ElevatorID = localelev_ID
	all_elevators[localelev_ID] = emptyElevator

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
	case floor := <-drvFloors:
		setFloor(floor, localelev_ID)
		elevio.SetFloorIndicator(floor)
	default:
		//init between floors
		if all_elevators[localelev_ID].Floor == -1 {
			setDir(elevio.MD_Down, localelev_ID)
			elevio.SetMotorDirection(all_elevators[localelev_ID].Direction)
			setState(types.ES_Moving, localelev_ID)
		}
	}
}


/*This is the only routine that is writing to all_elevators*/ //todo remove comment
func UpdateElevator(
	/*I/O Channels for handling FSM*/
	updateState <-chan types.ElevatorState, updateFloor <-chan int, updateDirection <-chan elevio.MotorDirection,
	clearFloor <-chan int, floorReached chan<- bool, localOrderAdded chan<- bool, drvButton <-chan elevio.ButtonEvent,
	/*I/O channels for interface/communicating with other elevators*/
	elevRx <-chan types.WrappedElevator, elevTx chan<- types.WrappedElevator, connectionUpdate <-chan types.ConnectionEvent, operationUpdate <-chan types.OperationEvent) {
	
	tick := time.NewTicker(time.Millisecond*constants.TRANSMIT_MS)
	for {
		select {
		case newState := <- updateState:
			setState(newState, localelev_ID)

		case newFloor := <-updateFloor:
			setFloor(newFloor, localelev_ID)
			floorReached <- true

		case newDir := <-updateDirection:
			setDir(newDir, localelev_ID)

		case buttonPress := <-drvButton:
			assignedOrder := orders.AssignOrder(all_elevators, localelev_ID, buttonPress)
			setOrdered(assignedOrder.Order.Floor, int(assignedOrder.Order.Button), assignedOrder.ElevatorID, false)
			if assignedOrder.ElevatorID == localelev_ID && assignedOrder.Order.Button == elevio.BT_Cab {
				localOrderAdded <- true
				lamps.SetAllLamps(localelev_ID, all_elevators)
			}
		case <- clearFloor:
			setOrderList(orders.ClearAtCurrentFloor(all_elevators[localelev_ID]), localelev_ID)
			lamps.SetAllLamps(localelev_ID, all_elevators)

		case received := <- elevRx:
			if received.ElevatorID == localelev_ID {
				break
			}
			if received.ElevatorID == "NOTFOUND" {
				break
			}

			if !keyExists(received.ElevatorID) {
				all_elevators[received.ElevatorID] = unwrapElevator(received)
				setOperational(true, received.ElevatorID)
				setConnected(true, received.ElevatorID)
				fmt.Printf("\nAdded new elevator!\n")
			} else {
				setFields(received.State, received.Floor, received.Direction, received.ElevatorID)
			}

			//Update orders
			orderMap, isNewLocalOrder, shouldLight := orders.MergeOrders(localelev_ID, getOrderMap(all_elevators), received.Orders)
			setFromOrderMap(orderMap)
			if(isNewLocalOrder){
				localOrderAdded <- true 
			}
			if(shouldLight){
				lamps.SetAllLamps(localelev_ID, all_elevators)
			}

		case update := <-connectionUpdate:
			setConnected(update.IsConnected, update.ElevatorID)
			if !update.IsConnected {
				fmt.Printf("\n%v has been disconnected\n", update.ElevatorID)

				elev, isNewLocalOrder := orders.OrderReassigner(update.ElevatorID, localelev_ID, all_elevators)
				setOrderList(elev.Orders, localelev_ID)
				if isNewLocalOrder {
					localOrderAdded <- true
				}
				lamps.SetAllLamps(localelev_ID, all_elevators)
			}

		case update := <-operationUpdate:
			setOperational(update.IsOperational, update.ElevatorID)

			if !update.IsOperational && update.ElevatorID != localelev_ID { 
				fmt.Printf("\nID: %v has been marked as not operational\n", update.ElevatorID)
				elev, isNewLocalOrder := orders.OrderReassigner(update.ElevatorID, localelev_ID, all_elevators)
				setOrderList(elev.Orders, localelev_ID)
				if isNewLocalOrder {
					localOrderAdded <- true
				}
				lamps.SetAllLamps(localelev_ID, all_elevators)
			}
		case <- tick.C:
			//transmit local elevator over UDP
			elevTx <- wrapElevator(localelev_ID)
		}
	}
}

func wrapElevator(elevatorID string) types.WrappedElevator {
	var wrapped types.WrappedElevator
	if keyExists(elevatorID) {
		temp := all_elevators[elevatorID] 
		wrapped.ElevatorID = temp.ElevatorID
		wrapped.State = temp.State
		wrapped.Floor = temp.Floor
		wrapped.Direction = temp.Direction
		wrapped.Orders = getOrderMap(all_elevators)
	} else {
		wrapped.ElevatorID = "NOTFOUND"
	}
	return wrapped
}

func unwrapElevator(wrapped types.WrappedElevator) types.Elevator {
	var unwrapped types.Elevator
	unwrapped.ElevatorID = wrapped.ElevatorID
	unwrapped.State = wrapped.State
	unwrapped.Floor = wrapped.Floor
	unwrapped.Direction = wrapped.Direction
	unwrapped.Orders = wrapped.Orders[wrapped.ElevatorID]
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

func getOrderMap(elevatorMap map[string]types.Elevator) map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order {
	orderMap := make(map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order)
	for key, val := range elevatorMap {
		orderMap[key] = val.Orders
	}
	return orderMap
}

func setFromOrderMap(orderMap map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order) {
	for key, val := range orderMap {
		setOrderList(val, key)
	}
}

/*Workaround functions because go does not allow setting structs in maps directly*/
func setID(ID string) {
	temp, is := all_elevators[ID]
	if is {
		temp.ElevatorID = ID
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
	} 
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
		temp.IsOperational = val
		all_elevators[ID] = temp
	}
}

func setConnected(val bool, ID string) {
	temp, is := all_elevators[ID]
	if is {
		temp.IsConnected = val
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