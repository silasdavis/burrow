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

package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"sync"

	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/genesis"
	"github.com/hyperledger/burrow/logging"
	dbm "github.com/tendermint/tmlibs/db"
)

// Blocks to average validator power over
const DefaultValidatorsWindowSize = 10

var stateKey = []byte("BlockchainState")

type TipInfo interface {
	ChainID() string
	LastBlockHeight() uint64
	LastBlockTime() time.Time
	LastBlockHash() []byte
	AppHashAfterLastBlock() []byte
	IterateValidators(iter func(id crypto.Addressable, power *big.Int) (stop bool)) (stopped bool)
	NumValidators() int
}

type BlockchainInfo interface {
	TipInfo
	GenesisHash() []byte
	GenesisDoc() genesis.GenesisDoc
}

type Root struct {
	genesisHash []byte
	genesisDoc  genesis.GenesisDoc
}

type Tip struct {
	chainID               string
	lastBlockHeight       uint64
	lastBlockTime         time.Time
	lastBlockHash         []byte
	appHashAfterLastBlock []byte
	validators            *Validators
	validatorsRing        *ValidatorsRing
}

type Blockchain struct {
	*Root
	*Tip
	sync.RWMutex
	db dbm.DB
}

var _ TipInfo = &Blockchain{}

type PersistedState struct {
	AppHashAfterLastBlock []byte
	LastBlockHeight       uint64
	GenesisDoc            genesis.GenesisDoc
	Validators            []PersistedValidator
}

func LoadOrNewBlockchain(db dbm.DB, genesisDoc *genesis.GenesisDoc,
	logger *logging.Logger) (*Blockchain, error) {

	logger = logger.WithScope("LoadOrNewBlockchain")
	logger.InfoMsg("Trying to load blockchain state from database",
		"database_key", stateKey)
	bc, err := loadBlockchain(db)
	if err != nil {
		return nil, fmt.Errorf("error loading blockchain state from database: %v", err)
	}
	if bc != nil {
		dbHash := bc.genesisDoc.Hash()
		argHash := genesisDoc.Hash()
		if !bytes.Equal(dbHash, argHash) {
			return nil, fmt.Errorf("GenesisDoc passed to LoadOrNewBlockchain has hash: 0x%X, which does not "+
				"match the one found in database: 0x%X", argHash, dbHash)
		}
		return bc, nil
	}

	logger.InfoMsg("No existing blockchain state found in database, making new blockchain")
	return newBlockchain(db, genesisDoc), nil
}

// Pointer to blockchain state initialised from genesis
func newBlockchain(db dbm.DB, genesisDoc *genesis.GenesisDoc) *Blockchain {
	vs := NewValidators()
	for _, gv := range genesisDoc.Validators {
		vs.AlterPower(gv.PublicKey, new(big.Int).SetUint64(gv.Amount))
	}
	root := NewRoot(genesisDoc)
	bc := &Blockchain{
		db:   db,
		Root: root,
		Tip:  NewTip(genesisDoc.ChainID(), root.genesisDoc.GenesisTime, root.genesisHash, vs),
	}
	return bc
}

func loadBlockchain(db dbm.DB) (*Blockchain, error) {
	buf := db.Get(stateKey)
	if len(buf) == 0 {
		return nil, nil
	}
	persistedState, err := Decode(buf)
	if err != nil {
		return nil, err
	}
	bc := newBlockchain(db, &persistedState.GenesisDoc)
	bc.lastBlockHeight = persistedState.LastBlockHeight
	bc.appHashAfterLastBlock = persistedState.AppHashAfterLastBlock
	return bc, nil
}

func NewRoot(genesisDoc *genesis.GenesisDoc) *Root {
	return &Root{
		genesisHash: genesisDoc.Hash(),
		genesisDoc:  *genesisDoc,
	}
}

// Create genesis Tip
func NewTip(chainID string, genesisTime time.Time, genesisHash []byte, initialValidators *Validators) *Tip {
	return &Tip{
		chainID:               chainID,
		lastBlockTime:         genesisTime,
		appHashAfterLastBlock: genesisHash,
		validators:            initialValidators,
		validatorsRing:        NewValidatorsRing(initialValidators, DefaultValidatorsWindowSize),
	}
}

//func (bc *Blockchain) FlushValidators() (powerChange *big.Int, totalFlow *big.Int, err error) {
//	return bc.validatorsWindow.FlushInto(bc.validators)
//}

func (bc *Blockchain) AlterPower(id crypto.Addressable, power *big.Int) (*big.Int, error) {
	return bc.validatorsRing.AlterPower(id, power)
}

func (bc *Blockchain) CommitBlock(blockTime time.Time, blockHash, appHash []byte) error {
	bc.Lock()
	defer bc.Unlock()
	// Checkpoint on the _previous_ block. If we die, this is where we will resume since we know it must have been
	// committed since we are committing the next block. If we fall over we can resume a safe committed state and
	// Tendermint will catch us up
	err := bc.save()
	if err != nil {
		return err
	}
	bc.lastBlockHeight += 1
	bc.lastBlockTime = blockTime
	bc.lastBlockHash = blockHash
	bc.appHashAfterLastBlock = appHash
	return nil
}

func (bc *Blockchain) save() error {
	if bc.db != nil {
		encodedState, err := bc.Encode()
		if err != nil {
			return err
		}
		bc.db.SetSync(stateKey, encodedState)
	}
	return nil
}

// TODO: this will need to be amino
func (bc *Blockchain) Encode() ([]byte, error) {
	persistedState := &PersistedState{
		GenesisDoc:            bc.genesisDoc,
		AppHashAfterLastBlock: bc.appHashAfterLastBlock,
		LastBlockHeight:       bc.lastBlockHeight,
	}
	encodedState, err := json.Marshal(persistedState)
	if err != nil {
		return nil, err
	}
	return encodedState, nil
}

func Decode(encodedState []byte) (*PersistedState, error) {
	persistedState := new(PersistedState)
	err := json.Unmarshal(encodedState, persistedState)
	if err != nil {
		return nil, err
	}
	return persistedState, nil
}

func (r *Root) GenesisHash() []byte {
	return r.genesisHash
}

func (r *Root) GenesisDoc() genesis.GenesisDoc {
	return r.genesisDoc
}

func (t *Tip) ChainID() string {
	return t.chainID
}

func (t *Tip) LastBlockHeight() uint64 {
	return t.lastBlockHeight
}

func (t *Tip) LastBlockTime() time.Time {
	return t.lastBlockTime
}

func (t *Tip) LastBlockHash() []byte {
	return t.lastBlockHash
}

func (t *Tip) AppHashAfterLastBlock() []byte {
	return t.appHashAfterLastBlock
}

func (t *Tip) IterateValidators(iter func(id crypto.Addressable, power *big.Int) (stop bool)) (stopped bool) {
	return t.validators.Iterate(iter)
}

func (t *Tip) NumValidators() int {
	return t.validators.Count()
}
