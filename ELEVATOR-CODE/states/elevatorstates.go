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

func InitElevators(localElevID string, drvFloors <-chan int) map[string]types.Elevator {
	allElevators := make(map[string]types.Elevator)

	//Initialize local elevator
	var emptyElevator types.Elevator
	emptyElevator.ElevatorID = localElevID
	allElevators[localElevID] = emptyElevator

	var ord [constants.N_FLOORS][constants.N_BUTTONS](types.Order)
	allElevators = setFields(types.ES_Idle, -1, elevio.MD_Stop, localElevID, allElevators)
	allElevators = setOrderList(ord, localElevID, allElevators)
	allElevators = setOperational(true, localElevID, allElevators)
	allElevators = setConnected(true, localElevID, allElevators)
	lamps.SetAllLamps(localElevID, allElevators)
	elevio.SetDoorOpenLamp(false)

	//wait to allow floor signal to arrive if we start on a floor
	time.Sleep(time.Millisecond * 50)
	select {
	case floor := <-drvFloors:
		allElevators = setFloor(floor, localElevID, allElevators)
		elevio.SetFloorIndicator(floor)
	default:
		//init between floors
		if allElevators[localElevID].Floor == -1 {
			allElevators = setDir(elevio.MD_Down, localElevID, allElevators)
			elevio.SetMotorDirection(allElevators[localElevID].Direction)
			allElevators = setState(types.ES_Moving, localElevID, allElevators)
		}
	}
	return allElevators
}

func UpdateElevator(
	localElevID string,
	/*I/O Channels for handling FSM*/
	updateState <-chan types.ElevatorState, updateFloor <-chan int, updateDirection <-chan elevio.MotorDirection,
	clearFloor <-chan int, floorReached chan<- bool, localOrderAdded chan<- bool, drvButton <-chan elevio.ButtonEvent,
	/*Channels for letting other modules read allElevators*/
	allElevsUpdateFSM chan<- types.Elevator, allElevsUpdateOperational chan<- map[string]types.Elevator,
	/*I/O channels for interface/communicating with other elevators*/
	elevRx <-chan types.WrappedElevator, elevTx chan<- types.WrappedElevator, connectionUpdate <-chan types.ConnectionEvent, operationUpdate <-chan types.OperationEvent) {
	
	tick := time.NewTicker(time.Millisecond*constants.TRANSMIT_MS)
	allElevators := InitElevators(localElevID, updateFloor)
	allElevsUpdateFSM <- allElevators[localElevID]
	allElevsUpdateOperational <- allElevators

	for {
		select {
		case newState := <- updateState:
			fmt.Printf("\nsent state\n")

			allElevators = setState(newState, localElevID, allElevators)
			allElevsUpdateFSM <- allElevators[localElevID]
			allElevsUpdateOperational <- allElevators

		case newFloor := <-updateFloor:
			fmt.Printf("\nrec drv floor\n")

			allElevators = setFloor(newFloor, localElevID, allElevators)

			allElevsUpdateFSM <- allElevators[localElevID]
			allElevsUpdateOperational <- allElevators
			floorReached <- true

		case newDir := <-updateDirection:
			fmt.Printf("\nsent dir\n")

			allElevators = setDir(newDir, localElevID, allElevators)
			allElevsUpdateFSM <- allElevators[localElevID]
			allElevsUpdateOperational <- allElevators

		case buttonPress := <-drvButton:
			fmt.Printf("\nsent button\n")

			assignedOrder := orders.AssignOrder(allElevators, localElevID, buttonPress)
			allElevators = setOrdered(assignedOrder.Order.Floor, int(assignedOrder.Order.Button), assignedOrder.ElevatorID, false, allElevators)
			allElevsUpdateFSM <- allElevators[localElevID]
			allElevsUpdateOperational <- allElevators

			if assignedOrder.ElevatorID == localElevID && assignedOrder.Order.Button == elevio.BT_Cab {
				localOrderAdded <- true
				lamps.SetAllLamps(localElevID, allElevators)
			}

		case <- clearFloor:
			fmt.Printf("\nsent cl\n")

			allElevators = setOrderList(orders.ClearAtCurrentFloor(allElevators[localElevID]), localElevID, allElevators)
			allElevsUpdateFSM <- allElevators[localElevID]
			allElevsUpdateOperational <- allElevators

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
				allElevators = setOperational(true, received.ElevatorID, allElevators)
				allElevators = setConnected(true, received.ElevatorID, allElevators)
				fmt.Printf("\nAdded new elevator!\n")
			} else {
				allElevators = setFields(received.State, received.Floor, received.Direction, received.ElevatorID, allElevators)
			}

			//Update orders
			orderMap, isNewLocalOrder, shouldLight := orders.MergeOrders(localElevID, getOrderMap(allElevators), received.Orders)
			allElevators = setFromOrderMap(orderMap, allElevators)
			allElevsUpdateFSM <- allElevators[localElevID]
			allElevsUpdateOperational <- allElevators

			if(isNewLocalOrder){
				localOrderAdded <- true 
			}
			if(shouldLight){
				lamps.SetAllLamps(localElevID, allElevators)
			}

		case update := <-connectionUpdate:
			allElevators = setConnected(update.IsConnected, update.ElevatorID, allElevators)
			allElevsUpdateFSM <- allElevators[localElevID]
			allElevsUpdateOperational <- allElevators

			if !update.IsConnected {
				fmt.Printf("\n%v has been disconnected\n", update.ElevatorID)

				elev, isNewLocalOrder := orders.OrderReassigner(update.ElevatorID, localElevID, allElevators)
				allElevators = setOrderList(elev.Orders, localElevID, allElevators)
				allElevsUpdateFSM <- allElevators[localElevID]
				allElevsUpdateOperational <- allElevators

				if isNewLocalOrder {
					localOrderAdded <- true
				}
				lamps.SetAllLamps(localElevID, allElevators)
			}

		case update := <-operationUpdate:
			allElevators = setOperational(update.IsOperational, update.ElevatorID, allElevators)
			allElevsUpdateFSM <- allElevators[localElevID]
			allElevsUpdateOperational <- allElevators

			if !update.IsOperational && update.ElevatorID != localElevID { 
				fmt.Printf("\nID: %v has been marked as not operational\n", update.ElevatorID)
				elev, isNewLocalOrder := orders.OrderReassigner(update.ElevatorID, localElevID, allElevators)
				allElevators = setOrderList(elev.Orders, localElevID, allElevators)
				allElevsUpdateFSM <- allElevators[localElevID]
				allElevsUpdateOperational <- allElevators

				if isNewLocalOrder {
					localOrderAdded <- true
				}
				lamps.SetAllLamps(localElevID, allElevators)
			}
		case <- tick.C:
			//transmit local elevator over UDP
			elevTx <- wrapElevator(localElevID, allElevators)
		}
	}
}

