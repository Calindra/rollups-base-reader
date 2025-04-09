package transaction

import "github.com/cartesi/rollups-graphql/pkg/commons"

type Sender interface {
	SubmitSigAndData(sigAndData commons.SigAndData) (string, error)
}
