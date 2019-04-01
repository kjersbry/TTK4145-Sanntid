//Heisen må kunne markere seg selv som notoperational!!!


package operation

import (
	"time"
	"../types"
	"../states"
	"../constants"
	"../elevio"
	"../orders"	
)


type operationInfo struct {
	lastFloor      int
	timeOfChange   time.Time
	isOperational bool
}

func OperationObserver(operationUpdate chan<- types.OperationEvent, localID string) {
	var elevMap map[string]types.Elevator 
	var update types.OperationEvent
	lastChange := make(map[string]operationInfo)
	var temp operationInfo

	for {
		elevMap = states.ReadAllElevators()
		for k, v := range elevMap {
			if _, keyExists := lastChange[k]; !keyExists {
				lastChange[k] = operationInfo{v.Floor, time.Now(), true}
			}
			
			if orders.ChooseDirection(v) != elevio.MD_Stop {
				if k == localID {
					continue
				} else {
					if lastChange[k].lastFloor != v.Floor {
						if lastChange[k].isOperational == false {
							update.ElevatorID = k
							update.IsOperational = true
							operationUpdate <- update
						}
						temp = operationInfo{v.Floor, time.Now(), true} 
						lastChange[k] = temp	
					} else if time.Now().Sub(lastChange[k].timeOfChange) > constants.ELEVATOR_TIMEOUT { 
						if lastChange[k].isOperational == true { 
							temp = operationInfo{lastChange[k].lastFloor, lastChange[k].timeOfChange, false}
							lastChange[k] = temp 
							update.ElevatorID = k
							update.IsOperational = false
							operationUpdate <- update				
						} 
					}

				}
			} else {
				temp = operationInfo{lastChange[k].lastFloor, time.Now(), true}  
				lastChange[k] = temp
			}
		}
		time.Sleep(time.Second * 5)
	}
}


