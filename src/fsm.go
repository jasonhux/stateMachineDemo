package main

import (
	"fmt"

	"github.com/stateMachineDemo/src/state"
	trans "github.com/stateMachineDemo/src/transition"
)

type FSM struct {
	currentState string
	stateMap     map[string]details
}

type details struct {
	isEnd         bool
	transitionMap map[string]transition
}

type transition struct {
	toState string
	action  func(transitionKey string) error
}

func NewFSM() *FSM {
	return &FSM{
		currentState: state.INIT,
		stateMap: map[string]details{
			state.INIT: details{
				transitionMap: map[string]transition{
					trans.START_UPDATE: transition{
						toState: state.UPDATING,
						action:  mockCallBack,
					},
				},
			},
			state.UPDATING: details{
				transitionMap: map[string]transition{
					trans.FINISH_UPDATE_WITH_SUCCESS: transition{
						toState: state.UPDATED,
					},
					trans.FINISH_UPDATE_WITH_ERROR: transition{
						toState: state.UPDATEFAILED,
					},
				},
			},
			state.UPDATED: details{
				transitionMap: map[string]transition{
					trans.COMPLETE_TASK: transition{
						toState: state.DONE,
					},
				},
			},
			state.UPDATEFAILED: details{
				transitionMap: map[string]transition{
					trans.START_UPDATE: transition{
						toState: state.UPDATING,
					},
					trans.COMPLETE_TASK: transition{
						toState: state.DONE,
					},
				},
			},
			state.DONE: details{
				isEnd: true,
			},
		},
	}
}

func (fsm *FSM) sendTransition(transitionKey string) (err error) {
	details, ok := fsm.stateMap[fsm.currentState]
	if !ok {
		return fmt.Errorf("current state <%s> does not have valid details", fsm.currentState)
	}
	if details.isEnd {
		return
	}
	transition, ok := details.transitionMap[transitionKey]
	if !ok {
		return fmt.Errorf("current state <%s> does not allow given transition <%s>", fsm.currentState, transitionKey)
	}
	if transition.action != nil {
		if err = transition.action(transitionKey); err != nil {
			return
		}
	}

	fmt.Printf("New state <%s> is created from <%s> by transition <%s>\n", transition.toState, fsm.currentState, transitionKey)
	fsm.currentState = transition.toState
	return
}

func mockCallBack(transitionKey string) error {
	//in real world, the transition action can be actually updating the cluster, here just put in a dumo func and return nil error
	fmt.Printf("action for transition <%s> is called\n", transitionKey)
	return nil
}
