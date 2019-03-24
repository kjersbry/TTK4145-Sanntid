package orderassigner

import (
	"../elevio"
	"../types"
	"../states"
)

//Kjersti's
func AssignOrder(drv_button <-chan elevio.ButtonEvent, add_order chan<- types.AssignedOrder){
	for{
		select{
		case order:= <- drv_button:
			//her legges assignment algorithm
			//en heis---> kun dette:
			dummyassigned := types.AssignedOrder{states.ReadLocalElevator().Elevator_ID, order}
			add_order <- dummyassigned //skriver resultat til order
		}
	}	
}


//Ole's
func AssignOrder(drv_button <-chan ButtonEvent, add_order chan<- AssignedOrder) {
	for {
		select {
		case order := <-drv_button:
			workingElevs = states.WorkingElevs 
			//WorkingElevs returns a slice with all elevators that are connected and operational. Must be made.
			var selected_elevator = assignAlg(order, workingElevs)
			assigned_order := AssignedOrder{order, selected_elevator}
			add_order <- assigned_order //skriver resultat til order
		}
	}
}
func assignAlg(new_order ButtonEvent, elevators []Elevator) int {
	var best_chooice int
	var best_duration float64

	var currentDuration float64

	for _, elev := range elevators {
		if elev.Operational {
			elev.Orders[new_order.Floor][new_order.Button].State = OS_AcceptedOrder
			currentDuration = timeToIdle(elev)

			if (currentDuration < best_duration) || best_duration == 0 {
				best_duration = currentDuration
				best_chooice = elev.Elevator_ID
			}
		}
	}

	return best_chooice
}

//TimeToIdle gives an estimate of how much that time that will elapse before an elevator as handled all its requests.
func timeToIdle(e Elevator) float64 {

	var duration float64
	duration = 0

	switch e.State { //The switch has been tested.
	case ES_Idle:
		e.Direction = ChooseDirection(e) //Put back orders.
		if e.Direction == MD_Stop {
			return duration
		}
		break

	case ES_DoorOpen:
		duration -= 3 / 2 //1. Find proper constant name (or use 3. sec) 2. Potential datatype problems?
		break

	case ES_Moving:
		duration += 5 / 2 //Find a proper name for the constant
		e.Floor = UpcommigFloor(e)
		break
	}

	//For loop and nested functionality remains untested
	for {
		if ShouldStop(e) {
			for button := 0; button < 3; button++ {
				e.Orders[e.Floor][button].State = OS_NoOrder
			}
			duration += 3                    //Put in proper constant name
			e.Direction = ChooseDirection(e) //put back orders.
			if e.Direction == MD_Stop {
				return duration
			}
		}

		e.Floor = UpcommingFloor(e)

		duration += 5 //TravelTime //Insert the proper operator
	}
}
