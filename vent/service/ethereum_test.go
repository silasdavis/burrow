// +build ethereum

package service

import (
	"fmt"
	"testing"

	"github.com/hyperledger/burrow/rpc/lib/jsonrpc"
	"github.com/hyperledger/burrow/rpc/web3"
	"github.com/hyperledger/burrow/rpc/web3/web3client"
	"github.com/stretchr/testify/require"
)

func TestGetLogs(t *testing.T) {
	client := jsonrpc.NewClient("http://127.0.0.1:9545")
	params := &web3.EthGetLogsParams{
		Filter: web3.Filter{
			FromBlock: "0x0",
			//ToBlock:   "3",
			//Address:   "0x1cf428867e7df1F215776c18EF01c3De2D64c525",
			//Topics:    nil,
		},
	}
	result, err := web3client.EthGetLogs(client, params)
	require.NoError(t, err)
	fmt.Println(result)
}
