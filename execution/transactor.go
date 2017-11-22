// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package execution

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	acm "github.com/hyperledger/burrow/account"
	"github.com/hyperledger/burrow/blockchain"
	"github.com/hyperledger/burrow/event"
	"github.com/hyperledger/burrow/execution/evm"
	"github.com/hyperledger/burrow/txs"
	"github.com/hyperledger/burrow/word"
	abci_types "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
)

type Call struct {
	Return  string `json:"return"`
	GasUsed int64  `json:"gas_used"`
	// TODO ...
}

type Transactor interface {
	Call(fromAddress, toAddress acm.Address, data []byte) (*Call, error)
	CallCode(fromAddress acm.Address, code, data []byte) (*Call, error)
	BroadcastTx(tx txs.Tx) (*txs.Receipt, error)
	BroadcastTxAsync(tx txs.Tx, callback func(res *abci_types.Response)) error
	Transact(privKey []byte, address acm.Address, data []byte, gasLimit, fee int64) (*txs.Receipt, error)
	TransactAndHold(privKey []byte, address acm.Address, data []byte, gasLimit, fee int64) (*evm.EventDataCall, error)
	Send(privKey []byte, toAddress acm.Address, amount int64) (*txs.Receipt, error)
	SendAndHold(privKey []byte, toAddress acm.Address, amount int64) (*txs.Receipt, error)
	TransactNameReg(privKey []byte, name, data string, amount, fee int64) (*txs.Receipt, error)
	SignTx(tx txs.Tx, privAccounts []acm.PrivateAccount) (txs.Tx, error)
}

// Transactor is the controller/middleware for the v0 RPC
type transactor struct {
	txMtx            *sync.Mutex
	blockchain       blockchain.Blockchain
	state            acm.StateReader
	eventEmitter     event.EventEmitter
	broadcastTxAsync func(tx txs.Tx, callback func(res *abci_types.Response)) error
}

var _ Transactor = &transactor{}

func NewTransactor(blockchain blockchain.Blockchain, state acm.StateReader, eventEmitter event.EventEmitter,
	broadcastTxAsync func(tx txs.Tx, callback func(res *abci_types.Response)) error) *transactor {
	return &transactor{
		blockchain:       blockchain,
		state:            state,
		eventEmitter:     eventEmitter,
		broadcastTxAsync: broadcastTxAsync,
	}
}

// Run a contract's code on an isolated and unpersisted state
// Cannot be used to create new contracts
func (trans *transactor) Call(fromAddress, toAddress acm.Address, data []byte) (*Call, error) {
	if evm.RegisteredNativeContract(toAddress.Word256()) {
		return nil, fmt.Errorf("attempt to call native contract at address "+
			"%X, but native contracts can not be called directly. Use a deployed "+
			"contract that calls the native function instead", toAddress)
	}
	// This was being run against CheckTx cache, need to understand the reasoning
	callee := acm.AsMutableAccount(trans.state.GetAccount(toAddress))
	if callee == nil {
		return nil, fmt.Errorf("account %s does not exist", toAddress)
	}
	caller := acm.ConcreteAccount{Address: fromAddress}.MutableAccount()
	txCache := NewTxCache(trans.state)
	params := vmParams(trans.blockchain)

	vmach := evm.NewVM(txCache, evm.DefaultDynamicMemoryProvider, params, caller.Address().Word256(), nil)
	vmach.SetFireable(trans.eventEmitter)

	gas := params.GasLimit
	ret, err := vmach.Call(caller, callee, callee.Code(), data, 0, &gas)
	if err != nil {
		return nil, err
	}
	gasUsed := params.GasLimit - gas
	// here return bytes are hex encoded; on the sibling function
	// they are not
	return &Call{Return: hex.EncodeToString(ret), GasUsed: gasUsed}, nil
}

