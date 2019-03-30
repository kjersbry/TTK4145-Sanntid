package merger

import (
	 "../constants"
	 "../types"
	 //"reflect"
	 //"fmt"
)



func MergeOrders(localID string, localElev map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order, elev2 map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order) (map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order, bool) { 
	orderMmap, isNewLocalOrder := combineMaps(localID, localElev, elev2)
	orderMap = removeDuplicates(orderMap)
	return orderMap, isNewLocalOrder
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
func MaxState(s1 types.OrderState, s2 types.OrderState) types.OrderState{
	if s2>s1 {
		return s2
	}
	return s1
}



func combineMaps(localID string, localMap map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order,
	 map2 map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order) (map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order, bool) { //returntype
		
		isNewLocalOrder := false
		
		for key, val:= range localMap {
			temp, keyExists := map2[key]
			tempLocal := val

			if keyExists {
			for i:= range val {
				for j:=range val[i]{
					s1 := tempLocal[i][j].State
					s1Was := s1
					s2 := temp[i][j].State
					if j==constants.N_BUTTONS-1 { //antar 0 indeksering, identifiserer cab call
						s2 = s1
						temp[i][j].Counter = tempLocal[i][j].Counter
					} else if tempLocal[i][j].Counter == temp[i][j].Counter { 
						switch {
						case (s1==types.OS_UnacceptedOrder && s2==types.OS_UnacceptedOrder):
							s1 = types.OS_AcceptedOrder
							s2 = s1
							break
						case ((s1-s2==2) || (s2-s1==2)):
							s1 = types.OS_AcceptedOrder
							s2 = s1
							break
						default :
							s1= MaxState(s1,s2)
							s2=s1
							break
							}	
						tempLocal[i][j].State=s1
						temp[i][j].State=s2
					} else {
						presedence:=getPresedence(tempLocal[i][j],temp[i][j])
						tempLocal[i][j].State = presedence.State
						temp[i][j].State = presedence.State
					}
					presedence := getPresedence(tempLocal[i][j], temp[i][j])
					tempLocal[i][j].Counter = presedence.Counter
					temp[i][j].Counter = presedence.Counter

					//Check if FSM should get local orders updated-signal
					if (key == localID) && (s1 != s1_was) && (s1 == types.OS_AcceptedOrder) {
						isNewLocalOrder = true
					} 
				}
				}
				map2[key] = temp
				localMap[key] = tempLocal
			}
			}
		
		return localMap, isNewLocalOrder
}

/*func removeDuplicates(elevOrders map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order) map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order {//returntype
	var tmpList [constants.N_ELEVATORS]string;
	
	index := 0

	for id, _ := range elev_orders{
		tmp_list[index]=id
		index++
	}
	for h:= 0; h < len(tmp_list); h++{
		elevator_1 := tmp_list[h]
		for i := h+1; i < len(tmp_list); i++ {
			elevator_2 := tmp_list[i]
			for j := range constants.N_FLOORS{
					for k := 0 ; k<constants.N_BUTTONS-1 {
						
						s1:=elev_orders[elevator_1][j][k].State
						s2:=elev_orders[elevator_2][j][k].State
						if (!(s1==types.OS_NoOrder))&&(!(s2==types.OS_NoOrder)){
							
							largest_id:= largestID(elevator_1,elevator_2)
							temp := elev_orders[largest_id]
							temp[j][k].State = types.OS_NoOrder
							temp[j][k].Counter++
							elev_orders[largest_id] = temp	
							
							}
						
					}
					
				}
			}

		
	}
	return elev_orders
}
*/
	

	
	
