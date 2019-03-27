package main

import (
	"./merger"
	"./types"
	"./constants"
)

func main(){
	id_liste:=[3]string {"1", "2" ,"3"}
	test_1 := make(map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order)
	test_2 := make(map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order)
	for i:=0; i<3; i++ {
		var ord [constants.N_FLOORS][constants.N_BUTTONS](types.Order)
		test_1[id_liste[i]]= ord
		test_2[id_liste[i]]= ord
	}
	var ord_1 [constants.N_FLOORS][constants.N_BUTTONS](types.Order)
	var ord_2 [constants.N_FLOORS][constants.N_BUTTONS](types.Order)
	ord_1[2][2].State=1
	ord_1[2][2].Counter=1
	ord_1[3][2].State=2
	ord_1[3][2].Counter=1
	ord_2[3][2].State=1
	ord_2[3][2].Counter=1

	test_1["3"]=ord_1
	test_2["2"]=ord_2

	order_map:=merger.Merger(test_1,test_2)
	
	types.PrintOrders2(order_map["3"])
	//map[string][constants.N_FLOORS][constants.N_BUTTONS]types.Order
	//types.PrintOrders2(map[key])

}