// Run the given code on an isolated and unpersisted state
// Cannot be used to create new contracts.
func (trans *transactor) CallCode(fromAddress acm.Address, code, data []byte) (*Call, error) {
	// This was being run against CheckTx cache, need to understand the reasoning
	callee := acm.ConcreteAccount{Address: fromAddress}.MutableAccount()
	caller := acm.ConcreteAccount{Address: fromAddress}.MutableAccount()
	txCache := NewTxCache(trans.state)
	params := vmParams(trans.blockchain)

	vmach := evm.NewVM(txCache, evm.DefaultDynamicMemoryProvider, params, caller.Address().Word256(), nil)
	gas := params.GasLimit
	ret, err := vmach.Call(caller, callee, code, data, 0, &gas)
	if err != nil {
		return nil, err
	}
	gasUsed := params.GasLimit - gas
	// here return bytes are hex encoded; on the sibling function
	// they are not
	return &Call{Return: hex.EncodeToString(ret), GasUsed: gasUsed}, nil
}

func (trans *transactor) BroadcastTxAsync(tx txs.Tx, callback func(res *abci_types.Response)) error {
	return trans.broadcastTxAsync(tx, callback)
}

// Broadcast a transaction.
func (trans *transactor) BroadcastTx(tx txs.Tx) (*txs.Receipt, error) {
	responseCh := make(chan *abci_types.Response, 1)
	err := trans.BroadcastTxAsync(tx, func(res *abci_types.Response) {
		responseCh <- res
	})

	if err != nil {
		return nil, err
	}
	response := <-responseCh
	checkTxResponse := response.GetCheckTx()
	if checkTxResponse == nil {
		return nil, fmt.Errorf("application did not return CheckTx response")
	}

	switch checkTxResponse.Code {
	case abci_types.CodeType_OK:
		receipt := new(txs.Receipt)
		err := wire.ReadBinaryBytes(checkTxResponse.Data, receipt)
		if err != nil {
			return nil, fmt.Errorf("could not deserialise transaction receipt: %s", err)
		}
		return receipt, nil
	case abci_types.CodeType_EncodingError, abci_types.CodeType_InternalError:
		return nil, fmt.Errorf("error code %s received, log: %s", checkTxResponse.Code, checkTxResponse.Log)
	default:
		return nil, fmt.Errorf("unknown error returned from Tendermint by BroadcastTxSync "+
			"ABCI code: %v, ABCI log: %v", checkTxResponse.Code, checkTxResponse.Log)
	}
}

// Orders calls to BroadcastTx using lock (waits for response from core before releasing)
func (trans *transactor) Transact(privKey []byte, address acm.Address, data []byte, gasLimit,
	fee int64) (*txs.Receipt, error) {
	if len(privKey) != 64 {
		return nil, fmt.Errorf("Private key is not of the right length: %d\n", len(privKey))
	}
	trans.txMtx.Lock()
	defer trans.txMtx.Unlock()
	pa := acm.GeneratePrivateAccountFromPrivateKeyBytes(privKey)
	// [Silas] This is puzzling, if the account doesn't exist the CallTx will fail, so what's the point in this?
	acc := trans.state.GetAccount(pa.Address())
	sequence := int64(1)
	if acc != nil {
		sequence = acc.Sequence() + 1
	}
	// TODO: [Silas] we should consider revising this method and removing fee, or
	// possibly adding an amount parameter. It is non-sensical to just be able to
	// set the fee. Our support of fees in general is questionable since at the
	// moment all we do is deduct the fee effectively leaking token. It is possible
	// someone may be using the sending of native token to payable functions but
	// they can be served by broadcasting a token.

	// We hard-code the amount to be equal to the fee which means the CallTx we
	// generate transfers 0 value, which is the most sensible default since in
	// recent solidity compilers the EVM generated will throw an error if value
	// is transferred to a non-payable function.
	txInput := &txs.TxInput{
		Address:  pa.Address(),
		Amount:   fee,
		Sequence: sequence,
		PubKey:   pa.PubKey(),
	}
	tx := &txs.CallTx{
		Input:    txInput,
		Address:  &address,
		GasLimit: gasLimit,
		Fee:      fee,
		Data:     data,
	}

	// Got ourselves a tx.
	txS, errS := trans.SignTx(tx, []acm.PrivateAccount{pa})
	if errS != nil {
		return nil, errS
	}
	return trans.BroadcastTx(txS)
}

