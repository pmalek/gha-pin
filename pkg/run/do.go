package run

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"

	"github.com/pmalek/gha-pin/pkg/process"
)

func getGitHubClient(ctx context.Context) *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return github.NewClient(nil)
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// Do processes the given list of files in parallel.
func Do(ctx context.Context, files []string) error {
	client := getGitHubClient(ctx)

	// Create a channel for file paths
	fileChan := make(chan string, len(files))
	resultChan := make(chan error, len(files))

	// Worker function
	worker := func() {
		for file := range fileChan {
			select {
			case <-ctx.Done():
				resultChan <- ctx.Err()
				return
			default:
				// Process the file and send any errors to the result channel
				err := process.File(ctx, client, file)
				resultChan <- err
			}
		}
	}

	// Start workers
	for i := 0; i < len(files); i++ {
		go worker()
	}

	// Send file paths to the file channel
	go func() {
		for _, file := range files {
			fileChan <- file
		}
		close(fileChan)
	}()

	// Collect results
	var errJoin error
	for _, file := range files {
		if err := <-resultChan; err != nil {
			if errJoin == nil {
				errJoin = errors.Join(errJoin, err)
			}
			log.Printf("Error: %v", err)
			continue
		}
		log.Printf("Processed: %s", file)
	}

	// Return the first error if any occurred
	return errJoin
}
