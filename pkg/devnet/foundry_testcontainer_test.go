package devnet

import (
	"context"
	"testing"

	"github.com/calindra/rollups-base-reader/pkg/commons"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type FoundryTestContainerSuite struct {
	suite.Suite
	ctx       context.Context
	ctxCancel context.CancelFunc
	container testcontainers.Container
}

// const timeout = 5 * time.Minute

func TestFoundryTestContainerSuite(t *testing.T) {
	suite.Run(t, new(FoundryTestContainerSuite))
}

func (s *FoundryTestContainerSuite) SetupTest() {
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	// s.ctx, s.ctxCancel = context.WithTimeout(context.Background(), timeout)
}

func (s *FoundryTestContainerSuite) TearDownTest() {
	testcontainers.CleanupContainer(s.T(), s.container)
	s.ctxCancel()
}
func (s *FoundryTestContainerSuite) TestFoundryContainer() {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	container, err := SetupFoundryNightly(ctx, testcontainers.WithLogConsumers(&commons.StdoutLogConsumer{}))
	s.container = container
	s.Require().NoError(err)
	s.NotNil(container)
	// s.NotEmpty(container.URI)
}
