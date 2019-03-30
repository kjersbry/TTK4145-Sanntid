package lamps

import (
	"../orders"
	"../elevio"
	"../types"
	"../constants"
)

func SetAllLamps(localID string, elevators map[string]types.Elevator){
	for i:= 0; i < constants.N_FLOORS; i++{
		for j:= 0; j < constants.N_BUTTONS - 1 ; j++ {
			isSomeHallCall := false
			for key, val := range elevators {
				if orders.IsOrder(val, i, elevio.ButtonType(j)) && val.Is_Operational {
					isSomeHallCall = true
				} 
				if key == localID {
				elevio.SetButtonLamp(elevio.BT_Cab, i, orders.IsOrder(elevators[localID], i, elevio.BT_Cab))
				}
			}
			elevio.SetButtonLamp(elevio.ButtonType(j), i, isSomeHallCall)
		}
	}
}
