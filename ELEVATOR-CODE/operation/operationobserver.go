//Heisen må kunne markere seg selv som notoperational!!!


package operation

import (
	"time"
	"../types"
	"../states"
	"../constants"
	"../elevio"
	"../orders"
	"fmt"
)

//Note: 1. The datatypes and the function probably don't need their own file. But, as of now, i don't see any better place to put them
//2. Nothing in this file has been properly tested.


//finding unoperational elevator

type floor_time struct { //Bad name.
	LastFloor      int
	timeOfChange   time.Time
	is_Operational bool //Her kan error state puttes -> Enumeration, som Kjersti foreslo
}

func OperationObserver(Operation_Update chan<- types.Operation_Event, localID string) {  //Operation observer running once for no reason?
	var elev_map map[string]types.Elevator 
	var update types.Operation_Event
	lastChange := make(map[string]floor_time)   operations
	var temp floor_time

	for {
		elev_map = states.ReadAllElevators()
		for k, v := range elev_map {							//Markerer først seg selv som operational
			if _, keyExists := lastChange[k]; !keyExists {
				lastChange[k] = floor_time{v.Floor, time.Now(), true}				
			}
			/*if v.Is_Operational == false { //Snakk med Kjersti om denne antagelsen  //Alternative??
						lastChange[k] = floor_time{v.Floor, time.Now(), false}
						fmt.Printf("\nThe elevator was marked as unoperational\n")
					} else {
						lastChange[k] = floor_time{v.Floor, time.Now(), true}
						fmt.Printf("\nThe elevator was marked as operational\n")
			}*/
																		//For rask med å markere?  En feil som dukket opp en gang er at den markerte som not operational for fort
			if orders.ChooseDirection(v) != elevio.MD_Stop {			//Printer ut Marked as not oper... gjenntatte ganger. Noenganger 2 på rappen
				if k == localID {
					continue
				} else {
					if lastChange[k].LastFloor != v.Floor {
						if lastChange[k].is_Operational == false {
							update.Elevator_ID = k
							update.Is_Operational = true
							Operation_Update <- update
							fmt.Printf("\nSet to operational at: %v\n", time.Now())
						}
						temp = floor_time{v.Floor, time.Now(), true}
						lastChange[k] = temp
						//Exact time is up for debate.     //Important error -> Markes the other elevator as not oper to fast.
					} else if time.Now().Sub(lastChange[k].timeOfChange) > constants.ELEVATOR_TIMEOUT { //Put in proper constant name
						if lastChange[k].is_Operational == true {
							temp = floor_time{lastChange[k].LastFloor, lastChange[k].timeOfChange, false}
							lastChange[k] = temp //Should this be here? This was added to accomodate for a previous error
							update.Elevator_ID = k
							update.Is_Operational = false
							Operation_Update <- update
							//a := time.Now().Sub(lastChange[k].timeOfChange)
							//fmt.Printf("\nSet to not operational at: %v\n Delta T = %v\n", time.Now(), time.Now().Sub(lastChange[k].timeOfChange) - constants.ELEVATOR_TIMEOUT)
							//fmt.Printf("\nSet to not operational at: %v\n time since last change = %v\n delta %v\n", time.Now(), a, time.Now().Sub(lastChange[k].timeOfChange).Sub(constants.ELEVATOR_TIMEOUT))
						} 
					}

				}
			} else {
				temp = floor_time{lastChange[k].LastFloor, time.Now(), true}
				lastChange[k] = temp
			}
		}
		time.Sleep(time.Second * 5) //Is this done correctly?
	}
}


