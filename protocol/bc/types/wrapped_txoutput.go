package types

import (
	"io"

	"github.com/bytom/encoding/blockchain"
)

type WrappedTxOutput interface {
	readFrom(*blockchain.Reader) (err error)
	writeTo(io.Writer) error
	writeCommitment(io.Writer) error
}
