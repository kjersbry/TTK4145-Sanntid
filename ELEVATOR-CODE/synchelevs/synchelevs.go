package synchstates

import (
	"../fsm"
	"../types"
	"../constants"
	"fmt"
	"time"
)

/*
func Dosomething() {
	var thiselev types.Elevator
	thiselev = fsm.ReadElevator()
	fmt.Printf("grr %d", e.Floor)
}*/



func TransmitElev(state_tx chan<- types.Elevator){
	for {
		state_tx <- fsm.ReadElevator()
		time.Sleep(time.Milliseconds*constants.TRANSMIT_MS)
	}
}

func ReceiveElevs(state_rx <-chan types.Elevator){
	
}
