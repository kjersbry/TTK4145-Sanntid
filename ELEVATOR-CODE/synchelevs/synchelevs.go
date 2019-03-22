package synchelevs

import (
	//"../states"
	"../types"
	//../constants"
	//fmt"
	//"time"
)

/*Jsonencodingen skal visstnok kunne ta inn map, men test det*/
/*phoenix: sender hele heisen eller all_elevs i heartbeat, hvis den ikke får noe må den lage ny*/

/* Tilstandene må oppdateres fra FSM. Men:
Husk å tenke gjennom om vi kan få lese-skrivefeil på elevinfopacket
Etter ordermerge må vi skrive til fsm
*/
 


func TransmitElev(elev_tx chan<- types.ElevInfoPacket){
	for {
		//elev_tx <- all_elevs
		//time.Sleep(time.Milliseconds*constants.TRANSMIT_MS)
	}
}

func ReceiveElevs(elev_rx <-chan types.ElevInfoPacket /*, updatefromfsm chan */){
	for{
		select{
		/*case received := <- elev_rx:
			for key, val := range received.Elev_map {
				if !states.KeyExists(key) {
					//all_elevs.Elev_map[key] = val
				} else {
					//update states
					if received.Elev_ID == key {//overwrite new states where elev = elev which sent
						//all_elevs.Elev_map[key].State = val.State
						//all_elevs.Elev_map[key].Floor = val.Floor
						//all_elevs.Elev_map[key].Direction = val.Direction
					}

				}
			}
			//merge orders
			//funksjon for det
			//her, eller ett hakk ut?
		//case fromfsm: update all_elevs.map[this elev]*/
		}
	}	
}

