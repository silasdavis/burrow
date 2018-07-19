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

package rpc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/state"
	"github.com/hyperledger/burrow/binary"
	bcm "github.com/hyperledger/burrow/blockchain"
	"github.com/hyperledger/burrow/consensus/tendermint/query"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution"
	"github.com/hyperledger/burrow/execution/names"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/logging/structure"
	"github.com/hyperledger/burrow/permission"
	"github.com/hyperledger/burrow/project"
	"github.com/hyperledger/burrow/txs"
	tmTypes "github.com/tendermint/tendermint/types"
)

// Magic! Should probably be configurable, but not shouldn't be so huge we
// end up DoSing ourselves.
const MaxBlockLookback = 1000

// Base service that provides implementation for all underlying RPC methods
type Service struct {
	state      state.IterableReader
	nameReg    names.IterableReader
	blockchain bcm.BlockchainInfo
	transactor *execution.Transactor
	nodeView   *query.NodeView
	logger     *logging.Logger
}

func NewService(state state.IterableReader, nameReg names.IterableReader,
	blockchain bcm.BlockchainInfo, transactor *execution.Transactor, nodeView *query.NodeView,
	logger *logging.Logger) *Service {

	return &Service{
		state:      state,
		nameReg:    nameReg,
		blockchain: blockchain,
		transactor: transactor,
		nodeView:   nodeView,
		logger:     logger.With(structure.ComponentKey, "Service"),
	}
}

// Get a Transactor providing methods for delegating signing and the core BroadcastTx function for publishing
// transactions to the network
func (s *Service) Transactor() *execution.Transactor {
	return s.transactor
}

func (s *Service) State() state.Reader {
	return s.state
}

func (s *Service) BlockchainInfo() bcm.BlockchainInfo {
	return s.blockchain
}

func (s *Service) ChainID() string {
	return s.blockchain.ChainID()
}

func (s *Service) ListUnconfirmedTxs(maxTxs int) (*ResultListUnconfirmedTxs, error) {
	// Get all transactions for now
	transactions, err := s.nodeView.MempoolTransactions(maxTxs)
	if err != nil {
		return nil, err
	}
	wrappedTxs := make([]*txs.Envelope, len(transactions))
	for i, tx := range transactions {
		wrappedTxs[i] = tx
	}
	return &ResultListUnconfirmedTxs{
		NumTxs: len(transactions),
		Txs:    wrappedTxs,
	}, nil
}

func (s *Service) Status() (*ResultStatus, error) {
	latestHeight := s.blockchain.LastBlockHeight()
	var (
		latestBlockMeta *tmTypes.BlockMeta
		latestBlockHash []byte
		latestBlockTime int64
	)
	if latestHeight != 0 {
		latestBlockMeta = s.nodeView.BlockStore().LoadBlockMeta(int64(latestHeight))
		latestBlockHash = latestBlockMeta.Header.Hash()
		latestBlockTime = latestBlockMeta.Header.Time.UnixNano()
	}
	publicKey, err := s.nodeView.PrivValidatorPublicKey()
	if err != nil {
		return nil, err
	}
	return &ResultStatus{
		NodeInfo:          s.nodeView.NodeInfo(),
		GenesisHash:       s.blockchain.GenesisHash(),
		PubKey:            publicKey,
		LatestBlockHash:   latestBlockHash,
		LatestBlockHeight: latestHeight,
		LatestBlockTime:   latestBlockTime,
		NodeVersion:       project.History.CurrentVersion().String(),
	}, nil
}

func (s *Service) ChainIdentifiers() (*ResultChainId, error) {
	return &ResultChainId{
		ChainName:   s.blockchain.GenesisDoc().ChainName,
		ChainId:     s.blockchain.ChainID(),
		GenesisHash: s.blockchain.GenesisHash(),
	}, nil
}

func (s *Service) Peers() (*ResultPeers, error) {
	peers := make([]*Peer, s.nodeView.Peers().Size())
	for i, peer := range s.nodeView.Peers().List() {
		peers[i] = &Peer{
			NodeInfo:   peer.NodeInfo(),
			IsOutbound: peer.IsOutbound(),
		}
	}
	return &ResultPeers{
		Peers: peers,
	}, nil
}

