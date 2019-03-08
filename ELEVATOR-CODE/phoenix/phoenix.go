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
  //process:= "backup"
  heartbeat_rx:= make(chan int)

  go bcast.Receiver(1234, heartbeat_rx)
  go heartbeat_receiver(heartbeat_rx)

  for{
    /*Fix this*/
    time.Sleep(time.Millisecond*1000000)
    time.Sleep(time.Millisecond*1000000)

    time.Sleep(time.Millisecond*1000000)

    time.Sleep(time.Millisecond*1000000)

    time.Sleep(time.Millisecond*1000000)

  }
}



func primary(counter int) {
  ctr_ch := make(chan int)
  heartbeat_tx := make(chan int)
  /* start secoundary process
  input -> -is_primary = "False"
  -counter = noe
  */
  newProcess := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run phoenix.go")
  err := newProcess.Run()
  if(err!= nil){

  }
  go incrementer(counter, ctr_ch)
  go bcast.Transmitter(1234,heartbeat_tx)
  go heartbeat_primary(ctr_ch, heartbeat_tx)

}

//Primary functions
func incrementer(counter int, ctr_ch chan<- int) {
  for{
    fmt.Printf("%d\n", counter)
    ctr_ch <- counter
    counter++
    time.Sleep(time.Millisecond*1000)
  }
}

func heartbeat_primary(ctr_ch <-chan int, heartbeat_tx chan<- int) {
  var cout int
  for {
    time.Sleep(time.Millisecond*20)
    select{
    case  cout =<- ctr_ch:
    }
    heartbeat_tx <- cout
    }
}

//Secoundary functions
func backup(){
  //if no heartbeat: spawn
  //var count int
  heartbeat_rx:= make(chan int)

  go bcast.Receiver(1234, heartbeat_rx)
  go heartbeat_receiver(heartbeat_rx)
}

func heartbeat_receiver(heartbeat_rx <-chan int){

    var backup int
    last_time:= time.Now()

    /*
    TODO

    Make for select and time check different threads.
    Seems like its waiting for input in the case


    */
    for {
      select{
      case backup = <- heartbeat_rx:
          last_time = time.Now()
          fmt.Printf("%d", backup)
      }
      fmt.Printf("\nheyyy\n") //kommer aldri hit

      if float64(time.Now().Sub(last_time)) > 100*math.Pow(10,6) {
        fmt.Printf("PRIMARY DIED  :(")
        //kill old primary in case it turns out it was alive
        //run primary procedure
        // with backup as input argument
      }
    }

}
