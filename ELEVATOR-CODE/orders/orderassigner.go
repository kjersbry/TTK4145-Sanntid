package orders

import (
	"../elevio"
	"../states"
	"../types"
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
				elevs := states.ReadAllElevators()
				var selectedElevator = assignAlgorithm(order, elevs)
				assignedOrder := types.AssignedOrder{selectedElevator, order}
				addOrder <- assignedOrder
			}
		}
	}
}

func OrderReassigner(faultyElevID string, localID string, allElevators map[string]types.Elevator) (types.Elevator, bool) {
	e := allElevators[faultyElevID]
	localElev := allElevators[localID]
	isNewLocalOrder := false
	for i := 0; i < constants.N_FLOORS; i++ {
		for j := 0; j < constants.N_BUTTONS-1; j++ {
			if e.Orders[i][j].State == types.OS_AcceptedOrder {
				localElev[i][j] = types.OS_AcceptedOrder
				isNewLocalOrder = true
			}
		}
	}
	return localElev, isNewLocalOrder
}

func assignAlgorithm(newOrder elevio.ButtonEvent, elevators map[string]types.Elevator) string {
	var bestChoice string
	var bestDuration float64
	var currentDuration float64

	i := 0
	for _, elev := range elevators {
		//if operational and connected
		elev.Orders[newOrder.Floor][newOrder.Button].State = types.OS_AcceptedOrder
		currentDuration = timeToIdle(elev)
		
		if (currentDuration < bestDuration) || i == 0 {
			bestDuration = currentDuration
			bestChoice = elev.Elevator_ID
		}
		i = 1
	}

	return bestChoice
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
		e.Floor = upcomingFloor(e)
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

		e.Floor = upcomingFloor(e)

		duration += constants.TRAVEL_TIME
	}
}


func upcomingFloor(e types.Elevator) int {
	if e.Direction == elevio.MD_Up {
		return e.Floor + 1
	} else if e.Direction == elevio.MD_Down {
		return e.Floor - 1
	} else {
		return e.Floor
	}
}