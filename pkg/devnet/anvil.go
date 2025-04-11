// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package devnet

import (
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/calindra/rollups-base-reader/pkg/supervisor"
)

// Default port for the Ethereum node.
const (
	AnvilDefaultAddress = "127.0.0.1"
	AnvilDefaultPort    = 8545
)

// Generate the devnet state and embed it in the Go binary.
//
//go:generate go run ./gen-devnet-state
//go:embed anvil_state.json
var devnetState []byte

//go:embed deployments/*.json
var bundleContracts embed.FS

const stateFileName = "anvil_state.json"

const anvilCommand = "anvil"

// Start the anvil process in the host machine.
type AnvilWorker struct {
	Address        string
	Port           int
	Verbose        bool
	AnvilCmd       string
	AnvilBlockTime time.Duration
}

// Define a struct to represent the structure of your JSON data
type ContractInfo struct {
	ContractName string `json:"contractName"`
	Address      string `json:"address"`
}

func (w AnvilWorker) String() string {
	return anvilCommand
}

// Try to install anvil
func InstallAnvil(ctx context.Context) (string, error) {
	anvilClient := NewAnvilRelease()
	anvilExec, err := HandleReleaseExecution(ctx, anvilClient)
	if err != nil {
		return "", fmt.Errorf("anvil: %w", err)
	}
	return anvilExec, nil
}

func GetAnvilVersion(anvilExec string) (string, error) {
	output, err := exec.Command(anvilExec, "--version").Output()
	if err != nil {
		return "", fmt.Errorf("anvil: failed to run anvil: %w", err)
	}
	anvilVersion := strings.TrimSpace(string(output))
	return anvilVersion, nil
}

// For while, we dont check if anvil is already installed.
// We always install anvil.
func CheckAnvilAndInstall(ctx context.Context) (string, error) {
	location, err := InstallAnvil(ctx)
	if err != nil {
		return "", fmt.Errorf("anvil: failed to install %w", err)
	}
	anvilVersion, err := GetAnvilVersion(location)
	if err != nil {
		return "", err
	}
	slog.Debug("anvil: installed anvil", "location", location, "version", anvilVersion)
	return location, nil
}

func GetContractInfo() ([]ContractInfo, error) {
	contracts := []ContractInfo{}

	// Read the JSON file
	err := fs.WalkDir(bundleContracts, ".", func(path string, d fs.DirEntry, err error) error {
		var contract ContractInfo

		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}

		content, err := fs.ReadFile(bundleContracts, path)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(content, &contract); err != nil {
			return fmt.Errorf("anvil: failed to unmarshal %s", path)
		}

		// Add the contract to the list if it is a Cartesi Rollups contract
		if strings.HasPrefix(filepath.Base(path), "cartesiRollups") {
			contracts = append(contracts, contract)
		}

		return nil
	})

	return contracts, err
}

func writeLine(writer io.WriteCloser, firstColumn, SecondColumn string) error {
	_, err := writer.Write(fmt.Appendf(nil, "%-28s %s\n", firstColumn, SecondColumn))
	if err != nil {
		return fmt.Errorf("anvil: failed to write line: %w", err)
	}
	return nil
}

const nameSpace = 28
const addressSpace = 42

func writeHeaderLine(writer io.WriteCloser) error {
	return writeLine(writer, strings.Repeat("─", nameSpace), strings.Repeat("─", addressSpace))
}

func ShowAddresses(writer io.WriteCloser) error {
	defer writer.Close()
	contracts, err := GetContractInfo()
	if err != nil {
		return err
	}
	contracts = append(contracts, ContractInfo{
		ContractName: "Application",
		Address:      ApplicationAddress,
	})

	// Sort things
	sort.Slice(contracts, func(i, j int) bool {
		return contracts[i].ContractName < contracts[j].ContractName
	})

	if err := writeLine(writer, "Contract", "Address"); err != nil {
		return err
	}

	if err := writeHeaderLine(writer); err != nil {
		return err
	}

	for _, contract := range contracts {
		if err := writeLine(writer, contract.ContractName, contract.Address); err != nil {
			return err
		}
	}

	return nil
}

func (w AnvilWorker) Start(ctx context.Context, ready chan<- struct{}) error {
	dir, err := makeStateTemp()
	if err != nil {
		return err
	}
	defer removeTemp(dir)
	slog.Debug("anvil: created temp dir with state file", "dir", dir)

	if w.AnvilCmd == "" {
		w.AnvilCmd = anvilCommand
	}

	seconds := uint64(w.AnvilBlockTime.Seconds())

	var server supervisor.ServerWorker
	server.Name = anvilCommand
	server.Command = w.AnvilCmd
	server.Port = w.Port
	server.Args = append(server.Args, "--host", fmt.Sprint(w.Address))
	server.Args = append(server.Args, "--port", fmt.Sprint(w.Port))
	server.Args = append(server.Args, "--load-state", filepath.Join(dir, stateFileName))
	if seconds > 0 {
		server.Args = append(server.Args, "--block-time", fmt.Sprint(seconds))
	}
	// server.Args = append(server.Args, "--tracing")
	if !w.Verbose {
		server.Args = append(server.Args, "--silent")
	}
	return server.Start(ctx, ready)
}

// Create a temporary directory with the state file in it.
// The directory should be removed by the callee.
func makeStateTemp() (string, error) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("anvil: failed to create temp dir: %w", err)
	}
	stateFile := filepath.Join(tempDir, stateFileName)
	const permissions = 0644
	err = os.WriteFile(stateFile, devnetState, permissions)
	if err != nil {
		return "", fmt.Errorf("anvil: failed to write state file: %w", err)
	}
	return tempDir, nil
}

// Delete the temporary directory.
func removeTemp(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		slog.Warn("anvil: failed to remove temp file", "error", err)
	}
}
