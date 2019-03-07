package main

import "elevio"
import "globalconstants"

func main(){
	//initialization
	elevio.Init("localhost:15657", N_FLOORS)
	//order.Init(numFloors, numElevators)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors  := make(chan int)
	order_added := make(chan bool) //for informing FSM about order update when idle
	add_order   := make(chan elevio.ButtonEvent) //send orders from assigner to orders
	clear_floor := make(chan int) //FSM tells order to clear order

	go elevio.PollFloorSensor(drv_floors)
	//go elevatorstates.PollAndSetDirection() //bare forslag
	elevatorstates.InitElevator(drv_floors)
	
	//run
	go elevio.PollButtons(drv_buttons)
	go fsm.FSM(drv_floors, clear_floor, order_added /*, chans.....*/)
	go orderassigner.AssignOrder(drv_buttons, add_order)
	go orders.OrdersServer(add_order, clear_floor, order_added)
	
}