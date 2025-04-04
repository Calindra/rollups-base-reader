package paiodecoder

import "github.com/ethereum/go-ethereum/common"

type AvailData struct {
	Transaction string         `json:"transaction"`
	Data        []byte         `json:"data"`
	AppContract common.Address `json:"app_contract"`
	MsgSender   common.Address `json:"msg_sender"`
}
