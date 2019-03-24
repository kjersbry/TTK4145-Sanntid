package main

import (
	"./elevio"
	"./states"
 	"./types"
 	"./fsm"
 	"./orderassigner"
	"./timer"
	"./constants"
	"./bcast"
	"operation" //This will probably be changed
)

func main(){
	//initialization
	//15657, 59334, 46342, 33922, 50945, 36732
	elevio.Init("localhost:46342", constants.N_FLOORS)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors  := make(chan int)
	order_added := make(chan int) //for informing FSM about order update
	add_order   := make(chan types.AssignedOrder) //send orders from assigner to orders
	door_timeout:= make(chan bool)
	start_door_timer:= make(chan bool)
	floor_reached := make(chan bool)

	//Server channels
	clear_floor := make(chan int) //FSM tells order to clear order
	//update_ID := make(chan string) //todo: kan hende vi bare bør droppe denne, vente og se
	update_state := make(chan types.ElevatorState)
	update_floor := make(chan int)
	update_direction := make(chan elevio.MotorDirection)

	elev_rx := make(chan types.Wrapped_Elevator)
	elev_tx := make(chan types.Wrapped_Elevator)
	sendwrap_request := make(chan string)

	go elevio.PollFloorSensor(drv_floors)
	states.InitSynchElevators(drv_floors)

	//run
	go elevio.PollButtons(drv_buttons)
	go states.UpdateElevator(update_state, drv_floors, update_direction, clear_floor, floor_reached, order_added, add_order, elev_rx, sendwrap_request, elev_tx)
	go fsm.FSM(floor_reached, clear_floor, order_added, start_door_timer, door_timeout, update_state, update_floor, update_direction/*, chans.....*/)
	go orderassigner.AssignOrder(drv_buttons, add_order)
	go timer.DoorTimer(start_door_timer, door_timeout)

	go states.TransmitElev(elev_tx)
	go bcast.Transmitter(33922, elev_tx)
	go bcast.Receiver(33922, elev_rx)

	go states.TestPrintAllElevators()

	operation_Update := make(chan Operation_Event) //Update elevator must use this <- Remember to update
	connection_Update := make(chan Connection_Event) //Update elevator must use this <- Remember to update
	//Connections
	go peers.ConnectionObserver(33922, update_connection)
	go operations.OperationObserver(ReadAllElevs, operation_Update)


	/*Infinite loop: */
	fin := make(chan int)
	for{
		select{
		case <- fin:
		}
	}
}
