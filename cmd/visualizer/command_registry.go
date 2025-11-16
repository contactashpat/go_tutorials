package main

import "fmt"

// Command represents something that can run with CLI args.
type Command interface {
	Run(args []string) error
}

// CommandFunc allows plain functions to satisfy Command.
type CommandFunc func(args []string) error

func (f CommandFunc) Run(args []string) error { return f(args) }

var commandFactories = map[string]func() Command{
	"see":    func() Command { return NewSeeCommand() },
	"decode": func() Command { return NewDecodeCommand() },
}

func registerCommand(name string, factory func() Command) {
	if _, exists := commandFactories[name]; exists {
		panic(fmt.Sprintf("command %s already registered", name))
	}
	commandFactories[name] = factory
}

func dispatchCommand(name string, args []string) error {
	factory, exists := commandFactories[name]
	if !exists {
		return fmt.Errorf("unknown command %s", name)
	}
	return factory().Run(args)
}
