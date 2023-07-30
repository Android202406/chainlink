package evm

import (
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

//go:generate mockery --quiet --name LoopRelayAdapter --output ./mocks/ --case=underscore
type LoopRelayAdapter interface {
	loop.Relayer
	Chain() evm.Chain
	Default() bool
}
type LoopRelayer struct {
	loop.Relayer
	x evm.EVMChainRelayerExtender
}

var _ loop.Relayer = &LoopRelayer{}

func NewLoopRelayAdapter(r *Relayer, cs evm.EVMChainRelayerExtender) *LoopRelayer {
	ra := relay.NewRelayerAdapter(r, cs)
	return &LoopRelayer{
		Relayer: ra,
		x:       cs,
	}
}

func (la *LoopRelayer) Chain() evm.Chain {
	return la.x.Chain()
}

func (la *LoopRelayer) Default() bool {
	return la.x.Default()
}

//TODO need service multi start/close, etc for the contained relayer and extender?
