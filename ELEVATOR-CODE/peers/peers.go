package peers

import (
	"fmt"
	"net"
	"time"

	"../conn"
	"../types"
)

const interval = 15 * time.Millisecond
const timeout = 50 * time.Millisecond

func ConnectionTransmitter(port int, id string) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	for {
		select {
		case <-time.After(interval):
		}
		conn.WriteTo([]byte(id), addr)
	}
}

//Detects when an elevator has been disconncted or reconnected. Does this need a sleep in the loop?
func ConnectionObserver(port int, connectionUpdate chan<- types.Connection_Event) {

	var buf [1024]byte
	var lostConnections []string
	var update types.Connection_Event
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// Adding new connection

		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {

				for i, ID := range lostConnections {
					if ID == id {
						update = types.Connection_Event{ID, true}
						connectionUpdate <- update
						lostConnections = append(lostConnections[:i], lostConnections[i+1:]...)
					}
				}
			}
			lastSeen[id] = time.Now()
		}

		//Removing lost connection

		for elevID, lastTime := range lastSeen {
			if time.Now().Sub(lastTime) > timeout { //Where is timeout???
				lostConnections = append(lostConnections, elevID)
				delete(lastSeen, elevID)
				update = types.Connection_Event{elevID, false}
				connectionUpdate <- update
			}
		}
	}

}
