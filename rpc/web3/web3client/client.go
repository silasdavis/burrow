package web3client

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/rpc"
	"github.com/hyperledger/burrow/rpc/web3"
)

const (
	GetLogsMethod = "eth_getLogs"
)

// Hex string optionally prefixed with 0x
type HexString string

func (hs HexString) Bytes() ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(string(hs), "0x"))
}

func (hs HexString) BigInt() (*big.Int, error) {
	bs, err := hs.Bytes()
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(bs), nil
}

func (hs HexString) Uint64() (uint64, error) {
	bi, err := hs.BigInt()
	if err != nil {
		return 0, err
	}
	if !bi.IsUint64() {
		return 0, fmt.Errorf("%v is not uint64", bi)
	}
	return bi.Uint64(), nil
}

func (hs HexString) Address() (crypto.Address, error) {
	bs, err := hs.Bytes()
	if err != nil {
		return crypto.Address{}, err
	}
	return crypto.AddressFromBytes(bs)
}

type EthLog struct {
	Topics []HexString `json:"topics"`
	// Hex representation of a Keccak 256 hash
	TransactionHash HexString `json:"transactionHash"`
	// Sender of the transaction
	Address HexString `json:"address"`
	// The hex representation of the Keccak 256 of the RLP encoded block
	BlockHash HexString `json:"blockHash"`
	// The hex representation of the block's height
	BlockNumber HexString `json:"blockNumber"`
	// Hex representation of a variable length byte array
	Data HexString `json:"data"`
	// Hex representation of the integer
	LogIndex HexString `json:"logIndex"`
	// Hex representation of the integer
	TransactionIndex HexString `json:"transactionIndex"`
}

type Filter struct {
	// The hex representation of the block's height
	FromBlock uint64 `json:"fromBlock"`
	// The hex representation of the block's height
	ToBlock uint64         `json:"toBlock"`
	Address crypto.Address `json:"address"`
	// Array of 32 Bytes DATA topics. Topics are order-dependent. Each topic can also be an array of DATA with 'or' options
	Topics []binary.Word256 `json:"topics"`
}

func EthGetLogs(client rpc.Client, params *web3.EthGetLogsParams) ([]EthLog, error) {
	var logs []EthLog
	_, err := client.Call(GetLogsMethod, map[string]interface{}{
		"fromBlock": params.FromBlock,
		"toBlock":   params.ToBlock,
		"address":   params.Address,
		"topics":    params.Topics,
	}, &logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func ConsumeBlockExecutions(logs []EthLog, consumer func(*exec.BlockExecution) error, continuityOptions ...exec.ContinuityOpt) error {
	for _, log := range logs {
		height, err := log.BlockNumber.Uint64()
		if err != nil {
			return fmt.Errorf("could not extract ethereum block height: %w", err)
		}
		var predecessor uint64
		if height > 0 {
			predecessor = height - 1
		}
		err = consumer(&exec.BlockExecution{
			Height:            height,
			PredecessorHeight: predecessor,
			Header:            nil,
			TxExecutions:      nil,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

//func ConsumeBlockExecutions(stream ExecutionEvents_StreamClient, consumer func(*exec.BlockExecution) error,
//	continuityOptions ...exec.ContinuityOpt) error {
//	var be *exec.BlockExecution
//	var err error
//	ba := exec.NewBlockAccumulator(continuityOptions...)
//	for be, err = ba.ConsumeBlockExecution(stream); err == nil; be, err = ba.ConsumeBlockExecution(stream) {
//		err = consumer(be)
//		if err != nil {
//			return err
//		}
//	}
//	return err
//}
