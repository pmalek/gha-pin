package process

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v68/github"
)

// File processes the content of the given file.
func File(ctx context.Context, client *github.Client, filePath string) error {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	// Process the content
	updatedContent, err := processFileContent(ctx, client, string(content))
	if err != nil {
		return fmt.Errorf("processing content: %w", err)
	}

	// Write the updated content back to the file
	if err := os.WriteFile(filePath, []byte(updatedContent), 0o644); err != nil {
		return fmt.Errorf("saving file: %w", err)
	}

	fmt.Printf("Successfully processed %s\n", filePath)
	return nil
}

func getCommitSHA(ctx context.Context, client *github.Client, owner, repo, ref string) (string, error) {
	commit, _, err := client.Repositories.GetCommitSHA1(ctx, owner, repo, ref, "")
	if err != nil {
		return "", fmt.Errorf("getting commit SHA: %w", err)
	}
	return commit, nil
}

func parseActionReference(actionRef string) (owner, repo, ref string) {
	parts := strings.Split(actionRef, "@")
	if len(parts) != 2 {
		return "", "", ""
	}
	repoPath := strings.Split(parts[0], "/")
	if len(repoPath) != 2 {
		return "", "", ""
	}
	return repoPath[0], repoPath[1], parts[1]
}

func processFileContent(ctx context.Context, client *github.Client, content string) (string, error) {
	// Regular expression to find "uses: owner/repo@ref"
	re := regexp.MustCompile(`(uses:\s*([^@]+)@([^#\s]+))`)

	// Function to replace matches with updated SHA references and add comments
	result := re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract details from the match
		parts := re.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match // Skip if format is unexpected
		}

		_, repoPath, ref := parts[1], parts[2], parts[3]
		owner, repo, version := parseActionReference(fmt.Sprintf("%s@%s", repoPath, ref))
		if owner == "" || repo == "" || version == "" || len(version) == 40 {
			return match // Skip invalid or already SHA-based references
		}

		sha, err := getCommitSHA(ctx, client, owner, repo, version)
		if err != nil {
			log.Printf("Failed to get commit SHA for %s/%s@%s: %v", owner, repo, version, err)
			return match
		}

		// Add a comment indicating the replaced version at the end of the line
		return fmt.Sprintf("uses: %s/%s@%s # %s", owner, repo, sha, version)
	})

	return result, nil
}
