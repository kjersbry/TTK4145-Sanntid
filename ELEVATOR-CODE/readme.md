This document gives a brief overveiw over the modules of the system. An overveiw over libraries that has not been written by the authors will also be given.

Modules:
The core of this system is main, elevatorstates, bcast and FSM. FSM (finite state machine) is responsible for normal operation of a single elevator. 
Bcast is used to transmitt and recive data from external systems. Both external and locla data is stored and handled in the elevatorstates module. 
Take special note of the map "allElevators" that is declared in elevatorstates. This map is only written to by one function, namley "updateElevator"
All other modules should be selfexplanetory.

Libraies used in this project:
This project has made frequent use of standard GO libraies. These aren't very esoteric and should be familiar to any decent GO programmers. Further, libraries written by the school faculty have also been frequently used. Most notable is "bcast" and "elevio". The modulue named "Peers" has been retrofitted to suit the needs of our design. Thus, it should be taken into account when reveiwing the quality the code.