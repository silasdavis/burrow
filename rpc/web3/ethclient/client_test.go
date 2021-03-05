// +build ethereum

package ethclient

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/rpc/rpcevents"

	"github.com/hyperledger/burrow/rpc/lib/jsonrpc"
	"github.com/stretchr/testify/require"
)

const truffleRemote = "http://127.0.0.1:9545"
const infuraRemote = "https://ropsten.infura.io/v3/7ed3059377654803a190fa44560d528f"

var client = jsonrpc.NewClient(truffleRemote)

//var client = jsonrpc.NewClient(infuraRemote)

func TestEthGetLogs(t *testing.T) {
	params := Filter{
		BlockRange: rpcevents.AbsoluteRange(1, 34340),
		Addresses: []crypto.Address{
			crypto.MustAddressFromHexString("a1e378f122fec6aa8c841397042e21bc19368768"),
			crypto.MustAddressFromHexString("f73aaa468496a87675d27638878a1600b0db3c71"),
		},
	}
	result, err := EthGetLogs(client, params)
	require.NoError(t, err)
	bs, err := json.Marshal(result)
	require.NoError(t, err)
	fmt.Printf("%s\n", string(bs))
}

func TestNetVersion(t *testing.T) {
	result, err := NetVersion(client)
	require.NoError(t, err)
	fmt.Printf("%#v\n", result)
}

func TestWeb3ClientVersion(t *testing.T) {
	result, err := Web3ClientVersion(client)
	require.NoError(t, err)
	fmt.Printf("%#v\n", result)
}

func TestEthSyncing(t *testing.T) {
	result, err := EthSyncing(client)
	require.NoError(t, err)
	fmt.Printf("%#v\n", result)
}

func TestEthBlockNumber(t *testing.T) {
	result, err := EthBlockNumber(client)
	require.NoError(t, err)
	fmt.Printf("%d\n", result)
}
