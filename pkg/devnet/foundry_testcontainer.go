package devnet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "embed"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

//go:embed anvil_state.json
var stateDefaultContent string

type FoundryOptions struct {
	stateContent string
	statePath    string
}

var _ testcontainers.ContainerCustomizer = (FoundryOption)(nil)

type FoundryOption func(*FoundryOptions) error

// Customize implements testcontainers.ContainerCustomizer.
func (f FoundryOption) Customize(req *testcontainers.GenericContainerRequest) error {
	return nil
}

func defaultOptions() FoundryOptions {
	return FoundryOptions{
		stateContent: stateDefaultContent,
	}
}

type FoundryContainer struct {
	testcontainers.Container
	stateContent string
}

type RPCRequest struct {
	Version string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
}

func bodyHealthCheck() (io.Reader, error) {
	body := RPCRequest{
		Version: "2.0",
		Method:  "web3_clientVersion",
		Params:  []string{},
		Id:      1,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bodyBytes), nil
}

func SetupFoundryV1(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*FoundryContainer, error) {
	return RunFoundry(ctx, "ghcr.io/foundry-rs/foundry:v1.0.0", opts...)
}

func SetupFoundryStable(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*FoundryContainer, error) {
	return RunFoundry(ctx, "ghcr.io/foundry-rs/foundry:stable", opts...)
}

func SetupFoundryNightly(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*FoundryContainer, error) {
	return RunFoundry(ctx, "ghcr.io/foundry-rs/foundry:nightly-2cdbfaca634b284084d0f86357623aef7a0d2ce3", opts...)
}

func WithFiles(content string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		files := []testcontainers.ContainerFile{
			{
				Reader:            bytes.NewBufferString(content),
				FileMode:          644,
				ContainerFilePath: "/usr/src/app/anvil_state.json",
			},
		}
		req.Files = files
		// req.Cmd = append(req.Cmd, "--load-state=/usr/src/app/anvil_state.json")
		return nil
	}
}

func WithStateContent(content []byte) FoundryOption {
	// create a temp json file
	stateParentDir, err := os.MkdirTemp("", "foundry-*")
	if err != nil {
		return func(o *FoundryOptions) error {
			return fmt.Errorf("failed to create temp dir: %v", err)
		}
	}
	statePath := filepath.Join(stateParentDir, "anvil_state.json")
	stateFile, err := os.Create(statePath)
	if err != nil {
		return func(o *FoundryOptions) error {
			return fmt.Errorf("failed to create temp file: %v", err)
		}
	}
	defer stateFile.Close()
	_, err = stateFile.Write([]byte(content))
	if err != nil {
		return func(o *FoundryOptions) error {
			return fmt.Errorf("failed to write to temp file: %v", err)
		}
	}
	return func(o *FoundryOptions) error {
		o.statePath = statePath
		return nil
	}
}

func RunFoundry(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (*FoundryContainer, error) {
	var err error
	// health check
	bodyRequest, err := bodyHealthCheck()
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        img,
			ExposedPorts: []string{"8546/tcp"},
			Cmd: []string{"anvil --host=0.0.0.0 --port=8546"},
			WaitingFor: wait.ForAll(
				wait.ForExposedPort(),
				wait.ForLog("Listening on 0.0.0.0"),
				wait.ForHTTP("/").WithHeaders(headers).WithMethod(http.MethodPost).WithBody(bodyRequest),
			).WithDeadline(1 * time.Minute).WithStartupTimeoutDefault(30 * time.Second),
		},
		Started: true,
	}

	// Gather all config options (defaults and then apply provided options)
	settings := defaultOptions()
	for _, opt := range opts {
		if apply, ok := opt.(FoundryOption); ok {
			if err := apply(&settings); err != nil {
				return nil, err
			}
		}
		if err := opt.Customize(&req); err != nil {
			return nil, err
		}
	}

	// Add volume with anvil state file
	err = WithFiles(settings.stateContent).Customize(&req)
	if err != nil {
		return nil, err
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	var foundryC *FoundryContainer
	if container != nil {
		foundryC = &FoundryContainer{
			Container:    container,
			stateContent: settings.stateContent,
		}
	}
	return foundryC, err
}

func (f *FoundryContainer) URI() (*string, error) {
	ctx := context.Background()
	container := f.Container
	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := container.MappedPort(ctx, "8545")
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("http://%s:%s", ip, port.Port())
	return &uri, nil
}
