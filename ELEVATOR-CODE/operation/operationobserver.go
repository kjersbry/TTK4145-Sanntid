package operation

import (
	"time"
	"types"
)

//Note: 1. The datatypes and the function probably don't need their own file. But, as of now, i don't see any better place to put them
//2. Nothing in this file has been properly tested.

//finding unoperational elevator

type floor_time struct {	//Bad name. 
	LastFloor int
	timeOfChange time.Time
	is_Operational bool			//Her kan error state puttes -> Enumeration, som Kjersti foreslo
}

func OperationObserver (elev_map map[string]types.Elevator, Operation_Update chan<- types.Operation_Event) {
	var update types.Operation_Event
	lastChange := make(map[string]floor_time)
	var localID = states.ReadLocalID() //Don't know the exact function name --> case senesitve
	
	for {			
		for k, v := range elev_map {
			if k = localID {
			continue
			}
			if _, keyExists := lastChange[k]; !keyExists {
				if k.is_Operational = false {			//Snakk med Kjersti om denne antagelsen
					lastChange[k] = floor_time{v.Floor, time.Now(), false}
				}else {
					lastChange[k] = floor_time{v.Floor, time.Now(), true}		
				}
			} else {
				if lastChange[k].LastFloor !=  v.Floor {
					if lastChange[k].is_Operational == false {
						lastChange[k].is_Operational == true
						update = {k, true}
						Operation_Update <- update
					}
					lastChange[k].LastFloor =  v.Floor
					lastChange[k].timeOfChange =  time.Now()
				} else if time.Now().Sub(lastChange[k].timeOfChange) > 15 { //Put in proper constant name
					lastChange[k].is_Operational = false					//The exact value is up for debate
					update = {k, false}
					Operation_Update = update
				}
				

			}		
		}
		time.Sleep(time.Second*5)  //Is this done correctly?
	}
}
