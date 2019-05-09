package types

import (
	"io"

	"github.com/bytom/encoding/blockchain"
	"github.com/bytom/errors"
	"github.com/bytom/protocol/bc"
)

// ClaimTxOutput is the top level struct of tx output.
type ClaimTxOutput struct {
	AssetVersion uint64
	OutputCommitment
	// Unconsumed suffixes of the commitment and witness extensible strings.
	CommitmentSuffix []byte
}

// NewClaimTxOutput create a new output struct
func NewClaimTxOutput(assetID bc.AssetID, amount uint64, controlProgram []byte) *ClaimTxOutput {
	return &ClaimTxOutput{
		AssetVersion: 1,
		OutputCommitment: OutputCommitment{
			AssetAmount: bc.AssetAmount{
				AssetId: &assetID,
				Amount:  amount,
			},
			VMVersion:      1,
			ControlProgram: controlProgram,
		},
	}
}

func (to *ClaimTxOutput) readFrom(r *blockchain.Reader) (err error) {
	if to.AssetVersion, err = blockchain.ReadVarint63(r); err != nil {
		return errors.Wrap(err, "reading asset version")
	}

	if to.CommitmentSuffix, err = to.OutputCommitment.readFrom(r, to.AssetVersion); err != nil {
		return errors.Wrap(err, "reading output commitment")
	}

	// read and ignore the (empty) output witness
	_, err = blockchain.ReadVarstr31(r)
	return errors.Wrap(err, "reading output witness")
}

func (to *ClaimTxOutput) writeTo(w io.Writer) error {
	if _, err := blockchain.WriteVarint63(w, to.AssetVersion); err != nil {
		return errors.Wrap(err, "writing asset version")
	}

	if err := to.writeCommitment(w); err != nil {
		return errors.Wrap(err, "writing output commitment")
	}

	if _, err := blockchain.WriteVarstr31(w, nil); err != nil {
		return errors.Wrap(err, "writing witness")
	}
	return nil
}

func (to *ClaimTxOutput) writeCommitment(w io.Writer) error {
	return to.OutputCommitment.writeExtensibleString(w, to.CommitmentSuffix, to.AssetVersion)
}
