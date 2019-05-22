package main

import "fmt"

func main() {
	fsm := NewFSM()
	endState, _ := fsm.startWithState("removed")
	fmt.Printf("End state <%s> is reached\n", endState)
	endState, _ = fsm.startWithState("deployed")
	fmt.Printf("End state <%s> is reached\n", endState)

}
