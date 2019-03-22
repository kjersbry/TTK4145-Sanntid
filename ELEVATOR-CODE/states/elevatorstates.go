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
var Local_elev_ID string

func InitSynchElevators(drv_floors <-chan int){
	all_elevs.Elev_map = make(map[string]types.Elevator)
	e := InitElevator()

	all_elevs.Elev_ID = e.Elevator_ID
	all_elevs.Elev_map[e.Elevator_ID] = e

	select{
	case floor:= <- drv_floors:
			setFloor(floor, e.Elevator_ID)
			fmt.Printf("\nhei\n")
	default:
			if(all_elevs.Elev_map[e.Elevator_ID].Floor == -1){
					//init between floors:
					setDir(elevio.MD_Down, e.Elevator_ID)
					elevio.SetMotorDirection(all_elevs.Elev_map[e.Elevator_ID].Direction)
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
	add_order <-chan types.AssignedOrder, elev_rx <-chan types.Elevator_test){
		for {
    	select{
      case new_state:= <- update_state:
				setState(new_state, all_elevs.Elev_ID)
				//types.PrintStates(all_elevs.Elev_map[all_elevs.Elev_ID])

      case new_floor:= <- update_floor:
				setFloor(new_floor, all_elevs.Elev_ID)
				floor_reached <- true
				fmt.Printf("\n elev floor set: %d\n", new_floor)

			case new_dir:= <- update_direction:
				setDir(new_dir, all_elevs.Elev_ID)

				//types.PrintStates(all_elevs.Elev_map[all_elevs.Elev_ID])

			case order:= <-add_order:
				ord := types.Order{-1/*todo: ORDER-ID*/, types.OS_UnacceptedOrder}
				setOrder(ord, order.Order.Floor, int(order.Order.Button), order.Elevator_ID)

				if(order.Elevator_ID == all_elevs.Elev_ID){
					order_added <- order.Order.Floor
					fmt.Printf("\nAdded order fl. %d\n", order.Order.Floor)
					lamps.SetAllLamps(all_elevs.Elev_map[all_elevs.Elev_ID]) //todo
				}
			case <- clear_floor:
				setOrderList(orders.ClearAtCurrentFloor(all_elevs.Elev_map[all_elevs.Elev_ID]), all_elevs.Elev_ID)
				lamps.SetAllLamps(all_elevs.Elev_map[all_elevs.Elev_ID])

			case received := <- elev_rx:
				fmt.Printf("\n\nID: %s\nDir: %s\nOrder1-cab: %s\n", received.Elevator_ID,
					types.DirToString(received.Direction), types.OrderToString(received.Orders[0][0][0]))
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
		//for key, val := range all_elevs.Elev_map {
		for key, _ := range all_elevs.Elev_map {
			/*TEST:*/
			fmt.Printf("\nElev: %s", key)
			//types.PrintStates(val)
		}
		time.Sleep(time.Second*5)
	}
}

//tar inn ID --> kan lett requeste og wrappe annen heis hvis det trengs
func wrap_elevator(elevator_ID string) Wrapped_Elevator {
	temp := all_elevs[elevator_ID] // = non-wrapped elev
	var wrapped types.Wrapped_Elevator
	wrapped.Elevator_ID = temp.Elevator_ID
	wrapped.State = temp.State
	wrapped.Floor = temp.Floor
	wrapped.Direction = temp.Direction

	//wrap orders
	for key, val := range all_elevs {
		//wrapped.Orders[i] = val
	}
}

func TransmitElev(elev_tx chan<- types.Elevator_test){
	//bruk wrapper(LOCAL_ID)

	temp := all_elevs.Elev_map[all_elevs.Elev_ID]
	var test_elev types.Elevator_test
	test_elev.Elevator_ID = temp.Elevator_ID

	for {
		temp = all_elevs.Elev_map[all_elevs.Elev_ID]

		test_elev.State = temp.State
		test_elev.Floor = temp.Floor
		test_elev.Direction = temp.Direction
		for i := 0; i < constants.N_ELEVATORS; i++{
			test_elev.Orders[i] = temp.Orders
		}
	//	for e := 0; e < constants.N_ELEVATORS - 1; e ++{
			//fmt.Printf("\n\nID: %s\nOrders dim: %dx%dx%d\n\n", test_elev.Elevator_ID, len(test_elev.Orders),
						//len(test_elev.Orders[1]), len(test_elev.Orders[1][1]))
			elev_tx <- test_elev
			time.Sleep(time.Millisecond*constants.TRANSMIT_MS)
	//	}
	}
}


func ReadElevator() types.Elevator {
	return all_elevs.Elev_map[all_elevs.Elev_ID]
}


func keyExists(key string) bool {
	_, exists := all_elevs.Elev_map[key]
	return exists
}

//Workaround functions because go does not allow setting structs in maps directly
func setFields(s types.ElevatorState,f int, d elevio.MotorDirection, ID string){
	setState(s, ID)
	setFloor(f, ID)
	setDir(d, ID)
}

func setState(s types.ElevatorState, ID string) {
	temp, is := all_elevs.Elev_map[ID]
	if(is){
		temp.State = s
		all_elevs.Elev_map[ID] = temp
	} /*else {
		//FATAL //todo, log
	}*/
}

func setFloor(f int, ID string){
	temp, is := all_elevs.Elev_map[ID]
	if(is){
		temp.Floor = f
		all_elevs.Elev_map[ID] = temp
	}
}

func setDir(d elevio.MotorDirection, ID string){
	temp, is := all_elevs.Elev_map[ID]
	if(is){
		temp.Direction = d
		all_elevs.Elev_map[ID] = temp
	}
}


func setOrder(order types.Order, floor int, button int, ID string){
	temp, is := all_elevs.Elev_map[ID]
	if(is){
		temp.Orders[floor][button] = order
		all_elevs.Elev_map[ID] = temp
	}
}

func setOrderList(list [constants.N_FLOORS][constants.N_BUTTONS]types.Order, ID string) {
	temp, is := all_elevs.Elev_map[ID]
	if(is){
		temp.Orders = list
		all_elevs.Elev_map[ID] = temp
	}
}
