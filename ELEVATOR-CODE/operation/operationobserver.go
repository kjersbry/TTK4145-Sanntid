package operation

import (
	"time"
	"../types"
	"../constants"
	"../elevio"
	"../orders"	
)

type operationInfo struct {
	lastFloor      int
	timeOfChange   time.Time
	isOperational bool
}

func OperationObserver(operationUpdate chan<- types.OperationEvent, localID string,  allElevsUpdate <-chan map[string]types.Elevator) {
	elevMap := make(map[string]types.Elevator)
	lastChange := make(map[string]operationInfo)
	var update types.OperationEvent
	var temp operationInfo

	tick := time.NewTicker(time.Second*5)

	for {
		select{
		case elevMap = <-allElevsUpdate:
		case <-tick.C:
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
		}
	}
}


