package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/encoding"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/rpc"
	"github.com/hyperledger/burrow/rpc/rpcevents"
	"github.com/hyperledger/burrow/rpc/rpcquery"
	"github.com/hyperledger/burrow/vent/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type BurrowChain struct {
	conn  *grpc.ClientConn
	query rpcquery.QueryClient
	exec  rpcevents.ExecutionEventsClient
}

var _ Chain = (*BurrowChain)(nil)

func NewBurrow(grpcAddr string) (*BurrowChain, error) {
	conn, err := encoding.GRPCDial(grpcAddr)
	if err != nil {
		return nil, err
	}
	return &BurrowChain{
		conn:  conn,
		query: rpcquery.NewQueryClient(conn),
		exec:  rpcevents.NewExecutionEventsClient(conn),
	}, nil
}

func (b *BurrowChain) ConsumeBlockExecutions(
	ctx context.Context,
	in *rpcevents.BlocksRequest,
	consumer func(Block) error,
	continuityOptions ...exec.ContinuityOpt) error {

	stream, err := b.exec.Stream(ctx, in)
	if err != nil {
		return fmt.Errorf("could not connect to block stream: %w", err)
	}

	return rpcevents.ConsumeBlockExecutions(stream, func(blockExecution *exec.BlockExecution) error {
		return consumer((*BurrowBlock)(blockExecution))
	}, continuityOptions...)
}

func (b *BurrowChain) Status(ctx context.Context, in *rpcquery.StatusParam) (*rpc.ResultStatus, error) {
	return b.query.Status(ctx, in)
}

func (b *BurrowChain) Connectivity() connectivity.State {
	return b.conn.GetState()
}

func (b *BurrowChain) GetMetadata(ctx context.Context, in *rpcquery.GetMetadataParam) (*rpcquery.MetadataResult, error) {
	return b.query.GetMetadata(ctx, in)
}

func (b *BurrowChain) Close() error {
	return b.conn.Close()
}

type BurrowBlock exec.BlockExecution

func NewBurrowBlock(block *exec.BlockExecution) *BurrowBlock {
	return (*BurrowBlock)(block)
}

func (b *BurrowBlock) GetMetadata(columns types.SQLColumnNames) (map[string]interface{}, error) {
	blockHeader, err := json.Marshal(b.Header)
	if err != nil {
		return nil, fmt.Errorf("could not marshal block header: %w", err)
	}

	return map[string]interface{}{
		columns.Height:      fmt.Sprintf("%v", b.Height),
		columns.BlockHeader: string(blockHeader),
	}, nil
}

var _ Block = (*BurrowBlock)(nil)

func (b *BurrowBlock) GetTime() time.Time {
	return b.Header.GetTime()
}

func (b *BurrowBlock) GetChainID() string {
	return b.Header.GetChainID()
}
func (b *BurrowBlock) GetHeight() uint64 {
	return b.Height
}

func (b *BurrowBlock) GetTxs() []Transaction {
	txs := make([]Transaction, len(b.TxExecutions))
	for i, tx := range b.TxExecutions {
		txs[i] = (*BurrowTx)(tx)
	}
	return txs
}

type BurrowTx exec.TxExecution

var _ Transaction = (*BurrowTx)(nil)

func (tx *BurrowTx) GetException() *errors.Exception {
	return tx.Exception
}

func (tx *BurrowTx) GetMetadata(columns types.SQLColumnNames) (map[string]interface{}, error) {
	// transaction raw data
	envelope, err := json.Marshal(tx.Envelope)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal envelope in tx %v: %v", tx, err)
	}

	events, err := json.Marshal(tx.Events)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal events in tx %v: %v", tx, err)
	}

	result, err := json.Marshal(tx.Result)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal result in tx %v: %v", tx, err)
	}

	receipt, err := json.Marshal(tx.Receipt)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal receipt in tx %v: %v", tx, err)
	}

	exception, err := json.Marshal(tx.Exception)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal exception in tx %v: %v", tx, err)
	}

	origin, err := json.Marshal(tx.Origin)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal origin in tx %v: %v", tx, err)
	}

	return map[string]interface{}{
		columns.Height:    tx.Height,
		columns.TxHash:    tx.TxHash.String(),
		columns.TxIndex:   tx.Index,
		columns.TxType:    tx.TxType.String(),
		columns.Envelope:  string(envelope),
		columns.Events:    string(events),
		columns.Result:    string(result),
		columns.Receipt:   string(receipt),
		columns.Origin:    string(origin),
		columns.Exception: string(exception),
	}, nil
}

func (tx *BurrowTx) GetHash() binary.HexBytes {
	return tx.TxHash
}

func (tx *BurrowTx) GetEvents() []Event {
	// All txs have events, but not all have LogEvents
	var events []Event
	for _, ev := range tx.Events {
		if ev.Log != nil {
			events = append(events, (*BurrowEvent)(ev))
		}
	}
	return events
}

type BurrowEvent exec.Event

var _ Event = (*BurrowEvent)(nil)

func (ev *BurrowEvent) GetTransactionHash() binary.HexBytes {
	return ev.Header.TxHash
}

func (ev *BurrowEvent) GetIndex() uint64 {
	return ev.Header.Index
}

func (ev *BurrowEvent) GetTopics() []binary.Word256 {
	return ev.Log.Topics
}

func (ev *BurrowEvent) GetData() []byte {
	return ev.Log.Data
}

func (ev *BurrowEvent) GetAddress() crypto.Address {
	return ev.Log.Address
}

func (ev *BurrowEvent) GetSolidityEventID() abi.EventID {
	return ev.Log.SolidityEventID()
}

// Tags
func (ev *BurrowEvent) Get(key string) (value interface{}, ok bool) {
	return (*exec.Event)(ev).Get(key)
}
