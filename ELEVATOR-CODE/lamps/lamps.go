package lamps

//may be moved to fsm or something

import (
	"../orders"
	"../elevio"
	"../types"
	"../constants"
)

func SetAllLamps(e types.Elevator){
	for i:= 0; i < constants.N_FLOORS; i++{
		for j:= 0; j < constants.N_BUTTONS; j++{
			if(orders.IsOrder(e, i, elevio.ButtonType(j))){ //todo change to oneliner
				elevio.SetButtonLamp(elevio.ButtonType(j), i, true)
			} else {
				elevio.SetButtonLamp(elevio.ButtonType(j), i, false)
			}
		}
	}
}
