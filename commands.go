package main

import (
	"blog_aggregator/internal/database"
	"context"
	"database/sql"
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
			"addfeed": func(state *State, command Command) error {
				return isLogged(state, command, handleAddFeed)
			},
			"feeds": handleGetFeeds,
			"follow": func(state *State, command Command) error {
				return isLogged(state, command, handleFollowFeed)
			},
			"following": func(state *State, command Command) error {
				return isLogged(state, command, handleGetFollowing)
			},
			"unfollow": func(state *State, command Command) error {
				return isLogged(state, command, handleUnfollowFeed)
			},
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

func isLogged(state *State, command Command, handler func(state *State, command Command) error) error {
	_, logged := state.db.GetUser(context.Background(), state.config.CurrentUserName)

	if logged != nil {
		return errors.New("user not logged in, to run this command, please log in first. go run . login <username>")
	}

	return handler(state, command)
}

func handlerLogin(state *State, command Command) error {

	if len(command.Args) == 0 {
		return errors.New("please provide an username to use the register command")

	}
	username := command.Args[0]
	state.config.SetUser(username)

	fmt.Println("-User:", username)
	return nil
}

func handlerRegister(state *State, command Command) error {
	fmt.Println("Entering Register State")

	if len(command.Args) == 0 {
		return errors.New("please provide an username to use the register command")

	}
	username := command.Args[0]
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
	fmt.Println("User registered with success. \n-ID:", user.ID, "\n-Name:", user.Name, "\n-CreatedAt:", user.CreatedAt, "\n-UpdatedAt:", user.UpdatedAt)
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
	users, err := state.db.GetAllUsers(context.Background())

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

	if len(command.Args) != 1 {
		return errors.New("error getting feed data")
	}

	duration := command.Args[0]
	fmt.Println("Collecting feed every ", duration)
	timeDuration, _ := time.ParseDuration(duration)
	ticker := time.NewTicker(timeDuration)
	for ; ; <-ticker.C {
		url, nextError := state.db.GetNextFeedToFetch(context.Background())
		feedRow, _ := state.db.GetFeedFromURL(context.Background(), url)

		if nextError != nil {
			return errors.New("unable to get the next URL to parse")
		}

		ScrapeFeed(state, url, feedRow.FeedID.String())
		fmt.Print("Finshed this round, trying again in the predetermined time\n\n")
		_ = UpdateLastFetchedAt(state, feedRow.FeedID.String())

	}
}

func ScrapeFeed(state *State, url string, feedID string) {
	results, feedErr := fetchFeed(context.Background(), url)

	if feedErr != nil {
		return
	}
	results.CleanFeed()

	for _, item := range results.Channel.Item {
		exists := checkIfPostExists(state, item.Description, item.Link)

		if exists {
			continue
		}

		_, postError := addToPosts(state, item, item.Link, feedID)
		if postError != nil {
			fmt.Println("failed to add to posts", postError)
		}
	}
}

func addToPosts(state *State, rssItem RSSItem, url string, feedID string) (database.Post, error) {
	params := database.InsertPostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       rssItem.Title,
		Url:         url,
		Description: sql.NullString{String: rssItem.Description, Valid: true},
		PublishedAt: sql.NullString{String: rssItem.PubDate, Valid: true},
		FeedID:      uuid.MustParse(feedID),
	}
	post, postError := state.db.InsertPost(context.Background(), params)

	return post, postError

}

func checkIfPostExists(state *State, description string, url string) bool {
	params := database.CheckIfPostExistsParams{
		Description: sql.NullString{String: description, Valid: true},
		Url:         url,
	}
	_, exists := state.db.CheckIfPostExists(context.Background(), params)

	return exists == nil
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
		ID:            uuid.New(),
		Name:          command.Args[0],
		Url:           command.Args[1],
		UserID:        userInfo.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: false},
	}
	results, insertFeedError := state.db.InsertFeed(context.Background(), arguments)

	if insertFeedError != nil {
		return errors.New("unable to add follow feed into table")
	}
	fmt.Println("Insert feed with success. \nID:", results.ID, "\n- Name:", results.Name, "\n- URL:", results.Url, "\n- UserID:", results.UserID)
	followError := handleFollowFeed(state, Command{
		Name: "follow",
		Args: []string{results.Url},
	},
	)

	if followError != nil {
		return errors.New("unable to follow feed after adding it")
	}
	return nil
}

func handleGetFeeds(state *State, command Command) error {
	feedRows, feedError := state.db.GetAllFeeds(context.Background())

	if feedError != nil {
		return errors.New("failed to get all feeds")
	}

	for _, feed := range feedRows {
		fmt.Println("Feed Name:", feed.Name, "\n- URL:", feed.Url, "\n- Username:", feed.Username)
	}
	return nil
}

func handleFollowFeed(state *State, command Command) error {

	if len(command.Args) != 1 {
		return errors.New("when following a feed, provide the feed link, go run . follow <feed_url>")
	}
	feedURL := command.Args[0]
	feedRow, feedError := state.db.GetFeedFromURL(context.Background(), feedURL)
	userInfo, userError := state.db.GetUser(context.Background(), state.config.GetCurrentUser())

	if feedError != nil {
		return errors.New("failed to get feed data or feed doesn't exist, check the full list with go run . feeds")
	}
	if userError != nil {
		return errors.New("unable to follow feed, user is not registered. register one with go run . / <username>")
	}
	insertFollowFeed := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userInfo.ID,
		FeedID:    feedRow.FeedID,
	}
	results, insertFeedError := state.db.CreateFeedFollow(context.Background(), insertFollowFeed)

	if insertFeedError != nil {
		return errors.New("unable to add follow feed into table")
	}

	fmt.Println("Follow feed insert success! \n-ID", results.ID, "\n- CreatedAt:", results.CreatedAt, "\n- UpdatedAt:", results.UpdatedAt, "\n- UserID:", results.UserID, "\n- FeedID:", results.FeedID)
	return nil
}

func handleGetFollowing(state *State, command Command) error {
	userInfo, userError := state.db.GetUser(context.Background(), state.config.GetCurrentUser())

	if userError != nil {
		return errors.New("unable to follow feed, user is not registered. register one with go run . / <username>")
	}

	followingRows, followingError := state.db.GetFollowedFeedsFromUser(context.Background(), userInfo.ID)

	if followingError != nil {
		return errors.New("failed to get all following feeds")
	}

	for _, following := range followingRows {
		fmt.Println("Feed Name:", following.Name)
	}
	return nil
}

func handleUnfollowFeed(state *State, command Command) error {

	if len(command.Args) != 1 {
		return errors.New("when unfollowing a feed, provide the feed link, go run . unfollow <feed_url>")
	}

	userInfo, userError := state.db.GetUser(context.Background(), state.config.GetCurrentUser())

	if userError != nil {
		return errors.New("unable to follow feed, user is not registered. register one with go run . / <username>")
	}
	arguments := database.UnfollowParams{
		Url:    command.Args[0],
		UserID: userInfo.ID,
	}
	results, unfollowError := state.db.Unfollow(context.Background(), arguments)

	if unfollowError != nil {
		return errors.New("unable to unfollow feed")
	}

	fmt.Println("Unfollow success!\n-User:", results.UserID, "\n-FeedID:", results.FeedID)
	return nil
}

func UpdateLastFetchedAt(state *State, feed_id string) error {
	params := database.UpdateFetchDateParams{
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:            uuid.MustParse(feed_id),
	}
	_, updateFetchErr := state.db.UpdateFetchDate(context.Background(), params)

	if updateFetchErr != nil {
		return errors.New("failed to update last fetched at date")
	}
	return nil
}
