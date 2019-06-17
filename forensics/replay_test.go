// +build forensics

package forensics

import (
	"encoding/hex"
	"fmt"
	"path"
	"testing"

	"github.com/magiconair/properties/assert"

	"github.com/hyperledger/burrow/config/source"
	"github.com/hyperledger/burrow/execution/state"
	"github.com/hyperledger/burrow/genesis"
	"github.com/hyperledger/burrow/logging"
	"github.com/stretchr/testify/require"
)

// This serves as a testbed for looking at non-deterministic burrow instances capture from the wild
// Put the path to 'good' and 'bad' burrow directories here (containing the config files and .burrow dir)
//const goodDir = "/home/silas/test-chain"
const goodDir = "/home/silas/burrows/production-t9/burrow-t9-studio-001-good"
const badDir = "/home/silas/burrows/production-t9/burrow-t9-studio-000-bad"
const criticalBlock uint64 = 6

func TestReplay_Compare(t *testing.T) {
	badReplay := newReplay(t, badDir)
	goodReplay := newReplay(t, goodDir)
	badRecaps, err := badReplay.Blocks(1, criticalBlock+1)
	require.NoError(t, err)
	goodRecaps, err := goodReplay.Blocks(1, criticalBlock+1)
	require.NoError(t, err)
	//for i, goodRecap := range goodRecaps {
	//	fmt.Printf("Good: %v\n", goodRecap)
	//	fmt.Printf("Bad: %v\n", badRecaps[i])
	//	assert.Equal(t, goodRecap, badRecaps[i])
	//	for i, txe := range goodRecap.TxExecutions {
	//		fmt.Printf("Tx %d: %v\n", i, txe.TxHash)
	//		fmt.Println(txe.Envelope)
	//	}
	//	fmt.Println()
	//}

	txe := goodRecaps[5].TxExecutions[0]
	assert.Equal(t, badRecaps[5].TxExecutions[0], txe)
	fmt.Printf("%v\n", txe.Envelope.Signatories[0])
}

func TestDecipher(t *testing.T) {
	hexmsg:= "7B22436861696E4944223A2270726F64756374696F6E2D74392D73747564696F2D627572726F772D364337333335222C2254797065223A2243616C6C5478222C225061796C6F6164223A7B22496E707574223A7B2241646472657373223A2236354139334431443333423633453932453942454335463938444633313638303033384530303431222C2253657175656E6365223A34307D2C2241646472657373223A2242413544333042313031393233363033444331333133313231334431334633443939354138344142222C224761734C696D6974223A393939393939392C2244617461223A224636373138374143303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303032303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030304534343635363136343643363936453635344637323631363336433635303030303030303030303030303030303030303030303030303030303030303030303030222C225741534D223A22227D7D"
	bs, err := hex.DecodeString(hexmsg)
	require.NoError(t, err)
	fmt.Println(string(bs))
}

func TestReplay_Good(t *testing.T) {
	replay := newReplay(t, goodDir)
	recaps, err := replay.Blocks(2, criticalBlock+1)
	require.NoError(t, err)
	for _, recap := range recaps {
		fmt.Println(recap.String())
	}
}

func TestReplay_Bad(t *testing.T) {
	replay := newReplay(t, badDir)
	recaps, err := replay.Blocks(1, criticalBlock+1)
	require.NoError(t, err)
	for _, recap := range recaps {
		fmt.Println(recap.String())
	}
}

func TestStateHashes_Bad(t *testing.T) {
	badReplay := newReplay(t, badDir)
	goodReplay := newReplay(t, goodDir)
	for i := uint64(0); i <= criticalBlock+1; i++ {
		fmt.Println("Good")
		goodSt, err := goodReplay.State(i)
		require.NoError(t, err)
		fmt.Printf("Good: Version: %d, Hash: %X\n", goodSt.Version(), goodSt.Hash())
		fmt.Println("Bad")
		badSt, err := badReplay.State(i)
		require.NoError(t, err)
		fmt.Printf("Bad: Version: %d, Hash: %X\n", badSt.Version(), badSt.Hash())
		fmt.Println()
	}
}

func TestReplay_Good_Block(t *testing.T) {
	replayBlock(t, goodDir, criticalBlock)
}

func TestReplay_Bad_Block(t *testing.T) {
	replayBlock(t, badDir, criticalBlock)
}

func TestCriticalBlock(t *testing.T) {
	badState := getState(t, badDir, criticalBlock)
	goodState := getState(t, goodDir, criticalBlock)
	require.Equal(t, goodState.Hash(), badState.Hash())
	fmt.Printf("good: %X, bad: %X\n", goodState.Hash(), badState.Hash())
	_, _, err := badState.Update(func(up state.Updatable) error {
		return nil
	})
	require.NoError(t, err)
	_, _, err = goodState.Update(func(up state.Updatable) error {
		return nil
	})
	require.NoError(t, err)

	fmt.Printf("good: %X, bad: %X\n", goodState.Hash(), badState.Hash())
}

func replayBlock(t *testing.T, burrowDir string, height uint64) {
	replay := newReplay(t, burrowDir)
	//replay.State()
	recap, err := replay.Block(height)
	require.NoError(t, err)
	recap.TxExecutions = nil
	fmt.Println(recap)
}

func getState(t *testing.T, burrowDir string, height uint64) *state.State {
	st, err := newReplay(t, burrowDir).State(height)
	require.NoError(t, err)
	return st
}

func newReplay(t *testing.T, burrowDir string) *Replay {
	genesisDoc := new(genesis.GenesisDoc)
	err := source.FromFile(path.Join(burrowDir, "genesis.json"), genesisDoc)
	require.NoError(t, err)
	return NewReplay(path.Join(burrowDir, ".burrow"), genesisDoc, logging.NewNoopLogger())
}
