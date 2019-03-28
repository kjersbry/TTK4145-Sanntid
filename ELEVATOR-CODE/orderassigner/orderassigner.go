package orderassigner

import (
	"../elevio"
	"../states"
	"../types"
	"../orders"
	"../constants"
//	"fmt"
	
)


//Assigns only to self. Good to have for testing purposes
/*func AssignOrder(drv_button <-chan elevio.ButtonEvent, add_order chan<- types.AssignedOrder, local_ID string) {
	for {
		select {
		case order := <-drv_button:
			//her legges assignment algorithm
			//en heis---> kun dette:
			dummyassigned := types.AssignedOrder{states.ReadLocalElevator().Elevator_ID, order}
			add_order <- dummyassigned //skriver resultat til order
		}
	}
}*/

//use this one
func AssignOrder(drv_button <-chan elevio.ButtonEvent, add_order chan<- types.AssignedOrder, local_ID string) {
	var previous_order elevio.ButtonEvent   //Todo -> Add a time restriction
	var previous_assigned_order types.AssignedOrder
	var assigned_order types.AssignedOrder
	
	for {
		select {
		case order := <-drv_button:
			if previous_order != order {
				if order.Button == elevio.BT_Cab {
					assigned_order = types.AssignedOrder{local_ID, order}
					add_order <- assigned_order //skriver resultat til order
				} else {
				workingElevs := states.WorkingElevs()
				var selected_elevator = assignAlg(order, workingElevs)
				assigned_order := types.AssignedOrder{selected_elevator, order}
				add_order <- assigned_order //skriver resultat til order
				}
			} else {
				add_order <- previous_assigned_order 
			}
			previous_order = order
			previous_assigned_order = assigned_order
		}
	}
}

func assignAlg(new_order elevio.ButtonEvent, elevators []types.Elevator) string {

	var best_chooice string
	var best_duration float64

	var currentDuration float64

	for i, elev := range elevators {
		if elev.Is_Operational {
			elev.Orders[new_order.Floor][new_order.Button].State = types.OS_AcceptedOrder
			currentDuration = timeToIdle(elev)

			if (currentDuration < best_duration) || i == 0 {
				best_duration = currentDuration
				best_chooice = elev.Elevator_ID
			}
		}
	}

	return best_chooice
}

//TimeToIdle gives an estimate of how much that time that will elapse before an elevator as handled all its requests.
func timeToIdle(e types.Elevator) float64 {

	var duration float64
	duration = 0

	switch e.State { 
	case types.ES_Idle:
		e.Direction = orders.ChooseDirection(e)
		if e.Direction == elevio.MD_Stop {
			return duration
		}
		break

	case types.ES_DoorOpen:
		duration -= 1.5 //1. Find proper constant name (or use 3. sec) 2. Potential datatype problems?
		break

	case types.ES_Moving:
		duration += 2.5 //Find a proper name for the constant
		e.Floor = states.UpcomingFloor(e)
		break
	}

	for {
		if orders.ShouldStop(e) {
			for button := 0; button < constants.N_BUTTONS; button++ {
				if e.Floor < 0 || e.Floor >= constants.N_FLOORS {

				}
				e.Orders[e.Floor][button].State = types.OS_NoOrder
			}
			
			e.Direction = orders.ChooseDirection(e)
			if e.Direction == elevio.MD_Stop {
				return duration
			}
			duration += 3     
		}

		e.Floor = states.UpcomingFloor(e)

		duration += 4 //TravelTime //Insert the proper operator
	}
}
