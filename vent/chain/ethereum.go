package chain

import (
	"context"

	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/rpc"
	"github.com/hyperledger/burrow/rpc/lib/jsonrpc"
	"github.com/hyperledger/burrow/rpc/rpcevents"
	"github.com/hyperledger/burrow/rpc/rpcquery"
	"google.golang.org/grpc/connectivity"
)

type EthereumChain struct {
	client rpc.Client
}

//var _ Chain = (*EthereumChain)(nil)

func NewEthereum(remote string) (*EthereumChain, error) {
	return &EthereumChain{
		client: jsonrpc.NewClient(remote),
	}, nil
}

func (e *EthereumChain) ConsumeBlockExecutions(
	ctx context.Context,
	in *rpcevents.BlocksRequest,
	consumer func(*exec.BlockExecution) error,
	continuityOptions ...exec.ContinuityOpt) error {

	//web3client.ConsumeBlockExecutions()
	panic("implement me")
}

func (e *EthereumChain) Status(ctx context.Context, in *rpcquery.StatusParam) (*rpc.ResultStatus, error) {
	panic("implement me")
}

func (e *EthereumChain) Connectivity() connectivity.State {
	panic("implement me")
}

func (e *EthereumChain) GetMetadata(ctx context.Context, in *rpcquery.GetMetadataParam) (*rpcquery.MetadataResult, error) {
	panic("implement me")
}

func (e *EthereumChain) Close() error {
	panic("implement me")
}
