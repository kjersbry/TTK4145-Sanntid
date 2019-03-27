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

func Merger(heis1 types.Elevator, heis2 types.Elevator){
	if isOperational(heis1) && isOerational(heis2) {
		order_map:=Merger(heis1.ordermap, heis2.ordermap)
		order_map:=removeDuplicates(order_map)

		
	}
}
 


func getPresedence(order1 types.Order,order2 types.Order) types.Order{
	
	if order1.Counter > order2.Counter {
		return order1
	}
	else {
		return order2
	}

}

func largestID(elev1 string, elev2 string) string {
	short1:= []rune(elev1)
	short2:= []rune(elev2)
	short1 := strings.Replace(short1, "-", "", -1)
	short1 := strings.Replace(short1,"peer","")
	short2 := strings.Replace(short2, "-", "", -1)
	short2 := strings.Replace(short2,"peer","")
	if short1>short2 {
		return elev1
	}
	return elev2
}



func combineMaps(local_map map[string][constants.N_FLOORS][constants.N_BUTTONS],
	 map2 map[string][constants.N_FLOORS][constants.N_BUTTONS]) 
	 map[string][constants.N_FLOORS][constants.N_BUTTONS]{ //returntype
	eq := reflect.DeepEqual(local_map, map2)
	if eq {
		return local_map
	}
	else {
		for key, val:= range local_map {
			for i:=range val {
				for j:=range val[i]{
					s1:=local_map[key][i][j].State
					s2:=map2[key][i][j].State
					
					if j==constants.N_BUTTONS-1 { //antar 0 indeksering, identifiserer cab call
						s2:=s1;
						map2[key][i][j].Counter:=local_map[key][i][j].Counter
					}

					
					else if local_map[key][i][j].Counter==map2[key][i][j].Counter { 
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
						map2[key][i][j].State:=s2
					}
					else {
						presedence=getPresedence(local_map[key][i][j],map2[key][i][j])
						local_map[key][i][j].State:=presedence.State
						map2[key][i][j].State:=presedence.State
					}
					presedence=getPresedence(local_map[key][i][j],map2[key][i][j])
						local_map[key][i][j].Counter:=presedence.Counter
						map2[key][i][j].Counter:=presedence.Counter
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

	

	
	
