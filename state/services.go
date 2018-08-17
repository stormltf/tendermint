package state

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
)

//------------------------------------------------------
// blockchain services types
// NOTE: Interfaces used by RPC must be thread safe!
//------------------------------------------------------

//------------------------------------------------------
// mempool

// Mempool defines the mempool interface as used by the ConsensusState.
// Updates to the mempool need to be synchronized with committing a block
// so apps can reset their transient state on Commit
type Mempool interface {
	Lock()
	Unlock()

	Size() int
	CheckTx(types.Tx, func(*abci.Response)) error
	ReapMaxBytes(max int) types.Txs
	Update(height int64, txs types.Txs) error
	Flush()
	FlushAppConn() error

	TxsAvailable() <-chan struct{}
	EnableTxsAvailable()

	SetFilter(func(types.Tx) bool)
}

// MockMempool is an empty implementation of a Mempool, useful for testing.
type MockMempool struct {
}

func (MockMempool) Lock()                                              {}
func (MockMempool) Unlock()                                            {}
func (MockMempool) Size() int                                          { return 0 }
func (MockMempool) CheckTx(tx types.Tx, cb func(*abci.Response)) error { return nil }
func (MockMempool) ReapMaxBytes(max int) types.Txs                     { return types.Txs{} }
func (MockMempool) Update(height int64, txs types.Txs) error           { return nil }
func (MockMempool) Flush()                                             {}
func (MockMempool) FlushAppConn() error                                { return nil }
func (MockMempool) TxsAvailable() <-chan struct{}                      { return make(chan struct{}) }
func (MockMempool) EnableTxsAvailable()                                {}
func (MockMempool) SetFilter(func(types.Tx) bool)                      {}

//------------------------------------------------------
// blockstore

// BlockStoreRPC is the block store interface used by the RPC.
type BlockStoreRPC interface {
	Height() int64

	LoadBlockMeta(height int64) *types.BlockMeta
	LoadBlock(height int64) *types.Block
	LoadBlockPart(height int64, index int) *types.Part

	LoadBlockCommit(height int64) *types.Commit
	LoadSeenCommit(height int64) *types.Commit
}

// BlockStore defines the BlockStore interface used by the ConsensusState.
type BlockStore interface {
	BlockStoreRPC
	SaveBlock(block *types.Block, blockParts *types.PartSet, seenCommit *types.Commit)
}

//-----------------------------------------------------------------------------------------------------
// evidence pool

// EvidencePool defines the EvidencePool interface used by the ConsensusState.
type EvidencePool interface {
	PendingEvidence(limit int) []types.Evidence
	AddEvidence(types.Evidence) error
	Update(*types.Block, State)
}

// MockMempool is an empty implementation of a Mempool, useful for testing.
type MockEvidencePool struct {
}

func (m MockEvidencePool) PendingEvidence(limit int) []types.Evidence { return nil }
func (m MockEvidencePool) AddEvidence(types.Evidence) error           { return nil }
func (m MockEvidencePool) Update(*types.Block, State)                 {}
