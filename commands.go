package main

import (
	"errors"
	"fmt"
)

type State struct {
	state *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	handlers map[string]func(state *State, command Command) error
}

func getCommands() Commands {
	commands := Commands{
		handlers: map[string]func(state *State, command Command) error{
			"login":    handlerLogin,
			"register": handlerRegister,
		},
	}
	return commands
}

func (c *Commands) run(state *State, command Command) error {
	function, ok := c.handlers[command.Name]
	if !ok {
		return errors.New("command not found")
	}
	return function(state, command)
}

func handlerLogin(state *State, command Command) error {

	if len(command.Args) == 0 {
		return errors.New("username is required")
	}

	user := command.Args[0]
	loginError := state.state.GetUserConfig(user)

	if loginError != nil {
		return errors.New("user not present, unable to login")
	}
	return nil
}

func handlerRegister(state *State, command Command) error {
	fmt.Println("Entering Register State")
	user := command.Args[0]
	addError := state.state.AddUser(user)

	if addError != nil {
		return errors.New("user already registered")
	}
	return nil
}
