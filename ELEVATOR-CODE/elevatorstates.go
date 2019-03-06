package elevatorstates
import "elevio"

typedef enum {
    ES_Idle,
    ES_DoorOpen,
    ES_Moving
} ElevatorState;


typedef struct {
    elevator_ID int;
    state ElevatorState;
    floor int;
    direction MotorDirection;   
} Elevator;

var elevator Elevator = {-1, ES_Idle, -1, MD_Stop} //NB midlertidig ID

//Forslag: detects changes in elevator direction variable and sets the motor
//da trenger vi ikke huske på å bruke setmotordir overalt
//men da må vi passe på chooseDirection i orders
/*func PollAndSetDirection(){
    prev:= elevator.direction
    for{
        if(prev!=elevator.direction){
            SetMotorDirection(elevator.direction)
        }
    }
}*/

func InitElevator(drv_floors <-chan int){
    //somehow initialize ID uniquely
    select{
    case floor:= <- drv_floors:
        elevator.floor = floor
    }
    if(elevator.floor == -1){
        //init between floors
        elevator.direction = MD_Down
        SetMotorDirection(elevator.direction)
        elevator.state = ES_Moving
    }
}


/*get-funcs for states? 
--> fsm can get instead of passing as arguments 
for instance: on_floor_arrival(floor), not passing current_dir
*/