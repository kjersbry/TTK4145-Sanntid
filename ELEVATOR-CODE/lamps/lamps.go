package lamps

//may be moved to fsm or something

import (
	"../orders"
	"../elevio"
	"../types"
	"../constants"
	//"fmt"
)

func SetAllLamps(local_ID string, elevators map[string]types.Elevator){
	for i:= 0; i < constants.N_FLOORS; i++{
		for j:= 0; j < constants.N_BUTTONS - 1 ; j++ {
			is_some_hallcall := false
			for key, val := range elevators {
				if orders.IsOrder(val, i, elevio.ButtonType(j)) {
					//fmt.Printf("\nSet light fl. %d, b %d\n", i, j)
					is_some_hallcall = true
				} 
				if key == local_ID {
				elevio.SetButtonLamp(elevio.BT_Cab, i, orders.IsOrder(elevators[local_ID], i, elevio.BT_Cab))
				}
			}
			elevio.SetButtonLamp(elevio.ButtonType(j), i, is_some_hallcall)
		}
	}
}
