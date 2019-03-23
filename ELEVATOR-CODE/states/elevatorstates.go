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

func UpdateElevator(
	/*FSM channels, single elevator*/
	update_state <-chan types.ElevatorState, update_floor <-chan int, update_direction <-chan elevio.MotorDirection,
	clear_floor <-chan int, floor_reached chan<- bool, order_added chan<- int,
	/*---*/
	add_order <-chan types.AssignedOrder, elev_rx <-chan types.Wrapped_Elevator){
		for {
    	select{
      case new_state:= <- update_state:
				setState(new_state, localelev_ID)
				//types.PrintStates(all_elevs[localelev_ID])

      case new_floor:= <- update_floor:
				setFloor(new_floor, localelev_ID)
				floor_reached <- true
				fmt.Printf("\n elev floor set: %d\n", new_floor)

			case new_dir:= <- update_direction:
				setDir(new_dir, localelev_ID)

				//types.PrintStates(all_elevs[localelev_ID])

			case order:= <-add_order:
				ord := types.Order{-1/*todo: ORDER-ID*/, types.OS_UnacceptedOrder}
				setOrder(ord, order.Order.Floor, int(order.Order.Button), order.Elevator_ID)

				if(order.Elevator_ID == localelev_ID){
					order_added <- order.Order.Floor
					fmt.Printf("\nAdded order fl. %d\n", order.Order.Floor)
					lamps.SetAllLamps(all_elevs[localelev_ID]) //todo
				}
			case <- clear_floor:
				setOrderList(orders.ClearAtCurrentFloor(all_elevs[localelev_ID]), localelev_ID)
				lamps.SetAllLamps(all_elevs[localelev_ID])

			case received := <- elev_rx:
				fmt.Printf("\n\nID: %s\nDir: %s\nOrder1-cab: %s\n", received.Elevator_ID,
					types.DirToString(received.Direction), types.OrderToString(received.Orders["hei1"][0][0]))
				/*if(received.Elev_ID == all_elevs.Elev_ID){
					break
				}
				for key, val := range received.Elev_map {
					//fmt.Println("Ding")
					if !keyExists(key) {
						all_elevs.Elev_map[key] = val
						fmt.Printf("\nAdded new elev!\n")
					} else {
						//update states
						if received.Elev_ID == key {//overwrite new states where elev = elev which sent
							setFields(val.State, val.Floor, val.Direction, key)
						}
					}
				}*/
				//merge orders
				//funksjon for det
				//her, eller ett hakk ut?
			}
    }
}

func TestPrintAllElevators(){
	for{
		fmt.Printf("\n\n")
		for key, _ := range all_elevs {
			/*TEST:*/
			fmt.Printf("\nElev: %s", key)
			//types.PrintStates(val)
		}
		time.Sleep(time.Second*5)
	}
}

//tar inn ID --> kan lett requeste og wrappe annen heis hvis det trengs
func wrapElevator(elevator_ID string) types.Wrapped_Elevator {
	temp := all_elevs[elevator_ID] // = the non-wrapped elev
	var wrapped types.Wrapped_Elevator
	wrapped.Elevator_ID = temp.Elevator_ID
	wrapped.State = temp.State
	wrapped.Floor = temp.Floor
	wrapped.Direction = temp.Direction

	//wrap orders
	for key, val := range all_elevs {
		wrapped.Orders[key] = val.Orders
	}
	return wrapped
}

func TransmitElev(elev_tx chan<- types.Wrapped_Elevator){
	for {
			elev_tx <- wrapElevator(localelev_ID)
			time.Sleep(time.Millisecond*constants.TRANSMIT_MS)
	}
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


func setOrder(order types.Order, floor int, button int, ID string){
	temp, is := all_elevs[ID]
	if(is){
		temp.Orders[floor][button] = order
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
