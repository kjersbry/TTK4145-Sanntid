package orderassigner
import "../elevio"


//This will be executed as a goroutine in main
func AssignOrder(drv_button <-chan elevio.ButtonEvent, add_order chan<- elevio.ButtonEvent){
	for{
		select{
		case order:= <- drv_button:
			//her legges assignment algorithm
			//en heis---> kun dette:
			add_order <- order //skriver resultat til order
		}
	}	
}