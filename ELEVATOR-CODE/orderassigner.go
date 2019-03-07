package orderassigner
import "elevio"


func assignOrders(drv_button <-chan ButtonEvent, add_order chan<- elevio.ButtonEvent){
	for{
		select{
		case order:= <- ButtonEvent:
			//her legges assignment algorithm
			//en heis---> kun dette:
			add_order <- order //skriver resultat til order
		}
	}	
}