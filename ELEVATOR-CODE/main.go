package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"./bcast" //This will probably be changed
	"./constants"
	"./elevio"
	"./fsm"
	"./orderassigner"
	"./peers"
	"./phoenix"
	"./states"
	"./timer"
	"./types"
)

/* TODO: sjekk om det er flere defaults på for{select{}} (men ikke fjern der det er select uten for)*/

const default_sport int = 15657
const default_pport int = 1234

var sport_suggestions = [6]int{default_sport, 59334, 46342, 33922, 50945, 36732}
var pport_suggestions = [6]int{default_pport, 1235, 1236, 1237, 1238, 1239}
var spawn_sim string
var server_port int
var phoenix_port int

func main() {
	flag.IntVar(&server_port, "sport", default_sport, "port for the elevator server")
	flag.IntVar(&phoenix_port, "pport", default_pport, "port for phoenix")
	flag.StringVar(&spawn_sim, "sim", "no", "set -sim=yes if you want to spawn simulator")
	flag.Parse()

	//runtime.GOMAXPROCS(runtime.NumCPU())

	//assume that this is the backup process
	phoenix.RunBackup(phoenix_port, server_port, runElevator)
}

func runElevator(local_ID string, server_port string) {
	if spawn_sim == "yes" {
		err := phoenix.SpawnTerminal("./SimElevatorServer --port " + server_port)
		if err != nil {
			fmt.Printf("\nCould not spawn simulator\n")
			return
		}
		time.Sleep(time.Second * 2)
	}

	//initialization
	elevio.Init("localhost:"+server_port, constants.N_FLOORS)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	order_added := make(chan bool)              //for informing FSM about order update
	add_order := make(chan types.AssignedOrder) //send orders from assigner to orders
	door_timeout := make(chan bool)
	start_door_timer := make(chan bool)
	floor_reached := make(chan bool)

	//Server channels
	clear_floor := make(chan int) //FSM tells order to clear order
	update_state := make(chan types.ElevatorState)
	update_floor := make(chan int)
	update_direction := make(chan elevio.MotorDirection)

	elev_rx := make(chan types.Wrapped_Elevator)
	elev_tx := make(chan types.Wrapped_Elevator)

	sendwrap_request := make(chan string)

	go elevio.PollFloorSensor(drv_floors)
	states.InitElevators(local_ID, drv_floors)

	//Connections
	operation_update := make(chan types.Operation_Event)   //Update elevator must use this <- Remember to update
	connection_update := make(chan types.Connection_Event) //Update elevator must use this <- Remember to update
	go peers.ConnectionObserver(33922, connection_update)
	go peers.ConnectionTransmitter(33922, local_ID)
	go operations.OperationObserver(ReadAllElevs, operation_update, local_ID)
	go bcast.Transmitter(33922, elev_tx)
	go bcast.Receiver(33922, elev_rx)

	//run
	go elevio.PollButtons(drv_buttons)
	go states.UpdateElevator(update_state, drv_floors, update_direction, clear_floor, floor_reached, order_added, add_order, elev_rx, connection_update, operation_update)
	go fsm.FSM(floor_reached, clear_floor, order_added, start_door_timer, door_timeout, update_state, update_floor, update_direction /*, chans.....*/)
	go orderassigner.AssignOrder(drv_buttons, add_order)
	go timer.DoorTimer(start_door_timer, door_timeout)
	go states.TransmitElev(elev_tx)

	go states.TestPrintAllElevators()

	/*Infinite loop: */
	fin := make(chan int)
	for {
		select {
		case <-fin:
		}
	}
}
