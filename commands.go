package main

import (
	"blog_aggregator/internal/database"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type State struct {
	db     *database.Queries
	config *Config
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
			"reset":    handleReset,
			"users":    handleGetUsers,
			"agg":      handleAgg,
			"addfeed":  handleAddFeed,
			"feeds":    handleGetFeeds,
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
		return errors.New("please provide an username to use the register command")

	}
	username := command.Args[0]

	user, loginError := state.db.GetUser(context.Background(), username)

	if loginError != nil {
		return errors.New("user not present, unable to login")
	}
	state.config.SetUser(username)

	fmt.Println("Got user:", user.Name, "With ID:", user.ID)
	return nil
}

func handlerRegister(state *State, command Command) error {
	fmt.Println("Entering Register State")

	if len(command.Args) == 0 {
		return errors.New("please provide an username to use the register command")

	}
	username := command.Args[0]

	_, getError := state.db.GetUser(context.Background(), username)

	if getError == nil {
		return errors.New("user already exists, unable to register again")
	}

	state.config.SetUser(username)
	arguments := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	user, addError := state.db.CreateUser(context.Background(), arguments)

	if addError != nil {
		return errors.New("error registering user")
	}
	fmt.Println("User registered with success:", user)
	return nil
}

func handleReset(state *State, command Command) error {
	err := state.db.CleanUsers(context.Background())

	if err != nil {
		return errors.New("unable to truncate the users table")
	}
	return nil
}

func handleGetUsers(state *State, command Command) error {
	users, err := state.db.GetUsers(context.Background())

	if err != nil {
		return errors.New("unable to get users from database")
	}

	for _, user := range users {
		if user == state.config.GetCurrentUser() {
			fmt.Println("*", user, "(current)")
			continue
		}

		fmt.Println("*", user)
	}

	return nil
}

func handleAgg(state *State, command Command) error {
	results, feedErr := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if feedErr != nil {
		return errors.New("failed to get results")
	}
	results.CleanFeed()
	fmt.Println(results)

	return nil
}

func handleAddFeed(state *State, command Command) error {

	if len(command.Args) != 2 {
		return errors.New("when adding a new feed, provide the feed name and link, go run . addfeed <feed_name> <feed_url>")
	}

	currentUser := state.config.GetCurrentUser()
	userInfo, userError := state.db.GetUser(context.Background(), currentUser)

	if userError != nil {
		return errors.New("before adding a feed to an user, please add an user first go run . register <username>")
	}

	arguments := database.InsertFeedParams{
		Name:   command.Args[0],
		Url:    command.Args[1],
		UserID: userInfo.ID,
	}
	state.db.InsertFeed(context.Background(), arguments)
	return nil
}

func handleGetFeeds(state *State, command Command) error {
	feedRows, feedError := state.db.GetFeeds(context.Background())

	if feedError != nil {
		return errors.New("failed to get all feeds")
	}

	for _, feed := range feedRows {
		fmt.Println("Feed Name:", feed.Name, "- URL:", feed.Url, "- Username:", feed.Username)
	}
	return nil
}