func (trans *transactor) TransactAndHold(privKey []byte, address acm.Address, data []byte, gasLimit,
	fee int64) (*evm.EventDataCall, error) {
	rec, tErr := trans.Transact(privKey, address, data, gasLimit, fee)
	if tErr != nil {
		return nil, tErr
	}
	var addr acm.Address
	if rec.CreatesContract {
		addr = rec.ContractAddr
	} else {
		addr = address
	}
	// We want non-blocking on the first event received (but buffer the value),
	// after which we want to block (and then discard the value - see below)
	wc := make(chan *evm.EventDataCall, 1)
	subId := fmt.Sprintf("%X", rec.TxHash)
	trans.eventEmitter.Subscribe(subId, evm.EventStringAccCall(addr),
		func(evt evm.EventData) {
			eventDataCall := evt.(evm.EventDataCall)
			if bytes.Equal(eventDataCall.TxID, rec.TxHash) {
				// Beware the contract of go-events subscribe is that we must not be
				// blocking in an event callback when we try to unsubscribe!
				// We work around this by using a non-blocking send.
				select {
				// This is a non-blocking send, but since we are using a buffered
				// channel of size 1 we will always grab our first event even if we
				// haven't read from the channel at the time we receive the first event.
				case wc <- &eventDataCall:
				default:
				}
			}
		})

	timer := time.NewTimer(300 * time.Second)
	toChan := timer.C

	var ret *evm.EventDataCall
	var rErr error

	select {
	case <-toChan:
		rErr = fmt.Errorf("Transaction timed out. Hash: " + subId)
	case e := <-wc:
		timer.Stop()
		if e.Exception != "" {
			rErr = fmt.Errorf("error when transacting: " + e.Exception)
		} else {
			ret = e
		}
	}
	trans.eventEmitter.Unsubscribe(subId)
	return ret, rErr
}

func (trans *transactor) Send(privKey []byte, toAddress acm.Address, amount int64) (*txs.Receipt, error) {
	if len(privKey) != 64 {
		return nil, fmt.Errorf("Private key is not of the right length: %d\n",
			len(privKey))
	}

	pk := &[64]byte{}
	copy(pk[:], privKey)
	trans.txMtx.Lock()
	defer trans.txMtx.Unlock()
	pa := acm.GeneratePrivateAccountFromPrivateKeyBytes(privKey)
	cache := trans.state
	acc := cache.GetAccount(pa.Address())
	sequence := int64(1)
	if acc != nil {
		sequence = acc.Sequence() + 1
	}

	tx := txs.NewSendTx()

	txInput := &txs.TxInput{
		Address:  pa.Address(),
		Amount:   amount,
		Sequence: sequence,
		PubKey:   pa.PubKey(),
	}

	tx.Inputs = append(tx.Inputs, txInput)

	txOutput := &txs.TxOutput{Address: toAddress, Amount: amount}

	tx.Outputs = append(tx.Outputs, txOutput)

	// Got ourselves a tx.
	txS, errS := trans.SignTx(tx, []acm.PrivateAccount{pa})
	if errS != nil {
		return nil, errS
	}
	return trans.BroadcastTx(txS)
}

func (trans *transactor) SendAndHold(privKey []byte, toAddress acm.Address,
	amount int64) (*txs.Receipt, error) {
	rec, tErr := trans.Send(privKey, toAddress, amount)
	if tErr != nil {
		return nil, tErr
	}

	wc := make(chan *txs.SendTx)
	subId := fmt.Sprintf("%X", rec.TxHash)

	trans.eventEmitter.Subscribe(subId, evm.EventStringAccOutput(toAddress),
		func(evt evm.EventData) {
			eventDataTx := evt.(evm.EventDataTx)
			tx := eventDataTx.Tx.(*txs.SendTx)
			wc <- tx
		})

	timer := time.NewTimer(300 * time.Second)
	toChan := timer.C

	var rErr error

	pa := acm.GeneratePrivateAccountFromPrivateKeyBytes(privKey)

	select {
	case <-toChan:
		rErr = fmt.Errorf("Transaction timed out. Hash: " + subId)
	case e := <-wc:
		if e.Inputs[0].Address == pa.Address() && e.Inputs[0].Amount == amount {
			timer.Stop()
			trans.eventEmitter.Unsubscribe(subId)
			return rec, rErr
		}
	}
	return nil, rErr
}

