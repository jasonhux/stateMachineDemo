package main

import "fmt"

type FSM struct {
	currentState string
	stateMap     map[string]details
}

type details struct {
	isEnd               bool
	determineTransition func() (transitionKey string, err error)
	transitionMap       map[string]transition
}

type transition struct {
	toState              string
	beforeTransitionFunc func() error
}

func NewFSM(initialState string) *FSM {
	return &FSM{
		currentState: initialState,
		stateMap: map[string]details{
			"removed": details{
				determineTransition: mockDetermineTransition,
				transitionMap: map[string]transition{
					"deploy": transition{
						toState:              "deployed",
						beforeTransitionFunc: mockBeforeTransitionFunc,
					},
				},
			},
			"deployed": details{
				determineTransition: mockDetermineTransition,
				transitionMap: map[string]transition{
					"configure": transition{
						toState:              "configured",
						beforeTransitionFunc: mockBeforeTransitionFunc,
					},
				},
			},
			"configured": details{
				determineTransition: mockDetermineTransition,
				transitionMap: map[string]transition{
					"complete": transition{
						toState: "completed",
					},
					"remove": transition{
						toState:              "removed",
						beforeTransitionFunc: mockBeforeTransitionFunc,
					},
				},
			},
			"completed": details{
				isEnd: true,
			},
		},
	}
}

func (fsm *FSM) run() (endState string, err error) {
	var transitionKey string

	for {
		stateDetails := fsm.stateMap[fsm.currentState]

		//if state is the end state, i.e. completed, then exit the fsm
		if stateDetails.isEnd {
			fmt.Printf("End state <%s> is reached\n", fsm.currentState)
			endState = fsm.currentState
			break
		}

		//in case a state has multiple next transitions, need to make a decision which transition will be chose next
		//this can be determined by a cluster status check, for example, after configuration we check the specs etc but still find something not correct -- validation, then we decide to remove the new created cluster etc
		//in demo, we simply choose the first one
		for k := range stateDetails.transitionMap {
			transitionKey = k
			break
		}

		//in real world, we should have the transition key determined inside of determineTransition func; again for demo purpose, the func is just to return the input transitin key;
		//stateDetails.determineTransition()
		transition, ok := stateDetails.transitionMap[transitionKey]
		if !ok {
			endState = fsm.currentState
			err = fmt.Errorf("For state <%s>, it can not take transition <%s>", fsm.currentState, transitionKey)
			break
		}
		//run before transition
		if transition.beforeTransitionFunc != nil {
			err := transition.beforeTransitionFunc()
			if err != nil {
				//can do retry on beforeTransitionFunc
				//here for demo purpose just break the loop
				break
			}
		}
		//only update the fsm state if before transition func is successful
		fmt.Printf("A transition <%s> is made from state <%s> to <%s>\n", transitionKey, fsm.currentState, transition.toState)
		fsm.currentState = transition.toState
		continue
	}
	return
}

func mockBeforeTransitionFunc() error {
	//in real world, the func will be vary for different transition, such as implement some work to deploy the cluster/node, return err if deploy not successful
	return nil
}

func mockDetermineTransition() (transitionKey string, err error) {
	//add logic here to determine next transition
	givenNextTransitionKey := "deploy"
	return givenNextTransitionKey, nil
}
