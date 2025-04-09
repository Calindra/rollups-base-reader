package commons

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DATestSuite struct {
	suite.Suite
}

func TestDATest(t *testing.T) {
	suite.Run(t, new(DATestSuite))
}

// func (s *DATestSuite) TestConvertDA() {
// 	// Test the convertDA function
// 	inputBoxAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
// 	expectedOutput := "0x1234567890abcdef1234567890abcdef12345678"

// 	output, err := convertDA(inputBoxAddress)
// 	s.NoError(err)
// 	s.Equal(expectedOutput, *output)
// }
