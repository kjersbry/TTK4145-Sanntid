package lamps

//handles all lighting and unlighting of order buttons
//may be moved to orders or something


import (
	"../globalconstants"
	"../orders"
	"../elevio"
	"../fsm"
)

func SetLamps(){
	for {
		for i:= 0; i < globalconstants.N_FLOORS; i++{
			for j:= 0; j < globalconstants.N_BUTTONS; j++{
				if orders.IsOrder(fsm.ReadElevator(), i, elevio.ButtonType(j)){ 
					//TODO: WHEN ADDING MORE ELEVATORS, CHANGE THE ABOVE CONDITION TO "if isOrder() AND orderstate = accepted"
					elevio.SetButtonLamp(elevio.ButtonType(j), i, true)
				} else {
					elevio.SetButtonLamp(elevio.ButtonType(j), i, false)
				}
			}
		}
	}
}