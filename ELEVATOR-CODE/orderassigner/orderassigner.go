package orderassigner

import (
	"../elevio"
	"../types"
	"../states"
)

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

//This function is responsible for reciving unassigned orders from elevator_io.
//The order is evaluated by assignAlg and assigned to an elevator.
//The the order and the ID of the assigned elevator is sent to a channel.
/*func AssignOrder(drv_button <-chan elevio.ButtonEvent, add_order chan<- types.AssignedOrder){
	for{
		select{
		case order:= <- drv_button:
			selected_elevator := assignAlg(order)
			assgined_order := AssignedOrder{order, selected_elevator}
			add_order <- assigned_order
		}
	}	
}*/

//This function adds the new order to the que of every avaible elevator. (Only done locally)
//the function timeToIdle is run once per elevator. 
//It returns the ID of the elevator with the lowest return value from timeToIdle
/*func assignAlg(new_order) int {
	var best_choice int 
	var best_duration int

	var currentElevator elevatorstates.Elevator
	var currentDuration int

	/*Note the following alternative: As of now, the readElevator function is called once for every elevator. We may be better off
	by making one function call, e.g. readAllElevators. And then go through the returned array. Such a function would be usefull for 
	the network mudule as well*/
/*
	for _, elevator := range states.WorkingElevators() {               //WorkingElevators must be created

		//add new_order to current elevator   <-The appropriate function should be made and saved in orders.
												//The sturcture of the order array should be clearly defined before this line is added.


		currentElevator = elevatorstates.ReadElevator(/*elevator <- Specifying witch elevator that will be examined*//*)
		currentDuration = timeToIdle(currentElevator)

		if(currentDuration < best_duration){
			best_chooice = currentElevator.Elevator_ID 
		} 
	}

	return best_choice
}*/
/*

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
