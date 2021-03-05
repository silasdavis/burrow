package ethclient

import (
	"fmt"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/rpc/rpcevents"
	"github.com/tmthrgd/go-hex"
	"math/big"
	"strings"
	bin "encoding/binary"

	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/rpc"
)

const (
	EthGetLogsMethod     = "eth_getLogs"
	EthSyncingMethod     = "eth_syncing"
	EthBlockNumberMethod = "eth_blockNumber"
	NetVersionMethod     = "net_version"
	Web3ClientVersionMethod  = "web3_clientVersion"
)

// Hex string optionally prefixed with 0x
type HexString string

func NewHexString(bs []byte) HexString {
	str := hex.EncodeToString(bs)
	// Ethereum expects leading zeros to be removed (SMH)
	if str[0] == '0' {
		str = str[1:]
	}
	return HexString("0x" + str)
}

func NewHexUint64(x uint64) HexString {
	bs := make([]byte, 8)
	bin.BigEndian.PutUint64(bs, x)
	return NewHexString(bs)
}

func (hs HexString) Bytes() ([]byte, error) {
	hexString := strings.TrimPrefix(string(hs), "0x")
	// Ethereum return odd-length hexString strings when it remove leading 0
	if len(hexString)%2 == 1 {
		hexString = "0" + hexString
	}
	return hex.DecodeString(hexString)
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

func (log *EthLog) GetTopics() ([]binary.Word256, error) {
	topics := make([]binary.Word256, len(log.Topics))
	for i, t := range log.Topics {
		bs, err := t.Bytes()
		if err != nil {
			return nil, err
		}
		topics[i] = binary.LeftPadWord256(bs)
	}
	return topics, nil
}

type Filter struct {
	*rpcevents.BlockRange
	Addresses []crypto.Address
	Topics    []binary.Word256
}

func (f *Filter) EthFilter() *EthFilter {
	topics := make([]HexString, len(f.Topics))
	for i, t := range f.Topics {
		topics[i] = NewHexString(t[:])
	}
	addresses := make([]HexString, len(f.Addresses))
	for i, a := range f.Addresses {
		addresses[i] = NewHexString(a.Bytes())
	}
	return &EthFilter{
		FromBlock: ethLogBound(f.GetStart()),
		ToBlock:   ethLogBound(f.GetEnd()),
		Addresses: addresses,
		Topics:    topics,
	}
}

type EthFilter struct {
	// The hex representation of the block's height
	FromBlock string `json:"fromBlock,omitempty"`
	// The hex representation of the block's height
	ToBlock string `json:"toBlock,omitempty"`
	// Yes this is JSON address since allowed to be singular
	Addresses []HexString `json:"address,omitempty"`
	// Array of 32 Bytes DATA topics. Topics are order-dependent. Each topic can also be an array of DATA with 'or' options
	Topics []HexString `json:"topics,omitempty"`
}

func EthGetLogs(client rpc.Client, filter Filter) ([]*EthLog, error) {
	var logs []*EthLog
	_, err := client.Call(EthGetLogsMethod, []*EthFilter{filter.EthFilter()}, &logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func EthSyncing(client rpc.Client) (bool, error) {
	syncing := new(bool)
	_, err := client.Call(EthSyncingMethod, nil, syncing)
	if err != nil {
		return false, err
	}
	return *syncing, nil
}

func EthBlockNumber(client rpc.Client) (uint64, error) {
	latestBlock := new(HexString)
	_, err := client.Call(EthBlockNumberMethod, nil, latestBlock)
	if err != nil {
		return 0, err
	}
	return latestBlock.Uint64()
}

// AKA ChainID
func NetVersion(client rpc.Client) (string, error) {
	version := new(string)
	_, err := client.Call(NetVersionMethod, nil, version)
	if err != nil {
		return "", err
	}
	return *version, nil
}

func Web3ClientVersion(client rpc.Client) (string, error) {
	version := new(string)
	_, err := client.Call(Web3ClientVersionMethod, nil, version)
	if err != nil {
		return "", err
	}
	return *version, nil
}

func ethLogBound(bound *rpcevents.Bound) string {
	if bound == nil {
		return ""
	}
	switch bound.Type {
	case rpcevents.Bound_FIRST:
		return "earliest"
	case rpcevents.Bound_LATEST:
		return "latest"
	case rpcevents.Bound_ABSOLUTE:
		return string(NewHexUint64(bound.Index))
	default:
		return ""
	}
}
