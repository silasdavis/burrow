// Abstracts over a Burrow GRPC connection and Ethereum json-rpc web3 connection for the purposes of vent

package chain

import (
	"context"
	"time"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/event/query"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/rpc"
	"github.com/hyperledger/burrow/rpc/rpcevents"
	"github.com/hyperledger/burrow/rpc/rpcquery"
	"github.com/hyperledger/burrow/vent/types"
	"google.golang.org/grpc/connectivity"
)

type Chain interface {
	ConsumeBlockExecutions(
		ctx context.Context,
		in *rpcevents.BlocksRequest,
		consumer func(Block) error,
		continuityOptions ...exec.ContinuityOpt) error
	Status(ctx context.Context, in *rpcquery.StatusParam) (*rpc.ResultStatus, error)
	Connectivity() connectivity.State
	GetMetadata(ctx context.Context, in *rpcquery.GetMetadataParam) (*rpcquery.MetadataResult, error)
	Close() error
}

type Block interface {
	GetHeight() uint64
	GetTxs() []Transaction
	GetTime() time.Time
	GetChainID() string
	GetMetadata(columns types.SQLColumnNames) (map[string]interface{}, error)
}

type Transaction interface {
	GetHash() binary.HexBytes
	GetIndex() uint64
	GetEvents() []Event
	GetException() *errors.Exception
	GetOrigin() *exec.Origin
	GetMetadata(columns types.SQLColumnNames) (map[string]interface{}, error)
}

type Event interface {
	query.Tagged
	GetIndex() uint64
	GetTransactionHash() binary.HexBytes
	GetSolidityEventID() abi.EventID
	GetAddress() crypto.Address
	GetTopics() []binary.Word256
	GetData() []byte
}
