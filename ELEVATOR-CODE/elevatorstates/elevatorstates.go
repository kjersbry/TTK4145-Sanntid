package elevatorstates
import "../elevio"


type ElevatorState int
const (
    ES_Idle ElevatorState = 0
    ES_DoorOpen           = 1
    ES_Moving             = 2
 )


type Elevator struct {
    Elevator_ID int
    State ElevatorState
    Floor int
    Direction elevio.MotorDirection   //does only change to stop when IDLE, not when stopping for order
   // orderlist-ENKELTHEIS Orderstype
}

var elevator Elevator = Elevator{-1, ES_Idle, -1, elevio.MD_Stop} //NB midlertidig ID

//Read function
func ReadElevator() Elevator {
    return elevator
}

func UpdateElevator(update_ID <-chan int, update_state <-chan ElevatorState,
    update_floor <-chan int, update_direction <-chan elevio.MotorDirection){
        for {
            select{
            case new_id:= <- update_ID:
                elevator.Elevator_ID=new_id
            case new_state:= <- update_state:
                elevator.State= new_state
            case new_floor:= <- update_floor:
                elevator.Floor=new_floor
            case new_dir:= <- update_direction:
                elevator.Direction = new_dir
            }
        }
}


func InitElevator(drv_floors <-chan int){
    //somehow initialize ID uniquely
    select{
    case floor:= <- drv_floors:
        elevator.Floor = floor
    }
    if(elevator.Floor == -1){
        //init between floors:
        elevator.Direction = elevio.MD_Down
        elevio.SetMotorDirection(elevator.Direction)
        elevator.State = ES_Moving
    }
}

//Forslag: detects changes in elevator direction variable and sets the motor
//da trenger vi ikke huske på å bruke setmotordir overalt
//men da må vi passe på chooseDirection i orders
//tror vi dropper dette
/*func PollAndSetDirection(){
    prev:= elevator.direction
    for{
        if(prev!=elevator.direction){
            SetMotorDirection(elevator.direction)
        }
    }
}*/
