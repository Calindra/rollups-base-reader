package repository

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AppRepositorySuite struct {
	suite.Suite
}

func TestAppRepository(t *testing.T) {
	suite.Run(t, new(AppRepositorySuite))
}