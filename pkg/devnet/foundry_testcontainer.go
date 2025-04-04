package devnet

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Only run on tests for paralell anvil

type foundryContainer struct {
	testcontainers.Container
	URI string
}

const image = "ghcr.io/foundry-rs/foundry:v1.0.0"

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

func headerHealthCheck() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}

// const startuptimeout = 5 * time.Minute
// const pollInterval = 30 * time.Second

func SetupFoundry(ctx context.Context) (*foundryContainer, error) {
	bodyRequest, err := bodyHealthCheck()
	if err != nil {
		return nil, err
	}
	headers := headerHealthCheck()

	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"8545/tcp"},
		Cmd:          []string{"anvil --host 0.0.0.0 --port 8545"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Listening on 0.0.0.0:8545"),
			wait.ForExposedPort(),
			wait.ForHTTP("/").WithHeaders(headers).WithMethod(http.MethodPost).WithBody(bodyRequest),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	var foundryC *foundryContainer
	if container != nil {
		foundryC = &foundryContainer{Container: container}
	}
	if err != nil {
		return foundryC, err
	}
	// ip, err := container.Host(ctx)
	// if err != nil {
	// 	return foundryC, err
	// }
	// port, err := container.MappedPort(ctx, "8545")
	// if err != nil {
	// 	return foundryC, err
	// }
	// foundryC.URI = fmt.Sprintf("http://%s:%s", ip, port.Port())
	return foundryC, nil
}
