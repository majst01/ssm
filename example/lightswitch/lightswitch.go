package main

import (
	"fmt"

	"github.com/majst01/ssm"
)

const (
	Off ssm.StateType = "Off"
	On  ssm.StateType = "On"

	SwitchOff ssm.EventType = "SwitchOff"
	SwitchOn  ssm.EventType = "SwitchOn"
)

// OffAction represents the action executed on entering the Off state.
type OffAction struct{}

func (a *OffAction) Execute(eventCtx ssm.EventContext) ssm.EventType {
	fmt.Println("The light has been switched off")
	return ssm.NoOp
}

// OnAction represents the action executed on entering the On state.
type OnAction struct{}

func (a *OnAction) Execute(eventCtx ssm.EventContext) ssm.EventType {
	fmt.Println("The light has been switched on")
	return ssm.NoOp
}

func newLightSwitchFSM() *ssm.StateMachine {
	return &ssm.StateMachine{
		States: ssm.States{
			ssm.Default: ssm.State{
				Events: ssm.Events{
					SwitchOff: Off,
				},
			},
			Off: ssm.State{
				Action: &OffAction{},
				Events: ssm.Events{
					SwitchOn: On,
				},
			},
			On: ssm.State{
				Action: &OnAction{},
				Events: ssm.Events{
					SwitchOff: Off,
				},
			},
		},
	}
}
