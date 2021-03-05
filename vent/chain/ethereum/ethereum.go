package ethereum

import (
	"context"
	"fmt"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/vent/chain"
	"time"

	"github.com/hyperledger/burrow/rpc/web3/ethclient"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/rpc"
	"github.com/hyperledger/burrow/rpc/lib/jsonrpc"
	"github.com/hyperledger/burrow/rpc/rpcevents"
	"github.com/hyperledger/burrow/vent/types"
	"google.golang.org/grpc/connectivity"
)

type EthereumChain struct {
	client  rpc.Client
	filter  *chain.Filter
	chainID string
	version string
	logger  *logging.Logger
}

func (c *EthereumChain) FilterBy(filter chain.Filter) error {
	panic("implement me")
}

func (c *EthereumChain) StatusMessage(ctx context.Context, lastProcessedHeight uint64) []interface{} {
	panic("implement me")
}

var _ chain.Chain = (*EthereumChain)(nil)

func NewEthereumChain(client rpc.Client, filter *chain.Filter, logger *logging.Logger) (*EthereumChain, error) {
	chainID, err := ethclient.NetVersion(client)
	if err != nil {
		return nil, fmt.Errorf("could not get Ethereum ChainID: %w", err)
	}
	version, err := ethclient.Web3ClientVersion(client)
	if err != nil {
		return nil, fmt.Errorf("could not get Ethereum node version: %w", err)
	}
	return &EthereumChain{
		client:  client,
		filter:  filter,
		chainID: chainID,
		version: version,
		logger:  logger,
	}, nil
}

func (c *EthereumChain) GetABI(ctx context.Context, address crypto.Address) (string, error) {
	// Unsupported by Ethereum
	return "", nil
}

func (c *EthereumChain) GetVersion() string {
	return c.version
}

func (c *EthereumChain) GetChainID() string {
	return c.chainID
}

func (c *EthereumChain) ConsumeBlocks(ctx context.Context, in *rpcevents.BlockRange, consumer func(chain.Block) error,
	continuityOptions ...exec.ContinuityOpt) error {
	return Consume(c.client, c.filter, in, c.logger, consumer)
}

func NewEthereum(remote string) (*EthereumChain, error) {
	return &EthereumChain{
		client: jsonrpc.NewClient(remote),
	}, nil
}

func (c *EthereumChain) Connectivity() connectivity.State {
	// Use ethSyncing
	panic("implement me")
}

func (c *EthereumChain) Close() error {
	// just a http.Client - nothing to free
	return nil
}

type EthereumBlock struct {
	Height       uint64
	Transactions []chain.Transaction
}

func NewEthereumBlock(log *EthereumEvent) *EthereumBlock {
	return &EthereumBlock{
		Height:       log.Height,
		Transactions: []chain.Transaction{NewEthereumTransaction(log)},
	}
}

var _ chain.Block = (*EthereumBlock)(nil)

func (b *EthereumBlock) GetHeight() uint64 {
	return b.Height
}

func (b *EthereumBlock) GetTxs() []chain.Transaction {
	return b.Transactions
}

func (b *EthereumBlock) GetTime() time.Time {
	panic("implement me")
}

func (b *EthereumBlock) GetMetadata(columns types.SQLColumnNames) (map[string]interface{}, error) {
	panic("implement me")
}

func (b *EthereumBlock) appendTransaction(log *EthereumEvent) {
	b.Transactions = append(b.Transactions, &EthereumTransaction{
		Index:  uint64(len(b.Transactions)),
		Hash:   log.TransactionHash,
		Events: []chain.Event{log},
	})
}

func (b *EthereumBlock) appendEvent(log *EthereumEvent) {
	tx := b.Transactions[len(b.Transactions)-1].(*EthereumTransaction)
	log.Index = uint64(len(tx.Events))
	tx.Events = append(tx.Events, log)
}

type EthereumTransaction struct {
	Height uint64
	Index  uint64
	Hash   binary.HexBytes
	Events []chain.Event
}

func NewEthereumTransaction(log *EthereumEvent) *EthereumTransaction {
	return &EthereumTransaction{
		Height: log.Height,
		Index:  0,
		Hash:   log.TransactionHash,
		Events: []chain.Event{log},
	}
}

func (tx *EthereumTransaction) GetHash() binary.HexBytes {
	return tx.Hash
}

func (tx *EthereumTransaction) GetIndex() uint64 {
	return tx.Index
}

func (tx *EthereumTransaction) GetEvents() []chain.Event {
	return tx.Events
}

func (tx *EthereumTransaction) GetException() *errors.Exception {
	// Ethereum does not retain an log from reverted transactions
	return nil
}

func (tx *EthereumTransaction) GetOrigin() *exec.Origin {
	// Origin refers to a previous dumped chain which is not a concept in Ethereum
	return nil
}

func (tx *EthereumTransaction) GetMetadata(columns types.SQLColumnNames) (map[string]interface{}, error) {
	return map[string]interface{}{
		columns.Height:  tx.Height,
		columns.TxHash:  tx.Hash.String(),
		columns.TxIndex: tx.Index,
		columns.TxType:  exec.TypeLog.String(),
	}, nil
}

var _ chain.Transaction = (*EthereumTransaction)(nil)

type EthereumEvent struct {
	exec.LogEvent
	Height uint64
	// Index of event in entire block (what ethereum provides us with
	IndexInBlock uint64
	// Index of event in transaction
	Index           uint64
	TransactionHash binary.HexBytes
}

var _ chain.Event = (*EthereumEvent)(nil)

func NewEthereumEvent(log *ethclient.EthLog) (*EthereumEvent, error) {
	height, err := log.BlockNumber.Uint64()
	if err != nil {
		return nil, fmt.Errorf("could not parse BlockNumber '%s': %w", log.BlockNumber, err)
	}
	indexInBlock, err := log.LogIndex.Uint64()
	if err != nil {
		return nil, fmt.Errorf("could not parse LogIndex '%s': %w", log.LogIndex, err)
	}
	txHash, err := log.TransactionHash.Bytes()
	if err != nil {
		return nil, fmt.Errorf("could not parse TransactionHash '%s': %w", log.TransactionHash, err)
	}
	topics, err := log.GetTopics()
	if err != nil {
		return nil, fmt.Errorf("could not parse Topics '%s': %w", log.Topics, err)
	}
	address, err := log.Address.Address()
	if err != nil {
		return nil, fmt.Errorf("could not parse Address '%s': %w", log.Address, err)
	}
	data, err := log.Data.Bytes()
	if err != nil {
		return nil, fmt.Errorf("could not parse BlockNumber '%s': %w", log.Data, err)
	}
	return &EthereumEvent{
		LogEvent: exec.LogEvent{
			Topics:  topics,
			Address: address,
			Data:    data,
		},
		Height:          height,
		IndexInBlock:    indexInBlock,
		TransactionHash: txHash,
	}, nil
}

func (ev *EthereumEvent) GetIndex() uint64 {
	return ev.Index
}

func (ev *EthereumEvent) GetTransactionHash() binary.HexBytes {
	return ev.TransactionHash
}

func (ev *EthereumEvent) GetAddress() crypto.Address {
	return ev.Address
}

func (ev *EthereumEvent) GetTopics() []binary.Word256 {
	return ev.Topics
}

func (ev *EthereumEvent) GetData() []byte {
	return ev.Data
}

func (ev *EthereumEvent) Get(key string) (value interface{}, ok bool) {
	return ev.LogEvent.Get(key)
}
