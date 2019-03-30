package merger

import (
	 "../constants"
	 "../types"
	 "../elevio"
	 "math"
)

func MergeOrders(localID string, localMap map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order,
	 otherMap map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order) (map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order, bool, bool) {
		
		isNewLocalOrder := false
		shouldLight := false

		for key, val := range localMap {
			temp, keyExists := otherMap[key]
			tempLocal := val

			if keyExists {
			for i := range val {
				for j := range val[i]{
					localState := tempLocal[i][j].State
					localStateWas := localState
					otherState := temp[i][j].State

					localCounter := tempLocal[i][j].Counter
					otherCounter := temp[i][j].Counter

					if j == elevio.BT_Cab {
						if localCounter == otherCounter {
							if localState != types.OS_NoOrder || otherState != types.OS_NoOrder {
								localState = types.OS_UnacceptedOrder
							}
						} else if otherCounter > localCounter {
							localState = otherState
						} 

					} else if localCounter == otherCounter { 
						switch {
						case (localState == types.OS_UnacceptedOrder && otherState == types.OS_UnacceptedOrder):
							localState = types.OS_AcceptedOrder
						case math.Abs(float64(localState - otherState)) == 2:
							localState = types.OS_AcceptedOrder
						default :
							localState = MaxState(localState, otherState)
						}	
					} else {
						presedence := getPresedence(tempLocal[i][j], temp[i][j])
						localState = presedence.State
					}
					tempLocal[i][j].State = localState

					presedence := getPresedence(tempLocal[i][j], temp[i][j])
					tempLocal[i][j].Counter = presedence.Counter

					//Check if FSM should get local-orders-updated-signal
					if localState != localStateWas {
						shouldLight = true
						if key == localID {
							if (localState == types.OS_AcceptedOrder) || ((j == constants.N_BUTTONS - 1) && localState != types.OS_NoOrder) {
								isNewLocalOrder = true
							}
						}
					}
					}
				}
				localMap[key] = tempLocal
			}
		}
		
	return localMap, isNewLocalOrder, shouldLight
}

func getPresedence(order1 types.Order, order2 types.Order) types.Order {
	if order2.Counter > order1.Counter {
		return order2
	} else {
		return order1
	}
}

func largestID(elev1 string, elev2 string) string {
	if elev1 > elev2 {
		return elev1
	}
	return elev2
}
func MaxState(state1 types.OrderState, state2 types.OrderState) types.OrderState{
	if state2 > state1 {
		return state2
	}
	return state1
}