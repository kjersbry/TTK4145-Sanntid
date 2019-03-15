package states
import (
    "../elevio"
    "../globalconstants"
    "fmt"
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


func orderToString(o Order) string {
  switch(o.State){
  case OS_NoOrder:
    return "No"
  case OS_AcceptedOrder:
  case OS_UnacceptedOrder:
  }
  return "Yes"
}

func PrintOrders(e Elevator){
  fmt.Printf("\n\n-----Queue------\n")
  for i:= 0; i < globalconstants.N_FLOORS; i++{
    for j:= 0; j < globalconstants.N_BUTTONS; j++{
      fmt.Printf("%s ", orderToString(e.Orders[i][j]))
    }
    fmt.Printf("\n")
  }
  fmt.Printf("\n----------------\n\n")
}
