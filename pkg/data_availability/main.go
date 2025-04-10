package main

import (
	"fmt"
	"math/big"

	"github.com/calindra/rollups-base-reader/pkg/contracts"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	addressInputBox := common.HexToAddress("0xBa3Cf8fB82E43D370117A0b7296f91ED674E94e3")

	parsedAbi, err := contracts.DataAvailabilityMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	inputBoxDA, err := parsedAbi.Pack("InputBox",
		addressInputBox)
	if err != nil {
		panic(err)
	}

	espressoDA, err := parsedAbi.Pack("InputBoxAndEspresso",
		addressInputBox,
		big.NewInt(0),
		uint32(0))
	if err != nil {
		panic(err)
	}

	fmt.Printf("0x%x\n", inputBoxDA)
	fmt.Printf("0x%x\n", espressoDA)
}
