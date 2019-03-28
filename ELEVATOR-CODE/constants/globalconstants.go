package constants
/*defines some global constants. Makes it easy to change number
of elevators and floors*/

import "time"

const N_FLOORS int = 4
const N_BUTTONS int = 3
const N_ELEVATORS int = 3

const DOOR_OPEN_SEC = 3
const TRANSMIT_MS = 200
const ELEVATOR_TIMEOUT = 15000 * time.Millisecond
