package ssm

import (
	"errors"
	"sync"
)

var (
	// ErrEventRejected is the error returned when the state machine cannot process
	// an event in the state that it is in.
	ErrEventRejected = errors.New("event rejected")
	// ErrConfigurationInvalid is the error returned when the state machine is configured incorrectly.
	ErrConfigurationInvalid = errors.New("state machine configuration invalid")
)

const (
	// Default represents the default state of the system.
	Default StateType = ""

	// NoOp represents a no-op event.
	NoOp EventType = "NoOp"
)

type (
	// StateType represents an extensible state type in the state machine.
	StateType string

	// EventType represents an extensible event type in the state machine.
	EventType string

	// EventContext represents the context to be passed to the action implementation.
	EventContext interface{}

	// Action represents the action to be executed in a given state.
	Action interface {
		Execute(eventCtx EventContext) EventType
	}

	// Events represents a mapping of events and states.
	Events map[EventType]StateType

	// State binds a state with an action and a set of events it can handle.
	State struct {
		Action Action
		Events Events
	}

	// States represents a mapping of states and their implementations.
	States map[StateType]State

	// StateMachine represents the state machine.
	StateMachine struct {
		// Previous represents the previous state.
		Previous StateType

		// Current represents the current state.
		Current StateType

		// States holds the configuration of states and events handled by the state machine.
		States States

		// mutex ensures that only 1 event is processed by the state machine at any given time.
		mutex sync.Mutex
	}
)

// getNextState returns the next state for the event given the machine's current
// state, or an error if the event can't be handled in the given state.
func (s *StateMachine) getNextState(event EventType) (StateType, error) {
	if state, ok := s.States[s.Current]; ok {
		if state.Events != nil {
			if next, ok := state.Events[event]; ok {
				return next, nil
			}
		}
	}
	return Default, ErrEventRejected
}

// SendEvent sends an event to the state machine.
func (s *StateMachine) SendEvent(event EventType, eventCtx EventContext) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for {
		// Determine the next state for the event given the machine's current state.
		nextState, err := s.getNextState(event)
		if err != nil {
			return ErrEventRejected
		}

		// Identify the state definition for the next state.
		state, ok := s.States[nextState]
		if !ok || state.Action == nil {
			// configuration error
			return ErrConfigurationInvalid
		}

		// Transition over to the next state.
		s.Previous = s.Current
		s.Current = nextState

		// Execute the next state's action and loop over again if the event returned
		// is not a no-op.
		nextEvent := state.Action.Execute(eventCtx)
		if nextEvent == NoOp {
			return nil
		}
		event = nextEvent
	}
}
