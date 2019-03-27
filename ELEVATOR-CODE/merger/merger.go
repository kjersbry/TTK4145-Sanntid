import (
	"time"
	"sync"
 	"net"
 	"fmt"
 	"../constants"
 	"reflect"
 	"math"
 	"strings"
)

func Merger(local_elev types.Wrapped_elevator, elev_2 types.Wrapped_elevator)
map[string][constants.N_FLOORS][constants.N_BUTTONS]{ 
	if local_elev.Is_Operational && elev_2.Is_Operational {
		order_map:=Merger(local_elv.ordermap, elev_2.ordermap)
		order_map:=removeDuplicates(order_map)
		return order_map
		
	}
}
 


func getPresedence(order_1 types.Order,order_2 types.Order) types.Order{
	
	if order_1.Counter > order_2.Counter {
		return order_1
	}
	else {
		return order_2
	}

}

func largestID(elev_1 string, elev_2 string) string {
	short_1:= []rune(elev_1)
	short_2:= []rune(elev_2)
	short_1 := strings.Replace(short_1, "-", "", -1)
	short_1 := strings.Replace(short_1,"peer","")
	short_2 := strings.Replace(short_2, "-", "", -1)
	short_2 := strings.Replace(short_2,"peer","")
	if short1>short2 {
		return elev_1
	}
	return elev_2
}



func combineMaps(local_map map[string][constants.N_FLOORS][constants.N_BUTTONS],
	 map_2 map[string][constants.N_FLOORS][constants.N_BUTTONS]) 
	 map[string][constants.N_FLOORS][constants.N_BUTTONS]{ //returntype
	eq := reflect.DeepEqual(local_map, map_2)
	if eq {
		return local_map
	}
	else {
		for key, val:= range local_map {
			for i:=range val {
				for j:=range val[i]{
					s1:=local_map[key][i][j].State
					s2:=map_2[key][i][j].State
					
					if j==constants.N_BUTTONS-1 { //antar 0 indeksering, identifiserer cab call
						s2:=s1;
						map_2[key][i][j].Counter:=local_map[key][i][j].Counter
					}

					
					else if local_map[key][i][j].Counter==map_2[key][i][j].Counter { 
						switch {
						case (s1=1 && s2=1):
								s1,s2 := 2
								break
							case math.Abs(s1-s2)==2:
								s1, s2 := 0
								break
							default :
								s1, s2:= math.Max(s1,s2):
								break
							}	
						local_map[key][i][j].State:=s1
						map_2[key][i][j].State:=s2
					}
					else {
						presedence=getPresedence(local_map[key][i][j],map_2[key][i][j])
						local_map[key][i][j].State:=presedence.State
						map_2[key][i][j].State:=presedence.State
					}
					presedence=getPresedence(local_map[key][i][j],map_2[key][i][j])
						local_map[key][i][j].Counter:=presedence.Counter
						map_2[key][i][j].Counter:=presedence.Counter
				}
						
					
				}
			}
		}
		return local_map
	 }
}

func removeDupliacates(elev_orders map[string][constants.N_FLOORS][constants.N_BUTTONS])
map[string][constants.N_FLOORS][constants.N_BUTTONS] {//returntype
	var tmp_list [constants.N_ELEVATORS]string;
	
	var index:=0;

	for id, val:= range elev_orders{
		tmp_list[index]:=id
		index++
	}
	for h:= range tmp_list-1{
		elevator_1:=tmp_list[h]
		for i:= h; i<=len(tmp_list-1)
			elevator_2:=tmp_list[i+1]
			for j:=range elev_orders[elevator_1]{
				floor:=elev_orders[elevator_1][j]{
					for k:= range floor {
						
						s1:=floor[k].State
						s2=elev_orders[elevator_2][j][k].State
						if (!(s1==0))&&(!(s2==0)){
							is_local, cab :=isLocalCab(k,elevator_1,elevator_2)	
							if is_local {
								elev_orders[cab][j][k].States
							}
							largest_id:=largestID(elevator_1,elevator_2)
							elev_orders[largest_id][j][k].State:=0
							elev_orders[largest_id][j][k].Counter++
								
							}
						
					}
					
				}
			}

		}
		return elev_orders
	}

	

	
	