func (s *Service) NetInfo() (*ResultNetInfo, error) {
	listening := s.nodeView.IsListening()
	var listeners []string
	for _, listener := range s.nodeView.Listeners() {
		listeners = append(listeners, listener.String())
	}
	peers, err := s.Peers()
	if err != nil {
		return nil, err
	}
	return &ResultNetInfo{
		Listening: listening,
		Listeners: listeners,
		Peers:     peers.Peers,
	}, nil
}

func (s *Service) Genesis() (*ResultGenesis, error) {
	return &ResultGenesis{
		Genesis: s.blockchain.GenesisDoc(),
	}, nil
}

// Accounts
func (s *Service) GetAccount(address crypto.Address) (*ResultGetAccount, error) {
	acc, err := s.state.GetAccount(address)
	if err != nil {
		return nil, err
	}
	return &ResultGetAccount{Account: acm.AsConcreteAccount(acc)}, nil
}

func (s *Service) ListAccounts(predicate func(acm.Account) bool) (*ResultListAccounts, error) {
	accounts := make([]*acm.ConcreteAccount, 0)
	s.state.IterateAccounts(func(account acm.Account) (stop bool) {
		if predicate(account) {
			accounts = append(accounts, acm.AsConcreteAccount(account))
		}
		return
	})

	return &ResultListAccounts{
		BlockHeight: s.blockchain.LastBlockHeight(),
		Accounts:    accounts,
	}, nil
}

func (s *Service) GetStorage(address crypto.Address, key []byte) (*ResultGetStorage, error) {
	account, err := s.state.GetAccount(address)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("UnknownAddress: %s", address)
	}

	value, err := s.state.GetStorage(address, binary.LeftPadWord256(key))
	if err != nil {
		return nil, err
	}
	if value == binary.Zero256 {
		return &ResultGetStorage{Key: key, Value: nil}, nil
	}
	return &ResultGetStorage{Key: key, Value: value.UnpadLeft()}, nil
}

func (s *Service) DumpStorage(address crypto.Address) (*ResultDumpStorage, error) {
	account, err := s.state.GetAccount(address)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("UnknownAddress: %X", address)
	}
	var storageItems []StorageItem
	s.state.IterateStorage(address, func(key, value binary.Word256) (stop bool) {
		storageItems = append(storageItems, StorageItem{Key: key.UnpadLeft(), Value: value.UnpadLeft()})
		return
	})
	return &ResultDumpStorage{
		StorageItems: storageItems,
	}, nil
}

func (s *Service) GetAccountHumanReadable(address crypto.Address) (*ResultGetAccountHumanReadable, error) {
	acc, err := s.state.GetAccount(address)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return &ResultGetAccountHumanReadable{}, nil
	}
	tokens, err := acc.Code().Tokens()
	if acc == nil {
		return &ResultGetAccountHumanReadable{}, nil
	}
	perms, err := permission.BasePermissionsToStringList(acc.Permissions().Base)
	if acc == nil {
		return &ResultGetAccountHumanReadable{}, nil
	}
	return &ResultGetAccountHumanReadable{
		Account: &AccountHumanReadable{
			Address:     acc.Address(),
			PublicKey:   acc.PublicKey(),
			Sequence:    acc.Sequence(),
			Balance:     acc.Balance(),
			Code:        tokens,
			Permissions: perms,
			Roles:       acc.Permissions().Roles,
		},
	}, nil
}

// Name registry
func (s *Service) GetName(name string) (*ResultGetName, error) {
	entry, err := s.nameReg.GetName(name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("name %s not found", name)
	}
	return &ResultGetName{Entry: entry}, nil
}

func (s *Service) ListNames(predicate func(*names.Entry) bool) (*ResultListNames, error) {
	var nms []*names.Entry
	s.nameReg.IterateNames(func(entry *names.Entry) (stop bool) {
		if predicate(entry) {
			nms = append(nms, entry)
		}
		return
	})
	return &ResultListNames{
		BlockHeight: s.blockchain.LastBlockHeight(),
		Names:       nms,
	}, nil
}

