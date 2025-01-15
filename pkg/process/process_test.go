package process

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-github/v68/github"
)

func TestFile(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		content   string
		wantErr   bool
		setupMock func(client *github.Client)
	}{
		{
			name:     "valid file",
			filePath: "testdata/valid_file.yaml",
			content:  "uses: owner/repo@v1.0.0",
			wantErr:  false,
			setupMock: func(client *github.Client) {
				// Setup mock for getCommitSHA
			},
		},
		{
			name:     "invalid file path",
			filePath: "testdata/invalid_file.yaml",
			content:  "",
			wantErr:  true,
			setupMock: func(client *github.Client) {
				// No mock setup needed
			},
		},
		{
			name:     "invalid action reference",
			filePath: "testdata/invalid_action_ref.yaml",
			content:  "uses: owner/repo@invalid",
			wantErr:  true,
			setupMock: func(client *github.Client) {
				// Setup mock for getCommitSHA to return error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := github.NewClient(nil)
			tt.setupMock(client)

			// Write test content to file
			if err := os.WriteFile(tt.filePath, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			err := File(context.Background(), client, tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("File() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Clean up test file
			os.Remove(tt.filePath)
		})
	}
}
