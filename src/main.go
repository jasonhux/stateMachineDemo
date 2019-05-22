package main

import (
	"fmt"

	trans "github.com/stateMachineDemo/src/transition"
)

func main() {
	fsm := NewFSM()
	err := fsm.sendTransition(trans.START_UPDATE)
	if err != nil {
		fmt.Println(err)
		return
	}
	//perform retry
	for i := 0; i < 3; i++ {
		isSuccessful := mockVerifyClusterUpdate(i)
		if isSuccessful {
			err = fsm.sendTransition(trans.FINISH_UPDATE_WITH_SUCCESS)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = fsm.sendTransition(trans.COMPLETE_TASK)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Successfully update the cluster with <%v> retry\n", i)
			return
		}
		err = fsm.sendTransition(trans.FINISH_UPDATE_WITH_ERROR)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = fsm.sendTransition(trans.START_UPDATE)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Failed to update the cluster")
}
func mockVerifyClusterUpdate(i int) (isSuccessful bool) {
	//mock verify func to return true on 2nd retry
	return i == 2
}
