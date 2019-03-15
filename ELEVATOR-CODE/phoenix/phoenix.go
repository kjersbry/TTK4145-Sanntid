package main

import (
  "fmt"
  "time"
  "../bcast"
  "os/exec"
  "math"
)


func main(){
  //assume that this process is backup
  heartbeat_rx:= make(chan int)

  go bcast.Receiver(1234, heartbeat_rx)
  go heartbeat_receiver(heartbeat_rx)

  for{
    /*fix?*/
    time.Sleep(time.Millisecond*1000000)
  }
}

/*******************BACKUP FUNC****************************/
func heartbeat_receiver(heartbeat_rx <-chan int){
    count_backup:= 0
    last_time:= time.Now()
    for {
      select{
      case count_backup = <- heartbeat_rx:
          last_time = time.Now()
          //fmt.Printf("\nback: %d\n", count_backup)
      default:
        if float64(time.Now().Sub(last_time)) > 100*math.Pow(10,6) {
          fmt.Printf("PRIMARY DIED  :(")
          primary(count_backup)
          return
        }
      }
    }
}

/*******************PRIMARY FUNCS****************************/
func primary(counter int) {
  ctr_ch := make(chan int)
  heartbeat_tx := make(chan int)
  fmt.Printf("\nStarting new process..\n")
  newProcess := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run phoenix.go")
  err := newProcess.Run()
  if(err!= nil){
    fmt.Printf("\nCould not spawn new process\n")
  }
  go incrementer(counter, ctr_ch)
  go bcast.Transmitter(1234,heartbeat_tx)
  go heartbeat_primary(ctr_ch, heartbeat_tx)
}

func incrementer(counter int, ctr_ch chan<- int) {
  for{
    fmt.Printf("counter: %d\n", counter)
    ctr_ch <- counter
    counter++
    time.Sleep(time.Millisecond*1000)
  }
}

func heartbeat_primary(ctr_ch <-chan int, heartbeat_tx chan<- int) {
  var count int
  for {
    time.Sleep(time.Millisecond*20)
    select{
    case  count =<- ctr_ch:
    default:
      heartbeat_tx <- count
    }
  }
}