func (trans *transactor) TransactNameReg(privKey []byte, name, data string,
	amount, fee int64) (*txs.Receipt, error) {

	if len(privKey) != 64 {
		return nil, fmt.Errorf("Private key is not of the right length: %d\n", len(privKey))
	}
	trans.txMtx.Lock()
	defer trans.txMtx.Unlock()
	pa := acm.GeneratePrivateAccountFromPrivateKeyBytes(privKey)
	cache := trans.state // XXX: DON'T MUTATE THIS CACHE (used internally for CheckTx)
	acc := cache.GetAccount(pa.Address())
	sequence := int64(1)
	if acc == nil {
		sequence = acc.Sequence() + 1
	}
	tx := txs.NewNameTxWithNonce(pa.PubKey(), name, data, amount, fee, sequence)
	// Got ourselves a tx.
	txS, errS := trans.SignTx(tx, []acm.PrivateAccount{pa})
	if errS != nil {
		return nil, errS
	}
	return trans.BroadcastTx(txS)
}

// Sign a transaction
func (trans *transactor) SignTx(tx txs.Tx, privAccounts []acm.PrivateAccount) (txs.Tx, error) {
	// more checks?

	for i, privAccount := range privAccounts {
		if privAccount == nil || privAccount.PrivKey().Unwrap() == nil {
			return nil, fmt.Errorf("invalid (empty) privAccount @%v", i)
		}
	}
	chainID := trans.blockchain.Root().ChainID()
	switch tx.(type) {
	case *txs.NameTx:
		nameTx := tx.(*txs.NameTx)
		nameTx.Input.PubKey = privAccounts[0].PubKey()
		nameTx.Input.Signature = acm.ChainSign(privAccounts[0], chainID, nameTx)
	case *txs.SendTx:
		sendTx := tx.(*txs.SendTx)
		for i, input := range sendTx.Inputs {
			input.PubKey = privAccounts[i].PubKey()
			input.Signature = acm.ChainSign(privAccounts[i], chainID, sendTx)
		}
	case *txs.CallTx:
		callTx := tx.(*txs.CallTx)
		callTx.Input.PubKey = privAccounts[0].PubKey()
		callTx.Input.Signature = acm.ChainSign(privAccounts[0], chainID, callTx)
	case *txs.BondTx:
		bondTx := tx.(*txs.BondTx)
		// the first privaccount corresponds to the BondTx pub key.
		// the rest to the inputs
		bondTx.Signature = acm.ChainSign(privAccounts[0], chainID, bondTx).
			Unwrap().(crypto.SignatureEd25519)
		for i, input := range bondTx.Inputs {
			input.PubKey = privAccounts[i+1].PubKey()
			input.Signature = acm.ChainSign(privAccounts[i+1], chainID, bondTx)
		}
	case *txs.UnbondTx:
		unbondTx := tx.(*txs.UnbondTx)
		unbondTx.Signature = acm.ChainSign(privAccounts[0], chainID, unbondTx).
			Unwrap().(crypto.SignatureEd25519)
	case *txs.RebondTx:
		rebondTx := tx.(*txs.RebondTx)
		rebondTx.Signature = acm.ChainSign(privAccounts[0], chainID, rebondTx).
			Unwrap().(crypto.SignatureEd25519)
	default:
		return nil, fmt.Errorf("Object is not a proper transaction: %v\n", tx)
	}
	return tx, nil
}

func vmParams(blockchain blockchain.Blockchain) evm.Params {
	tip := blockchain.Tip()
	return evm.Params{
		BlockHeight: int64(tip.LastBlockHeight()),
		BlockHash:   word.LeftPadWord256(tip.LastBlockHash()),
		BlockTime:   tip.LastBlockTime().Unix(),
		GasLimit:    GasLimit,
	}
}
