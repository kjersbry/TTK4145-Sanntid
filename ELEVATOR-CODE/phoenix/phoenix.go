package phoenix

import (
  "../bcast"
  "../localip"
  "fmt"
  "time"
  "math"
  "strconv"
  "os/exec"
 // "os"
)
/*Module for running a backup and primary process for any elevator algorithm elevFunc.
elevFunc: the process that should be run in primary*/

var server_port int
var phoenix_port int

type RunElevFunc func(string, string)
var elevFunc RunElevFunc

/*******************BACKUP FUNCS****************************/
func RunBackup(phx_port int, srv_port int, elevatorFunc RunElevFunc){
  phoenix_port = phx_port
  server_port = srv_port
  heartbeat_rx:= make(chan string)
  end_backup := make(chan int)

  elevFunc = elevatorFunc

  go bcast.Receiver(phoenix_port, heartbeat_rx)
  go backupHeartbeatReceiver(heartbeat_rx, end_backup)

	for{
		select{
    case <- end_backup:
      return
		}
	}
}

func backupHeartbeatReceiver(heartbeat_rx <-chan string, end_backup chan<- int){
    local_ID := GetPeerID()
    last_time:= time.Now()
    for {
      select{
      case local_ID = <- heartbeat_rx:
          last_time = time.Now()
      default:
        if float64(time.Now().Sub(last_time)) > 100*math.Pow(10,6) {
          fmt.Printf("\nback: %s\n", local_ID)

          runPrimary(local_ID)
          end_backup <- 1
          return
        }
      }
    }
}

func GetPeerID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISElevatorCONNECTED"
	}
	return fmt.Sprintf("peer-%s", localIP/*, os.Getpid()*/)
}

/*******************PRIMARY FUNCS****************************/
func runPrimary(local_ID string) {
  spawnNewProcess()

  heartbeat_tx := make(chan string)
  go bcast.Transmitter(phoenix_port, heartbeat_tx)
  go primarySendHeartbeats(local_ID, heartbeat_tx)
  elevFunc(local_ID, strconv.Itoa(server_port)) 

	/*Infinite loop:, tror ikke den er nødvendig siden det er det inni elevFunc*/
	end_primary := make(chan int)
	for{
		select{
		case <- end_primary:
		}
	}
}

func SpawnTerminal(arg string) error {
  newProcess := exec.Command("gnome-terminal", "-x", "sh", "-c", arg)
  err := newProcess.Run()
  return err
}

func spawnNewProcess(){
  fmt.Printf("\nStarting new process..\n")
  err := SpawnTerminal("go run main.go -sport=" + strconv.Itoa(server_port) + " -pport=" + strconv.Itoa(phoenix_port))
  
  if(err!= nil){
    fmt.Printf("\nCould not spawn new process\n")
  }
}

func primarySendHeartbeats(local_ID string, heartbeat_tx chan<- string) {
  for {
    time.Sleep(time.Millisecond*20)
      heartbeat_tx <- local_ID
  }
}
