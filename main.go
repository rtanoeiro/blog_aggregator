package main

import (
	"fmt"
	"log"
	"os"
)

const configFile = "/.gatorconfig.json"

func main() {
	state := &State{}
	commands := getCommands()
	appConfig, _ := ReadConfigFile(configFile)
	state.state = &appConfig
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	command := Command{Name: os.Args[1], Args: os.Args[2:]}
	fmt.Println("Provided command:", command)
	errorCommand := commands.run(state, command)

	if errorCommand != nil {
		log.Fatal(errorCommand)
	}

}
