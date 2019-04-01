package main

import (
	"flag"
	"fmt"
	"time"

	"./bcast"
	"./constants"
	"./elevio"
	"./fsm"
	"./peers"
	"./states"
	"./timer"
	"./types"
	"./operation"
	"./localip"
	"os/exec"
)

/*Simulator port suggestions:
15657, 59334, 46342, 33922, 50945, 36732*/

func main() {
	var spawn_sim string
	var server_port string

	flag.StringVar(&server_port, "sport", "15657", "port for the elevator server")
	flag.StringVar(&spawn_sim, "sim", "no", "set -sim=yes if you want to spawn simulator")
	flag.Parse()

	if spawn_sim == "yes" {
		newProcess := exec.Command("gnome-terminal", "-x", "sh", "-c", "./SimElevatorServer --port " + server_port)
		err := newProcess.Run()
		if err != nil {
			fmt.Printf("\nCould not spawn simulator\n")
			return
		}
		time.Sleep(time.Second * 2)
	}

	//initialization
	local_ID := localip.GetPeerID()
	elevio.Init("localhost:" + server_port, constants.N_FLOORS)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	order_added := make(chan bool)
	door_timeout := make(chan bool)
	start_door_timer := make(chan bool)
	floor_reached := make(chan bool)

	//Server channels
	clear_floor := make(chan int)
	update_state := make(chan types.ElevatorState)
	update_floor := make(chan int)
	update_direction := make(chan elevio.MotorDirection)

	elev_rx := make(chan types.Wrapped_Elevator)
	elev_tx := make(chan types.Wrapped_Elevator)

	go elevio.PollFloorSensor(drv_floors)
	states.InitElevators(local_ID, drv_floors)

	//Connections
	operation_update := make(chan types.Operation_Event)  
	connection_update := make(chan types.Connection_Event) 
	go peers.ConnectionObserver(33924, connection_update, local_ID)
	go peers.ConnectionTransmitter(33924, local_ID)
	go operation.OperationObserver(operation_update, local_ID)
	go bcast.Transmitter(33922, elev_tx)
	go bcast.Receiver(33922, elev_rx)

	//run
	go elevio.PollButtons(drv_buttons)
	go states.UpdateElevator(update_state, drv_floors, update_direction, clear_floor, floor_reached, order_added, drv_buttons, elev_rx, elev_tx, connection_update, operation_update)
	go fsm.FSM(floor_reached, clear_floor, order_added, start_door_timer, door_timeout, update_state, update_floor, update_direction)
	go timer.DoorTimer(start_door_timer, door_timeout)

	/*Infinite loop: */
	fin := make(chan int)
	for {
		select {
		case <-fin:
		}
	}
}