func wrapElevator(elevatorID string, elevators map[string]types.Elevator) types.WrappedElevator {
	var wrapped types.WrappedElevator
	if _, exists := elevators[elevatorID] ; exists  {
		temp := elevators[elevatorID] 
		wrapped.ElevatorID = temp.ElevatorID
		wrapped.State = temp.State
		wrapped.Floor = temp.Floor
		wrapped.Direction = temp.Direction
		wrapped.Orders = getOrderMap(elevators)
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

func getOrderMap(elevatorMap map[string]types.Elevator) map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order {
	orderMap := make(map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order)
	for key, val := range elevatorMap {
		orderMap[key] = val.Orders
	}
	return orderMap
}

func setFromOrderMap(orderMap map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order, elevators map[string]types.Elevator) map[string]types.Elevator {
	for key, val := range orderMap {
		elevators = setOrderList(val, key, elevators)
	}
	return elevators
}

/*Workaround functions because go does not allow setting structs in maps directly*/
func setID(ID string, elevators map[string]types.Elevator) map[string]types.Elevator {
	temp, is := elevators[ID]
	if is {
		temp.ElevatorID = ID
		elevators[ID] = temp
	}
	return elevators
}

func setFields(s types.ElevatorState, f int, d elevio.MotorDirection, ID string, elevators map[string]types.Elevator) map[string]types.Elevator {
	elevators = setState(s, ID, elevators)
	elevators = setFloor(f, ID, elevators)
	elevators = setDir(d, ID, elevators)
	return elevators
}

func setState(s types.ElevatorState, ID string, elevators map[string]types.Elevator) map[string]types.Elevator  {
	temp, is := elevators[ID]
	if is {
		temp.State = s
		elevators[ID] = temp
	} 
	return elevators
}

func setFloor(f int, ID string, elevators map[string]types.Elevator) map[string]types.Elevator {
	temp, is := elevators[ID]
	if is {
		temp.Floor = f
		elevators[ID] = temp
	}
	return elevators
}

func setDir(d elevio.MotorDirection, ID string, elevators map[string]types.Elevator) map[string]types.Elevator {
	temp, is := elevators[ID]
	if is {
		temp.Direction = d
		elevators[ID] = temp
	}
	return elevators
}

func setOperational(val bool, ID string, elevators map[string]types.Elevator) map[string]types.Elevator {
	temp, is := elevators[ID]
	if is {
		temp.IsOperational = val
		elevators[ID] = temp
	}
	return elevators
}

func setConnected(val bool, ID string, elevators map[string]types.Elevator) map[string]types.Elevator {
	temp, is := elevators[ID]
	if is {
		temp.IsConnected = val
		elevators[ID] = temp
	}
	return elevators
}

func setOrdered(floor int, button int, ID string, accepted bool, elevators map[string]types.Elevator) map[string]types.Elevator {
	temp, is := elevators[ID]
	if is {
		if accepted {
			temp.Orders[floor][button].State = types.OS_AcceptedOrder
		} else {
			temp.Orders[floor][button].State = types.OS_UnacceptedOrder
		}
		elevators[ID] = temp
	}
	return elevators
}

func setOrderList(list [constants.N_FLOORS][constants.N_BUTTONS]types.Order, ID string, elevators map[string]types.Elevator) map[string]types.Elevator {
	temp, is := elevators[ID]
	if is {
		temp.Orders = list
		elevators[ID] = temp
	}
	return elevators
}