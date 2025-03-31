package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/calindra/rollups-base-reader/pkg/commons"
)

// Using from https://github.com/cartesi/rollups-node/blob/c748b205fe35f217868d98e8e64909afdd11d9af/internal/model/models.go

const (
	// URL of the models.go file in GitHub
	// modelsURL = "https://raw.githubusercontent.com/cartesi/rollups-node/next/2.0/internal/model/models.go"
	modelsURL = "https://raw.githubusercontent.com/cartesi/rollups-node/d8063929ef309df6925ff45883cce482ec3d9e18/internal/model/models.go"
	// Local destination path
	localPath = "node_models.go"
)

const timeout time.Duration = 5 * time.Minute

func main() {
	commons.ConfigureLog(slog.LevelDebug)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create the directory path if it doesn't exist
	dir := filepath.Dir(localPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %v", err)
		}
	}

	// Download the file
	slog.Info("Downloading models", "url", modelsURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, modelsURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to download file, status code: %d", resp.StatusCode)
	}

	// Create the destination file
	out, err := os.Create(localPath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer out.Close()

	// Modify the package declaration
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Replace the original package declaration with our package name
	modified := replacePackage(string(content), "model", "models")

	// Write the modified content to the file
	_, err = out.WriteString(modified)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	slog.Info("Successfully downloaded and saved models", "path", localPath)
}

// replacePackage replaces the package declaration in the given content
func replacePackage(content, oldPkg, newPkg string) string {
	packageLine := fmt.Sprintf("package %s", oldPkg)
	newPackageLine := fmt.Sprintf("package %s", newPkg)

	// Add a comment to indicate this is a generated file
	header := "// Code generated from GitHub. DO NOT EDIT.\n// Source: " + modelsURL + "\n\n"

	// Replace the package line and add the header
	modified := content
	if len(content) >= len(packageLine) && content[:len(packageLine)] == packageLine {
		modified = newPackageLine + content[len(packageLine):]
	}

	return header + modified
}
