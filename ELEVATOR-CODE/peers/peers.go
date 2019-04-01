package peers

import (
	"fmt"
	"net"
	"time"

	"../conn"
	"../types"
)

const interval = 15 * time.Millisecond
const timeout = 500 * time.Millisecond

func ConnectionTransmitter(port int, localID string) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	for {
		select {
		case <-time.After(interval):
		}
		conn.WriteTo([]byte(localID), addr)
	}
}

func ConnectionObserver(port int, connectionUpdate chan<- types.ConnectionEvent, localID string) {

	var buf [1024]byte
	var lostConnections []string
	var update types.ConnectionEvent
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		if (id != "") && (id != localID) {
			if _, idExists := lastSeen[id]; !idExists {

				for i, savedID := range lostConnections {
					if savedID == id {
						update = types.ConnectionEvent{savedID, true}
						connectionUpdate <- update
						lostConnections = append(lostConnections[:i], lostConnections[i+1:]...)
					}
				}
			}
			lastSeen[id] = time.Now()
		}

		for elevID, lastTime := range lastSeen {
			if time.Now().Sub(lastTime) > timeout {
				update = types.ConnectionEvent{elevID, false}
				connectionUpdate <- update
				lostConnections = append(lostConnections, elevID)
				delete(lastSeen, elevID)
			}
		}
	}

}