func (s *Service) GetBlock(height uint64) (*ResultGetBlock, error) {
	return &ResultGetBlock{
		Block:     &Block{s.nodeView.BlockStore().LoadBlock(int64(height))},
		BlockMeta: &BlockMeta{s.nodeView.BlockStore().LoadBlockMeta(int64(height))},
	}, nil
}

// Returns the current blockchain height and metadata for a range of blocks
// between minHeight and maxHeight. Only returns maxBlockLookback block metadata
// from the top of the range of blocks.
// Passing 0 for maxHeight sets the upper height of the range to the current
// blockchain height.
func (s *Service) ListBlocks(minHeight, maxHeight uint64) (*ResultListBlocks, error) {
	latestHeight := s.blockchain.LastBlockHeight()

	if minHeight == 0 {
		minHeight = 1
	}
	if maxHeight == 0 || latestHeight < maxHeight {
		maxHeight = latestHeight
	}
	if maxHeight > minHeight && maxHeight-minHeight > MaxBlockLookback {
		minHeight = maxHeight - MaxBlockLookback
	}

	var blockMetas []*tmTypes.BlockMeta
	for height := maxHeight; height >= minHeight; height-- {
		blockMeta := s.nodeView.BlockStore().LoadBlockMeta(int64(height))
		blockMetas = append(blockMetas, blockMeta)
	}

	return &ResultListBlocks{
		LastHeight: latestHeight,
		BlockMetas: blockMetas,
	}, nil
}

func (s *Service) ListValidators() (*ResultListValidators, error) {
	concreteValidators := make([]*acm.ConcreteValidator, 0, s.blockchain.NumValidators())
	s.blockchain.IterateValidators(func(id crypto.Addressable, power *big.Int) (stop bool) {
		concreteValidators = append(concreteValidators, &acm.ConcreteValidator{
			Address:   id.Address(),
			PublicKey: id.PublicKey(),
			Power:     power.Uint64(),
		})
		return
	})
	return &ResultListValidators{
		BlockHeight:         s.blockchain.LastBlockHeight(),
		BondedValidators:    concreteValidators,
		UnbondingValidators: nil,
	}, nil
}

func (s *Service) DumpConsensusState() (*ResultDumpConsensusState, error) {
	peerRoundState, err := s.nodeView.PeerRoundStates()
	if err != nil {
		return nil, err
	}
	return &ResultDumpConsensusState{
		RoundState:      s.nodeView.RoundState().RoundStateSimple(),
		PeerRoundStates: peerRoundState,
	}, nil
}

func (s *Service) GeneratePrivateAccount() (*ResultGeneratePrivateAccount, error) {
	privateAccount, err := acm.GeneratePrivateAccount()
	if err != nil {
		return nil, err
	}
	return &ResultGeneratePrivateAccount{
		PrivateAccount: privateAccount.ConcretePrivateAccount(),
	}, nil
}

func (s *Service) LastBlockInfo(blockWithin string) (*ResultLastBlockInfo, error) {
	res := &ResultLastBlockInfo{
		LastBlockHeight: s.blockchain.LastBlockHeight(),
		LastBlockHash:   s.blockchain.LastBlockHash(),
		LastBlockTime:   s.blockchain.LastBlockTime(),
	}
	if blockWithin == "" {
		return res, nil
	}
	duration, err := time.ParseDuration(blockWithin)
	if err != nil {
		return nil, fmt.Errorf("could not parse blockWithin duration to determine whether to throw error: %v", err)
	}
	// Take neg abs in case caller is counting backwards (not we add later)
	if duration > 0 {
		duration = -duration
	}
	blockTimeThreshold := time.Now().Add(duration)
	if res.LastBlockTime.After(blockTimeThreshold) {
		// We've created blocks recently enough
		return res, nil
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		resJSON = []byte("<error: could not marshal last block info>")
	}
	return nil, fmt.Errorf("no block committed within the last %s (cutoff: %s), last block info: %s",
		blockWithin, blockTimeThreshold.Format(time.RFC3339), string(resJSON))
}
