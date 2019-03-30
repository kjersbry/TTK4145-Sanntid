package orderassigner

import (
	"../elevio"
	"../states"
	"../types"
	"../orders"
	"../constants"	
)

func AssignOrder(drvButton <-chan elevio.ButtonEvent, addOrder chan<- types.AssignedOrder, localID string) {
	
	var assignedOrder types.AssignedOrder

	for {
		select {
		case order := <-drvButton:
			if order.Button == elevio.BT_Cab {
				assignedOrder = types.AssignedOrder{localID, order}
				addOrder <- assignedOrder
			} else {
				workingElevs := states.WorkingElevs()
				var selectedElevator = assignAlg(order, workingElevs)
				assignedOrder := types.AssignedOrder{selectedElevator, order}
				addOrder <- assignedOrder
			}
		}
	}
}

func assignAlg(newOrder elevio.ButtonEvent, elevators []types.Elevator) string {

	var bestChooice string
	var bestDuration float64

	var currentDuration float64

	for i, elev := range elevators {
		elev.Orders[newOrder.Floor][newOrder.Button].State = types.OS_AcceptedOrder
		currentDuration = timeToIdle(elev)
		
		if (currentDuration < bestDuration) || i == 0 {
			bestDuration = currentDuration
			bestChooice = elev.Elevator_ID
		}
		
	}

	return bestChooice
}

//Estimates of how much time that will elapse before the elevator get idle.
func timeToIdle(e types.Elevator) float64 {

	var duration float64

	switch e.State { 
	case types.ES_Idle:
		e.Direction = orders.ChooseDirection(e)
		if e.Direction == elevio.MD_Stop {
			return duration
		}

	case types.ES_DoorOpen:
		duration -= constants.AVERAGE_REMAINING_DOOR_OPEN

	case types.ES_Moving:
		duration += constants.AVERAGE_REMAINING_TRAVEL_TIME
		e.Floor = states.UpcomingFloor(e)
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
			duration += constants.DOOR_OPEN_SEC     
		}

		e.Floor = states.UpcomingFloor(e)

		duration += constants.TRAVEL_TIME
	}
}
