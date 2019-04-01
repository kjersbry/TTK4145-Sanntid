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
	var spawnSim string
	var serverPort string

	flag.StringVar(&serverPort, "sport", "15657", "port for the elevator server")
	flag.StringVar(&spawnSim, "sim", "no", "set -sim=yes if you want to spawn simulator")
	flag.Parse()

	if spawnSim == "yes" {
		newProcess := exec.Command("gnome-terminal", "-x", "sh", "-c", "./SimElevatorServer --port " + serverPort)
		err := newProcess.Run()
		if err != nil {
			fmt.Printf("\nCould not spawn simulator\n")
			return
		}
		time.Sleep(time.Second * 2)
	}

	//initialization
	localID := localip.GetPeerID()
	elevio.Init("localhost:" + serverPort, constants.N_FLOORS)
	drvButtons 		:= make(chan elevio.ButtonEvent)
	drvFloors 		:= make(chan int)
	orderAdded 		:= make(chan bool)
	doorTimeout 	:= make(chan bool)
	startDoorTimer 	:= make(chan bool)
	floorReached 	:= make(chan bool)

	//Server channels
	clearFloor 		:= make(chan int)
	updateState 	:= make(chan types.ElevatorState)
	updateDirection := make(chan elevio.MotorDirection)

	elevRX 			:= make(chan types.WrappedElevator)
	elevTX 			:= make(chan types.WrappedElevator)

	go elevio.PollFloorSensor(drvFloors) 
	states.InitElevators(localID, drvFloors) 

	//Connections
	operationUpdate := make(chan types.OperationEvent)  
	connUpdate 		:= make(chan types.ConnectionEvent) 
	go peers.ConnectionObserver(33924, connUpdate, localID)
	go peers.ConnectionTransmitter(33924, localID) 
	go operation.OperationObserver(operationUpdate, localID)
	go bcast.Transmitter(33922, elevTX)
	go bcast.Receiver(33922, elevRX)

	//run
	go elevio.PollButtons(drvButtons)
	go states.UpdateElevator(updateState, drvFloors, updateDirection, clearFloor, floorReached, orderAdded, drvButtons, elevRX, elevTX, connUpdate, operationUpdate)
	go fsm.FSM(floorReached, clearFloor, orderAdded, startDoorTimer, doorTimeout, updateState, updateDirection)
	go timer.DoorTimer(startDoorTimer, doorTimeout)

	/*Infinite loop: */
	fin := make(chan int)
	for {
		select {
		case <-fin:
		}
	}
}