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

//Only UpdateElevators() is writing to allElevators
var allElevators map[string]types.Elevator
var localElevID string


func InitElevators(localID string, drvFloors <-chan int) {
	localElevID = localID
	allElevators = make(map[string]types.Elevator)

	//Initialize local elevator
	var emptyElevator types.Elevator
	emptyElevator.ElevatorID = localElevID
	allElevators[localElevID] = emptyElevator

	var ord [constants.N_FLOORS][constants.N_BUTTONS](types.Order)
	setFields(types.ES_Idle, -1, elevio.MD_Stop, localElevID)
	setOrderList(ord, localElevID)
	setOperational(true, localElevID)
	setConnected(true, localElevID)
	lamps.SetAllLamps(localElevID, allElevators)
	elevio.SetDoorOpenLamp(false)

	//wait to allow floor signal to arrive if we start on a floor
	time.Sleep(time.Millisecond * 50)
	select {
	case floor := <-drvFloors:
		setFloor(floor, localElevID)
		elevio.SetFloorIndicator(floor)
	default:
		//init between floors
		if allElevators[localElevID].Floor == -1 {
			setDir(elevio.MD_Down, localElevID)
			elevio.SetMotorDirection(allElevators[localElevID].Direction)
			setState(types.ES_Moving, localElevID)
		}
	}
}


/*This is the only routine that is writing to allElevators*/
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
			setState(newState, localElevID)

		case newFloor := <-updateFloor:
			setFloor(newFloor, localElevID)
			floorReached <- true

		case newDir := <-updateDirection:
			setDir(newDir, localElevID)

		case buttonPress := <-drvButton:
			assignedOrder := orders.AssignOrder(allElevators, localElevID, buttonPress)
			setOrdered(assignedOrder.Order.Floor, int(assignedOrder.Order.Button), assignedOrder.ElevatorID, false)
			if assignedOrder.ElevatorID == localElevID && assignedOrder.Order.Button == elevio.BT_Cab {
				localOrderAdded <- true
				lamps.SetAllLamps(localElevID, allElevators)
			}
		case <- clearFloor:
			setOrderList(orders.ClearAtCurrentFloor(allElevators[localElevID]), localElevID)
			lamps.SetAllLamps(localElevID, allElevators)

		case received := <- elevRx:
			if received.ElevatorID == localElevID {
				break
			}
			if received.ElevatorID == "NOTFOUND" {
				break
			}

			if _, exists := allElevators[received.ElevatorID] ; !exists {
				allElevators[received.ElevatorID] = unwrapElevator(received)
				setOperational(true, received.ElevatorID)
				setConnected(true, received.ElevatorID)
				fmt.Printf("\nAdded new elevator!\n")
			} else {
				setFields(received.State, received.Floor, received.Direction, received.ElevatorID)
			}

			//Update orders
			orderMap, isNewLocalOrder, shouldLight := orders.MergeOrders(localElevID, getOrderMap(allElevators), received.Orders)
			setFromOrderMap(orderMap)
			if(isNewLocalOrder){
				localOrderAdded <- true 
			}
			if(shouldLight){
				lamps.SetAllLamps(localElevID, allElevators)
			}

		case update := <-connectionUpdate:
			setConnected(update.IsConnected, update.ElevatorID)
			if !update.IsConnected {
				fmt.Printf("\n%v has been disconnected\n", update.ElevatorID)

				elev, isNewLocalOrder := orders.OrderReassigner(update.ElevatorID, localElevID, allElevators)
				setOrderList(elev.Orders, localElevID)
				if isNewLocalOrder {
					localOrderAdded <- true
				}
				lamps.SetAllLamps(localElevID, allElevators)
			}

		case update := <-operationUpdate:
			setOperational(update.IsOperational, update.ElevatorID)

			if !update.IsOperational && update.ElevatorID != localElevID { 
				fmt.Printf("\nID: %v has been marked as not operational\n", update.ElevatorID)
				elev, isNewLocalOrder := orders.OrderReassigner(update.ElevatorID, localElevID, allElevators)
				setOrderList(elev.Orders, localElevID)
				if isNewLocalOrder {
					localOrderAdded <- true
				}
				lamps.SetAllLamps(localElevID, allElevators)
			}
		case <- tick.C:
			//transmit local elevator over UDP
			elevTx <- wrapElevator(localElevID)
		}
	}
}

func wrapElevator(elevatorID string) types.WrappedElevator {
	var wrapped types.WrappedElevator
	if _, exists := allElevators[elevatorID] ; exists {
		temp := allElevators[elevatorID] 
		wrapped.ElevatorID = temp.ElevatorID
		wrapped.State = temp.State
		wrapped.Floor = temp.Floor
		wrapped.Direction = temp.Direction
		wrapped.Orders = getOrderMap(allElevators)
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
	return allElevators[localElevID]
}

func ReadAllElevators() map[string]types.Elevator {
	return allElevators
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
	temp, is := allElevators[ID]
	if is {
		temp.ElevatorID = ID
		allElevators[ID] = temp
	}
}

func setFields(s types.ElevatorState, f int, d elevio.MotorDirection, ID string) {
	setState(s, ID)
	setFloor(f, ID)
	setDir(d, ID)
}

func setState(s types.ElevatorState, ID string) {
	temp, is := allElevators[ID]
	if is {
		temp.State = s
		allElevators[ID] = temp
	} 
}

func setFloor(f int, ID string) {
	temp, is := allElevators[ID]
	if is {
		temp.Floor = f
		allElevators[ID] = temp
	}
}

func setDir(d elevio.MotorDirection, ID string) {
	temp, is := allElevators[ID]
	if is {
		temp.Direction = d
		allElevators[ID] = temp
	}
}

func setOperational(val bool, ID string) {
	temp, is := allElevators[ID]
	if is {
		temp.IsOperational = val
		allElevators[ID] = temp
	}
}

func setConnected(val bool, ID string) {
	temp, is := allElevators[ID]
	if is {
		temp.IsConnected = val
		allElevators[ID] = temp
	}
}

func setOrdered(floor int, button int, ID string, accepted bool) {
	temp, is := allElevators[ID]
	if is {
		if accepted {
			temp.Orders[floor][button].State = types.OS_AcceptedOrder
		} else {
			temp.Orders[floor][button].State = types.OS_UnacceptedOrder
		}
		allElevators[ID] = temp
	}
}

func setOrderList(list [constants.N_FLOORS][constants.N_BUTTONS]types.Order, ID string) {
	temp, is := allElevators[ID]
	if is {
		temp.Orders = list
		allElevators[ID] = temp
	}
}