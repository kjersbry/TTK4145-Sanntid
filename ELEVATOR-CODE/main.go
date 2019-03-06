package main

import "./elevio"


func main(){
	//initialization
	numFloors := 4
	elevio.Init("localhost:15657", numFloors)
	//order.Init(numFloors, numElevators)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors  := make(chan int)
	order_added := make(chan bool) //for informing FSM about order update when idle

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	//go elevatorstates.PollAndSetDirection() //bare forslag
	elevatorstates.InitElevator(drv_floors)
	
	//run
	go fsm.FSM(drv_floors, order_added /*, chans.....*/)
	go orderassigner.AssignOrder(drv_buttons/*, requestwriteorder <-chan*/)
	go orders.OrderServer(/*write request chan<-, */ order_added)
	
}