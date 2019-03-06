package orderassigner
import "elevio"


func assignOrders(drv_button <-chan ButtonEvent/*, write_order_request chan<- order (some type)*/){
	for{
		select{
		case order:= <- ButtonEvent:
			//her legges assignment algorithm
			//en heis---> kun dette:
			//write_order_request <- order //skriver resultat til order
		}
	}	
}