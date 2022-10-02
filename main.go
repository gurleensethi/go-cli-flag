// Available Commands:
// - search-repos: Search for github repos
// - search-users: Serach for users on github.
//
// Flags:
// - Top level flags:
//   - debug: Print the debug information as executing command
//
// Example:
// - go run main.go -debug search-repos golang
// - go run main.go -debug search-users gurleensethi

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	debug = flag.Bool("debug", false, "log out all the debug information")

	usage = `Specify a command to execute:
  - search-repos: Search for github repos
  - search-users: Serach for users on github.`
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	command := flag.Args()[0]

	err := executeCommand(command, flag.Args()[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func executeCommand(command string, args []string) error {
	printDebug(fmt.Sprintf("Command: %s", command))
	printDebug(fmt.Sprintf("Args: %v", args))

	switch command {
	case "search-repos":
		return executeSearchRepos(args)
	case "search-users":
		return executeSearchUsers(args)
	default:
		return fmt.Errorf("invalid command: '%s'", command)
	}
}

func executeSearchRepos(args []string) error {
	flagSet := flag.NewFlagSet("search-repos", flag.ExitOnError)
	flagSet.Parse(args)

	printDebug(fmt.Sprintf("[search-repos] Args: %s", flagSet.Args()))

	if len(flagSet.Args()) == 0 {
		return errors.New("provide a search term for searching repos: search-repos <search_term>")
	}

	searchTerm := flagSet.Args()[0]

	printDebug(fmt.Sprintf("[search-repos] Search Term: %s", searchTerm))

	repos, err := findRepos(searchTerm)
	if err != nil {
		return err
	}

	fmt.Println(strings.Join(repos, ", "))

	return nil
}

func executeSearchUsers(args []string) error {
	flagSet := flag.NewFlagSet("search-repos", flag.ExitOnError)
	flagSet.Parse(args)

	printDebug(fmt.Sprintf("[search-repos] Args: %s", flagSet.Args()))

	if len(flagSet.Args()) == 0 {
		return errors.New("provide a search term for searching repos: search-repos <search_term>")
	}

	searchTerm := flagSet.Args()[0]

	printDebug(fmt.Sprintf("[search-repos] Search Term: %s", searchTerm))

	users, err := findUsers(searchTerm)
	if err != nil {
		return err
	}

	fmt.Println(strings.Join(users, ", "))

	return nil
}

func findRepos(term string) ([]string, error) {
	type repo struct {
		FullName string `json:"full_Name"`
	}

	type searchResult struct {
		Items []repo `json:"items"`
	}

	// Prepare github repository search url.
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/search/repositories", nil)
	if err != nil {
		printDebug(fmt.Sprintf("%v", err))
		return nil, errors.New("failed to connect to github")
	}

	query := req.URL.Query()
	query.Set("q", term)
	req.URL.RawQuery = query.Encode()

	// Make http request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		printDebug(fmt.Sprintf("%v", err))
		return nil, errors.New("failed to connect to github")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("failed to connect to github")
	}

	// Parse the json response.
	results := searchResult{}

	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		printDebug(fmt.Sprintf("%v", err))
		return nil, errors.New("failed to connect to github")
	}

	// Extract out the repo names.
	repos := make([]string, 0)

	for _, r := range results.Items {
		repos = append(repos, r.FullName)
	}

	return repos, nil
}

func findUsers(term string) ([]string, error) {
	type user struct {
		Login string `json:"login"`
	}

	type searchResult struct {
		Items []user `json:"items"`
	}

	// Prepare github repository search url.
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/search/users", nil)
	if err != nil {
		printDebug(fmt.Sprintf("%v", err))
		return nil, errors.New("failed to connect to github")
	}

	query := req.URL.Query()
	query.Set("q", term)
	req.URL.RawQuery = query.Encode()

	// Make http request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		printDebug(fmt.Sprintf("%v", err))
		return nil, errors.New("failed to connect to github")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("failed to connect to github")
	}

	// Parse the json response.
	results := searchResult{}

	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		printDebug(fmt.Sprintf("%v", err))
		return nil, errors.New("failed to connect to github")
	}

	// Extract out the repo names.
	repos := make([]string, 0)

	for _, r := range results.Items {
		repos = append(repos, r.Login)
	}

	return repos, nil
}

func printDebug(msg string) {
	if *debug {
		fmt.Printf("[DEBUG]: %s\n", msg)
	}
}
