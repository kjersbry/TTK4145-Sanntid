package merger

import (
	 "../constants"
	 "../types"
	 //"reflect"
	 //"fmt"
)

/*TODO!!!!!:
- test uten duplikat
- slukke lys ved sletting fra duplikat*/
//er counter alltid riktig?

func MergeOrders(local_ID string, local_elev map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order, elev_2 map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order) (map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order, bool) { 
	order_map, is_new_local_order := combineMaps(local_ID, local_elev, elev_2)
	order_map = removeDuplicates(order_map)
	return order_map, is_new_local_order
}

func getPresedence(order_1 types.Order, order_2 types.Order) types.Order {
	if order_2.Counter > order_1.Counter {
		return order_2
	} else {
		return order_1
	}

}

func largestID(elev_1 string, elev_2 string) string {
	/*short_1 := []rune(elev_1)
	short_2 := []rune(elev_2)
	short_1 = strings.Replace(short_1, "-", "", -1)
	short_1 = strings.Replace(short_1,"peer","")
	short_2 = strings.Replace(short_2, "-", "", -1)
	short_2 = strings.Replace(short_2,"peer","")*/
	if elev_1 > elev_2 {
		return elev_1
	}
	return elev_2
}
func MaxState(s1 types.OrderState, s2 types.OrderState) types.OrderState{
	if s2>s1 {
		return s2
	}
	return s1
}



func combineMaps(local_ID string, local_map map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order,
	 map_2 map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order) (map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order, bool) { //returntype
		
		is_new_local_order := false
		
		for key, val:= range local_map {
			temp, keyExists := map_2[key]
			temp_local := val

			if keyExists {
			for i:= range val {
				for j:=range val[i]{
					s1 := temp_local[i][j].State
					s1_was := s1
					s2 := temp[i][j].State
					if j==constants.N_BUTTONS-1 { //antar 0 indeksering, identifiserer cab call
						s2 = s1
						temp[i][j].Counter = temp_local[i][j].Counter
					} else if temp_local[i][j].Counter == temp[i][j].Counter { 
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
						temp_local[i][j].State=s1
						temp[i][j].State=s2
					} else {
						presedence:=getPresedence(temp_local[i][j],temp[i][j])
						temp_local[i][j].State = presedence.State
						temp[i][j].State = presedence.State
					}
					presedence := getPresedence(temp_local[i][j], temp[i][j])
					temp_local[i][j].Counter = presedence.Counter
					temp[i][j].Counter = presedence.Counter

					//Check if FSM should get local orders updated-signal
					if (key == local_ID) && (s1 != s1_was) && (s1 == types.OS_AcceptedOrder) {
						is_new_local_order = true
					} 
				}
				}
				map_2[key] = temp
				local_map[key] = temp_local
			}
			}
		
		return local_map, is_new_local_order
}

func removeDuplicates(elev_orders map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order) map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order {//returntype
	var tmp_list [constants.N_ELEVATORS]string;
	
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
	

	
	
