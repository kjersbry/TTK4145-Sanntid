package lamps

//may be moved to fsm or something

import (
	"../orders"
	"../elevio"
	"../types"
	"../constants"
)

func SetAllLamps(local_ID string, elevators map[string]types.Elevator){
	for key, val := range elevators {
		for i:= 0; i < constants.N_FLOORS; i++{
			for j:= 0; j < constants.N_BUTTONS - 1 ; j++ {
				elevio.SetButtonLamp(elevio.ButtonType(j), i, orders.IsOrder(val, i, elevio.ButtonType(j)))
			}
			if key == local_ID {
				elevio.SetButtonLamp(elevio.BT_Cab, i, orders.IsOrder(elevators[local_ID], i, elevio.BT_Cab))
			}
		}
	}
}
