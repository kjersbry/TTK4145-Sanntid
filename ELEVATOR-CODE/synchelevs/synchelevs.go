package synchelevs

import (
	"../fsm"
	"../types"
	//"../constants"
	"fmt"
	//"time"
)


func Dosomething() {
	var thiselev types.Elevator
	thiselev = fsm.ReadElevator()
	fmt.Printf("grr %d", thiselev.Floor)
}
