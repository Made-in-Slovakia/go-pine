package pine

import (
	"errors"
)

// Exchange is a struct that represents single exchange between the client and the emulator.
//
// Core of exchange is triplet [Command] [Answer] and callback function. Each triplet has an index
// that can be used to access it.
//
// While pair of [Command] and [Answer] is oblivious, callback function is used by function 
// [Exchange.RunCallbacks] to process data returned by emulator.
type Exchange struct {
	client    *Client
	commands  []Command
	answers   []Answer
	callbacks []func(answer Answer)
	closed    bool
}

// Returns a new [Exchange]. The given capacity is used as capacity for slices used internally by
// [Exchange].
func NewExchange(client *Client, capacity int) *Exchange {
	return &Exchange{
		client:    client,
		commands:  make([]Command, 0, capacity),
		answers:   make([]Answer, 0, capacity),
		callbacks: make([]func(Answer), 0, capacity),
		closed:    false,
	}
}

// Adds given [Command] to [Exchange] and returns index that represents position of the [Command] in
// internal store.
func (e *Exchange) AddCommand(command Command) (int, error) {
	return e.AddCommandWithCallback(command, func(answer Answer) {})
}

// Adds given [Command] and callback function to [Exchange] and returns index that represents
// position of the [Command] in internal store.
func (e *Exchange) AddCommandWithCallback(command Command, callback func(Answer)) (int, error) {
	if e.closed {
		return 0, errors.New("can't add command to closed exchange")
	}
	e.commands = append(e.commands, command)
	e.callbacks = append(e.callbacks, callback)
	return len(e.commands) - 1, nil
}

// Reports whether the [Exchange] is closed.
//
// [Exchange] is closed after execcution. Closed [Exchange] do not accept new [Command] and [Answer]
// can only be read from the closed [Exchange].
func (e *Exchange) IsClosed() bool {
	return e.closed
}

// Executes echange between the client and the emultor. Closes [Exchange].
func (e *Exchange) Execute() error {
	if e.closed {
		// When exchange is already closed, we can safely ignore this call.
		return nil
	}
	answers, err := e.client.SendCommands(e.commands)
	if err != nil {
		return err
	}

	e.answers = answers
	e.closed = true

	return nil
}

// Runs all callbacks.
//
// Callbacks can be executed only when [Exchange] is closed.
func (e *Exchange) RunCallbacks() error {
	if !e.closed {
		return errors.New("can't read from open exchange")
	}
	if len(e.callbacks) != len(e.answers) {
		// This error should never happen, it is here just for safety
		return errors.New("different amount of answers than callbacks")
	}

	for i, a := range e.answers {
		e.callbacks[i](a)
	}
	return nil
}

// Runs [Exchange.Execute] followed by [Exchange.RunCallbacks] and calls the given callback function
// with error or nil.
func (e *Exchange) ExecuteAndRunCallbacks(callback func(error)) {
	// TODO check errors
	if err := e.Execute(); err != nil {
		callback(err)
		return
	}
	if err := e.RunCallbacks(); err != nil {
		callback(err)
		return
	}
	callback(nil)
}

// Returns [Answer] content at the given position as uint8.
func (e *Exchange) ReadUint8(position int) (uint8, error) {
	err := e.beforeReadChecks(position)
	if err != nil {
		return 0, err
	}
	return e.answers[position].ContentAsUint8()
}

// Returns [Answer] content at the given position as uint16.
func (e *Exchange) ReadUint16(position int) (uint16, error) {
	err := e.beforeReadChecks(position)
	if err != nil {
		return 0, err
	}
	return e.answers[position].ContentAsUint16()
}

// Returns [Answer] content at the given position as uint32.
func (e *Exchange) ReadUint32(position int) (uint32, error) {
	err := e.beforeReadChecks(position)
	if err != nil {
		return 0, err
	}
	return e.answers[position].ContentAsUint32()
}

// Returns [Answer] content at the given position as uint64.
func (e *Exchange) ReadUint64(position int) (uint64, error) {
	err := e.beforeReadChecks(position)
	if err != nil {
		return 0, err
	}
	return e.answers[position].ContentAsUint64()
}

// Returns [Answer] content at the given position as string.
//
// TODO: test behavior for empty [Answer] and add checks if needed.
func (e *Exchange) ReadString(position int) (string, error) {
	err := e.beforeReadChecks(position)
	if err != nil {
		return "", err
	}
	return e.answers[position].ContentAsString()
}

func (e *Exchange) beforeReadChecks(position int) error {
	if !e.closed {
		return errors.New("can't read from open exchange")
	}
	if position < 0 || position >= len(e.answers) {
		return errors.New("can't read from position out of bounds")
	}
	return nil
}
