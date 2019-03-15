package main

import (
	"./elevio"
 	"./globalconstants"
 	"./states"
 	"./fsm"
 	"./orderassigner"
	 "./lamps"
	 "./timer"
)

func main(){
	//initialization
	elevio.Init("localhost:15657", globalconstants.N_FLOORS)
	//order.Init(numFloors, numElevators)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors  := make(chan int)
	order_added := make(chan bool) //for informing FSM about order update when idle
	add_order   := make(chan elevio.ButtonEvent) //send orders from assigner to orders
	door_timeout:= make(chan bool)
	start_door_timer:= make(chan bool)

	//Server channels
	clear_floor := make(chan int) //FSM tells order to clear order
	update_ID := make(chan int) //todo: kan hende vi bare bør droppe denne, vente og se
	update_state := make(chan states.ElevatorState)
	update_floor := make(chan int)
	update_direction := make(chan elevio.MotorDirection)

	go elevio.PollFloorSensor(drv_floors)
	//go elevatorstates.PollAndSetDirection() //bare forslag, tror vi dropper
	fsm.InitElevator(drv_floors)
	
	//run
	go elevio.PollButtons(drv_buttons)
	go fsm.FSM(drv_floors/*, clear_floor*/, order_added, start_door_timer, door_timeout, update_state, update_floor, update_direction/*, chans.....*/)
	go orderassigner.AssignOrder(drv_buttons, add_order)
	go timer.DoorTimer(start_door_timer, door_timeout)
	go lamps.SetLamps()

	//Servers
	//go orders.UpdateOrders(add_order, clear_floor, order_added)
	go fsm.UpdateElevator(update_ID, update_state, update_floor, update_direction, add_order, clear_floor, order_added)
	
	
}