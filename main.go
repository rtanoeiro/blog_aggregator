package main

import (
	"blog_aggregator/internal/database"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const configFile = "/.gatorconfig.json"

func main() {
	state := &State{}
	commands := getCommands()
	appConfig, _ := ReadConfigFile(configFile)
	state.config = &appConfig

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	db, err := sql.Open("postgres", appConfig.DBURL)
	if err != nil {
		log.Fatal("Unable to connect to Database")
	}
	dbQueries := database.New(db)
	state.db = dbQueries

	command := Command{Name: os.Args[1], Args: os.Args[2:]}
	fmt.Println("Provided command:", command)
	errorCommand := commands.run(state, command)

	if errorCommand != nil {
		log.Fatal(errorCommand)
	}
}
