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

package evm

import (
	"fmt"
	"time"

	"github.com/hyperledger/burrow/account"
	. "github.com/hyperledger/burrow/word"

	"github.com/hyperledger/burrow/txs"
	"github.com/tendermint/go-wire"
	tm_types "github.com/tendermint/tendermint/types" // Block
)

// Functions to generate eventId strings

func EventStringAccInput(addr account.Address) string  { return fmt.Sprintf("Acc/%s/Input", addr) }
func EventStringAccOutput(addr account.Address) string { return fmt.Sprintf("Acc/%s/Output", addr) }
func EventStringAccCall(addr account.Address) string   { return fmt.Sprintf("Acc/%s/Call", addr) }
func EventStringLogEvent(addr account.Address) string  { return fmt.Sprintf("Log/%s", addr) }
func EventStringPermissions(name string) string        { return fmt.Sprintf("Permissions/%s", name) }
func EventStringNameReg(name string) string            { return fmt.Sprintf("NameReg/%s", name) }
func EventStringBond() string                          { return "Bond" }
func EventStringUnbond() string                        { return "Unbond" }
func EventStringRebond() string                        { return "Rebond" }
func EventStringDupeout() string                       { return "Dupeout" }
func EventStringNewBlock() string                      { return "NewBlock" }
func EventStringFork() string                          { return "Fork" }

func EventStringNewRound() string         { return fmt.Sprintf("NewRound") }
func EventStringTimeoutPropose() string   { return fmt.Sprintf("TimeoutPropose") }
func EventStringCompleteProposal() string { return fmt.Sprintf("CompleteProposal") }
func EventStringPolka() string            { return fmt.Sprintf("Polka") }
func EventStringUnlock() string           { return fmt.Sprintf("Unlock") }
func EventStringLock() string             { return fmt.Sprintf("Lock") }
func EventStringRelock() string           { return fmt.Sprintf("Relock") }
func EventStringTimeoutWait() string      { return fmt.Sprintf("TimeoutWait") }
func EventStringVote() string             { return fmt.Sprintf("Vote") }

//----------------------------------------

const (
	EventDataTypeNewBlock       = byte(0x01)
	EventDataTypeFork           = byte(0x02)
	EventDataTypeTx             = byte(0x03)
	EventDataTypeCall           = byte(0x04)
	EventDataTypeLog            = byte(0x05)
	EventDataTypeNewBlockHeader = byte(0x06)

	EventDataTypeRoundState = byte(0x11)
	EventDataTypeVote       = byte(0x12)
)

type EventData interface {
	AssertIsEventData()
}

var _ = wire.RegisterInterface(
	struct{ EventData }{},
	wire.ConcreteType{EventDataNewBlockHeader{}, EventDataTypeNewBlockHeader},
	wire.ConcreteType{EventDataNewBlock{}, EventDataTypeNewBlock},
	// wire.ConcreteType{EventDataFork{}, EventDataTypeFork },
	wire.ConcreteType{EventDataTx{}, EventDataTypeTx},
	wire.ConcreteType{EventDataCall{}, EventDataTypeCall},
	wire.ConcreteType{EventDataLog{}, EventDataTypeLog},
	wire.ConcreteType{EventDataRoundState{}, EventDataTypeRoundState},
	wire.ConcreteType{EventDataVote{}, EventDataTypeVote},
)

// Most event messages are basic types (a block, a transaction)
// but some (an input to a call tx or a receive) are more exotic

type EventDataNewBlock struct {
	Block *tm_types.Block `json:"block"`
}

type EventDataNewBlockHeader struct {
	Header *tm_types.Header `json:"header"`
}

// All txs fire EventDataTx, but only CallTx might have Return or Exception
type EventDataTx struct {
	Tx        txs.Tx `json:"tx"`
	Return    []byte `json:"return"`
	Exception string `json:"exception"`
}

// EventDataCall fires when we call a contract, and when a contract calls another contract
type EventDataCall struct {
	CallData  *CallData `json:"call_data"`
	Origin    []byte    `json:"origin"`
	TxID      []byte    `json:"tx_id"`
	Return    []byte    `json:"return"`
	Exception string    `json:"exception"`
}

type CallData struct {
	Caller []byte `json:"caller"`
	Callee []byte `json:"callee"`
	Data   []byte `json:"data"`
	Value  uint64  `json:"value"`
	Gas    uint64  `json:"gas"`
}

// EventDataLog fires when a contract executes the LOG opcode
type EventDataLog struct {
	Address Word256   `json:"address"`
	Topics  []Word256 `json:"topics"`
	Data    []byte    `json:"data"`
	Height  uint64     `json:"height"`
}

// We fire the most recent round state that led to the event
// (ie. NewRound will have the previous rounds state)
type EventDataRoundState struct {
	CurrentTime time.Time `json:"current_time"`

	Height        int                `json:"height"`
	Round         int                `json:"round"`
	Step          string             `json:"step"`
	StartTime     time.Time          `json:"start_time"`
	CommitTime    time.Time          `json:"commit_time"`
	Proposal      *tm_types.Proposal `json:"proposal"`
	ProposalBlock *tm_types.Block    `json:"proposal_block"`
	LockedRound   int                `json:"locked_round"`
	LockedBlock   *tm_types.Block    `json:"locked_block"`
	POLRound      int                `json:"pol_round"`
}

type EventDataVote struct {
	Index   int
	Address account.Address
	Vote    *tm_types.Vote
}

func (_ EventDataNewBlock) AssertIsEventData()       {}
func (_ EventDataNewBlockHeader) AssertIsEventData() {}
func (_ EventDataTx) AssertIsEventData()             {}
func (_ EventDataCall) AssertIsEventData()           {}
func (_ EventDataLog) AssertIsEventData()            {}
func (_ EventDataRoundState) AssertIsEventData()     {}
func (_ EventDataVote) AssertIsEventData()           {}
