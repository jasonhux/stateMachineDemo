package main

func main() {
	fsm := NewFSM("removed")
	fsm.run()
	fsm = NewFSM("deployed")
	fsm.run()
}
