package states
import (
    "../elevio"
    "../globalconstants"
)


type ElevatorState int
const (
    ES_Idle ElevatorState = 0
    ES_DoorOpen           = 1
    ES_Moving             = 2
 )


type OrderState int
const (
    OS_NoOrder OrderState = 0
    OS_UnacceptedOrder    = 1
    OS_AcceptedOrder      = 2
)

type Order struct { 
    //litt midlertidig oppsett, hvordan skal denne være?
    ID int64
    State OrderState
}

type Elevator struct {
    Elevator_ID int
    State ElevatorState
    Floor int
    Direction elevio.MotorDirection   //does only change to stop when IDLE, not when stopping for order
    Orders [globalconstants.N_FLOORS][globalconstants.N_BUTTONS]Order
}



