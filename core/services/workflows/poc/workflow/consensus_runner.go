package workflow

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

var _ capability = &consensusRunner[any, any]{}

type consensusRunner[I, O any] struct {
	nonTriggerCapability
	capabilities.Consensus[I, O]
}

func (c consensusRunner[I, O]) CapabilityType() commoncap.CapabilityType {
	return commoncap.CapabilityTypeConsensus
}

func (c consensusRunner[I, O]) Run(_ string, value values.Value) (values.Value, bool, error) {
	observations := value.(*values.Map).Underlying["observations"]
	// Seems odd that consensus implementations must all know that this can be a list or not a list...
	if _, ok := observations.(*values.List); !ok {
		observations, _ = values.NewList([]any{observations})
	}
	inputs, err := capabilities.UnwrapValue[*[]I](observations)
	if err != nil {
		return nil, false, err
	}
	consensus, err := c.Invoke(*inputs)
	if err != nil {
		return nil, false, err

	}
	result, err := values.Wrap(consensus.Results())
	return result, err == nil, err
}
