// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This program gets the devnet state from the devnet Docker image.
// To do that, it creates a container from the image, copies the state file, and deletes the
// container.
package main

import (
	"bufio"
	"log/slog"
	"os/exec"

	"github.com/cartesi/rollups-graphql/v2/pkg/commons"
)

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// Read the output from stdout and stderr
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			slog.Debug("[DOCKER]", "line", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			slog.Debug("[DOCKER]", "line", scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		slog.Error("Command failed", "error", err)
	}

	slog.Debug("Command finished", "name", name, "args", args)
}

// containerExists checks if a Docker container exists
func containerExists(containerName string) bool {
	cmd := exec.Command("docker", "container", "inspect", containerName)
	err := cmd.Run()
	// If the command returns an error, the container doesn't exist
	return err == nil
}

func main() {
	commons.ConfigureLog(slog.LevelDebug)

	// Remove previous container only if it exists
	if containerExists("temp-devnet") {
		slog.Info("Removing existing temp-devnet container")
		run("docker", "rm", "temp-devnet")
	}

	// you can see the tags on
	// https://github.com/cartesi/cli/pkgs/container/sdk
	// update me when the image is updated
	slog.Info("Creating temporary container")

	// like __filename
	// _, xdir, _, ok := runtime.Caller(0)
	// if !ok {
	// 	panic("no found dirname")
	// }

	// Custom Docker anvil state
	// dockerfile := filepath.Join(filepath.Dir(xdir), "Dockerfile.sdk")
	// run("docker", "build", "-t", "temp-devnet:devrel", ".", "-f", dockerfile)
	// run("docker", "create", "--name", "temp-devnet", "temp-devnet:devrel")

	// Production ready anvil state
	run("docker", "create", "--name", "temp-devnet", "ghcr.io/cartesi/sdk:0.12.0-alpha.14")

	slog.Info("Copying the state file")
	defer func() {
		run("docker", "rm", "temp-devnet")
		// run("docker", "rmi", "temp-devnet:devrel")
		slog.Info("Finished copying the state file")
	}()
	run("docker", "cp", "temp-devnet:/usr/share/cartesi/anvil_state.json", ".")
	run("docker", "cp", "temp-devnet:/usr/share/cartesi/deployments/", ".")
}
