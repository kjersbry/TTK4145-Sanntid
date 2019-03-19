package main

import (
	"./elevio"
 	"./globalconstants"
 	"./states"
 	"./fsm"
 	"./orderassigner"
	 "./timer"
	 "time"
)

func main(){
	//initialization
	elevio.Init("localhost:15657", globalconstants.N_FLOORS)
	//order.Init(numFloors, numElevators)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors  := make(chan int)
	order_added := make(chan int) //for informing FSM about order update
	add_order   := make(chan elevio.ButtonEvent) //send orders from assigner to orders
	door_timeout:= make(chan bool)
	start_door_timer:= make(chan bool)
	floor_reached := make(chan bool)

	//Server channels
	clear_floor := make(chan int) //FSM tells order to clear order
	update_ID := make(chan int) //todo: kan hende vi bare bør droppe denne, vente og se
	update_state := make(chan states.ElevatorState)
	update_floor := make(chan int)
	update_direction := make(chan elevio.MotorDirection)

	go elevio.PollFloorSensor(drv_floors)
	fsm.InitElevator(drv_floors)

	//run
	go elevio.PollButtons(drv_buttons)
	go fsm.UpdateElevator(update_ID, update_state, drv_floors, update_direction, floor_reached, add_order, clear_floor, order_added)
	go fsm.FSM(floor_reached, clear_floor, order_added, start_door_timer, door_timeout, update_state, update_floor, update_direction/*, chans.....*/)
	go orderassigner.AssignOrder(drv_buttons, add_order)
	go timer.DoorTimer(start_door_timer, door_timeout)

	//Servers
	//go orders.UpdateOrders(add_order, clear_floor, order_added)


	for{
		time.Sleep(time.Second*100)
	}
}
