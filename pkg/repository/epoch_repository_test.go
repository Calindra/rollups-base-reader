package repository

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type EpochRepositorySuite struct {
	suite.Suite
}

func TestEpochRepository(t *testing.T) {
	suite.Run(t, new(EpochRepositorySuite))
}