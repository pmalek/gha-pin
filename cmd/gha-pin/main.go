package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/pmalek/gha-pin/pkg/run"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s workflow-file [workflow-file...]\n\n"+
			"Environment variables:\n"+
			"  GITHUB_TOKEN    GitHub personal access token (optional)\n",
			os.Args[0])
	}

	var files []string
	for _, arg := range os.Args[1:] {
		matches, err := filepath.Glob(arg)
		if err != nil {
			log.Fatalf("Invalid glob pattern %q: %v", arg, err)
		}
		if len(matches) == 0 {
			log.Printf("Warning: no files match pattern %q", arg)
			continue
		}
		files = append(files, matches...)
	}

	if len(files) == 0 {
		log.Fatal("No files to process")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Cleaning up...")
		cancel()
	}()

	if err := run.Do(ctx, files); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully processed %d file(s)\n", len(files))
}
