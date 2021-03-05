package burrow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/burrow/event/query"
	"github.com/hyperledger/burrow/vent/chain"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/encoding"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/rpc/rpcevents"
	"github.com/hyperledger/burrow/rpc/rpcquery"
	"github.com/hyperledger/burrow/vent/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type BurrowChain struct {
	conn    *grpc.ClientConn
	filter  query.Query
	query   rpcquery.QueryClient
	exec    rpcevents.ExecutionEventsClient
	chainID string
	version string
}

var _ chain.Chain = (*BurrowChain)(nil)

func NewBurrow(grpcAddr string) (*BurrowChain, error) {
	conn, err := encoding.GRPCDial(grpcAddr)
	if err != nil {
		return nil, err
	}
	client := rpcquery.NewQueryClient(conn)
	status, err := client.Status(context.Background(), &rpcquery.StatusParam{})
	if err != nil {
		return nil, fmt.Errorf("could not get initial status from Burrow: %w", err)
	}
	return &BurrowChain{
		conn:    conn,
		query:   client,
		exec:    rpcevents.NewExecutionEventsClient(conn),
		chainID: status.ChainID,
		version: status.BurrowVersion,
	}, nil
}

func (b *BurrowChain) FilterBy(filter chain.Filter) error {
	qb := query.NewBuilder()
	for _, address := range filter.Addresses {
		qb = qb.AndEquals("Address", address)
	}
	for i, topic := range filter.Topics {
		qb = qb.AndEquals(exec.LogNKey(i), topic)
	}
	var err error
	b.filter, err = qb.Query()
	if err != nil {
		return fmt.Errorf("could not build Vent filter query: %w", err)
	}
	return nil
}

func (b *BurrowChain) GetChainID() string {
	return b.chainID
}

func (b *BurrowChain) GetVersion() string {
	return b.version
}

func (b *BurrowChain) StatusMessage(ctx context.Context, lastProcessedHeight uint64) []interface{} {
	var catchUpRatio float64
	status, err := b.query.Status(ctx, &rpcquery.StatusParam{})
	if err != nil {
		err = fmt.Errorf("could not get Burrow chain status: %w", err)
		return []interface{}{
			"msg", "status",
			"error", err.Error(),
		}
	}
	if status.SyncInfo.LatestBlockHeight > 0 {
		catchUpRatio = float64(lastProcessedHeight) / float64(status.SyncInfo.LatestBlockHeight)
	}
	return []interface{}{
		"msg", "status",
		"last_processed_height", lastProcessedHeight,
		"fraction_caught_up", catchUpRatio,
		"burrow_latest_block_height", status.SyncInfo.LatestBlockHeight,
		"burrow_latest_block_duration", status.SyncInfo.LatestBlockDuration,
		"burrow_latest_block_hash", status.SyncInfo.LatestBlockHash,
		"burrow_latest_app_hash", status.SyncInfo.LatestAppHash,
		"burrow_latest_block_time", status.SyncInfo.LatestBlockTime,
		"burrow_latest_block_seen_time", status.SyncInfo.LatestBlockSeenTime,
		"burrow_node_info", status.NodeInfo,
		"burrow_catching_up", status.CatchingUp,
	}
}

func (b *BurrowChain) ConsumeBlocks(
	ctx context.Context,
	in *rpcevents.BlockRange,
	consumer func(chain.Block) error,
	continuityOptions ...exec.ContinuityOpt) error {

	stream, err := b.exec.Stream(ctx, &rpcevents.BlocksRequest{
		BlockRange: in,
		Query:      b.filter.String(),
	})
	if err != nil {
		return fmt.Errorf("could not connect to block stream: %w", err)
	}

	return rpcevents.ConsumeBlockExecutions(stream, func(blockExecution *exec.BlockExecution) error {
		return consumer((*BurrowBlock)(blockExecution))
	}, continuityOptions...)
}

func (b *BurrowChain) Connectivity() connectivity.State {
	return b.conn.GetState()
}

func (b *BurrowChain) GetABI(ctx context.Context, address crypto.Address) (string, error) {
	result, err := b.query.GetMetadata(ctx, &rpcquery.GetMetadataParam{
		Address: &address,
	})
	if err != nil {
		return "", err
	}
	return result.Metadata, nil
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

var _ chain.Block = (*BurrowBlock)(nil)

func (b *BurrowBlock) GetTime() time.Time {
	return b.Header.GetTime()
}

func (b *BurrowBlock) GetHeight() uint64 {
	return b.Height
}

func (b *BurrowBlock) GetTxs() []chain.Transaction {
	txs := make([]chain.Transaction, len(b.TxExecutions))
	for i, tx := range b.TxExecutions {
		txs[i] = (*BurrowTx)(tx)
	}
	return txs
}

type BurrowTx exec.TxExecution

var _ chain.Transaction = (*BurrowTx)(nil)

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

func (tx *BurrowTx) GetEvents() []chain.Event {
	// All txs have events, but not all have LogEvents
	var events []chain.Event
	for _, ev := range tx.Events {
		if ev.Log != nil {
			events = append(events, (*BurrowEvent)(ev))
		}
	}
	return events
}

type BurrowEvent exec.Event

var _ chain.Event = (*BurrowEvent)(nil)

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

// Tags
func (ev *BurrowEvent) Get(key string) (value interface{}, ok bool) {
	return (*exec.Event)(ev).Get(key)
}
