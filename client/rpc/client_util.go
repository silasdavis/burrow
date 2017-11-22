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
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/tendermint/go-crypto"

	acm "github.com/hyperledger/burrow/account"
	"github.com/hyperledger/burrow/client"
	"github.com/hyperledger/burrow/keys"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/permission"
	ptypes "github.com/hyperledger/burrow/permission/types"
	"github.com/hyperledger/burrow/txs"
)

//------------------------------------------------------------------------------------
// sign and broadcast convenience

// tx has either one input or we default to the first one (ie for send/bond)
// TODO: better support for multisig and bonding
func signTx(keyClient keys.KeyClient, chainID string, tx_ txs.Tx) (acm.Address, txs.Tx, error) {
	signBytesString := fmt.Sprintf("%X", acm.SignBytes(chainID, tx_))
	var inputAddr acm.Address
	var sigED crypto.SignatureEd25519
	switch tx := tx_.(type) {
	case *txs.SendTx:
		inputAddr = tx.Inputs[0].Address
		defer func(s *crypto.SignatureEd25519) { tx.Inputs[0].Signature = s.Wrap() }(&sigED)
	case *txs.NameTx:
		inputAddr = tx.Input.Address
		defer func(s *crypto.SignatureEd25519) { tx.Input.Signature = s.Wrap() }(&sigED)
	case *txs.CallTx:
		inputAddr = tx.Input.Address
		defer func(s *crypto.SignatureEd25519) { tx.Input.Signature = s.Wrap() }(&sigED)
	case *txs.PermissionsTx:
		inputAddr = tx.Input.Address
		defer func(s *crypto.SignatureEd25519) { tx.Input.Signature = s.Wrap() }(&sigED)
	case *txs.BondTx:
		inputAddr = tx.Inputs[0].Address
		defer func(s *crypto.SignatureEd25519) {
			tx.Signature = *s
			tx.Inputs[0].Signature = s.Wrap()
		}(&sigED)
	case *txs.UnbondTx:
		inputAddr = tx.Address
		defer func(s *crypto.SignatureEd25519) { tx.Signature = *s }(&sigED)
	case *txs.RebondTx:
		inputAddr = tx.Address
		defer func(s *crypto.SignatureEd25519) { tx.Signature = *s }(&sigED)
	}
	sig, err := keyClient.Sign(signBytesString, inputAddr)
	if err != nil {
		return acm.Address{}, nil, err
	}
	// TODO: [ben] temporarily address the type conflict here, to be cleaned up
	// with full type restructuring
	var sig64 [64]byte
	copy(sig64[:], sig)
	sigED = crypto.SignatureEd25519(sig64)
	return inputAddr, tx_, nil
}

func decodeAddressPermFlag(addrS, permFlagS string) (addr acm.Address, pFlag ptypes.PermFlag, err error) {
	var addrBytes []byte
	if addrBytes, err = hex.DecodeString(addrS); err != nil {
		copy(addr[:], addrBytes)
		return
	}
	if pFlag, err = permission.PermStringToFlag(permFlagS); err != nil {
		return
	}
	return
}

func checkCommon(nodeClient client.NodeClient, keyClient keys.KeyClient, pubkey, addr, amtS,
	nonceS string) (pub crypto.PubKey, amt uint64, nonce uint64, err error) {

	if amtS == "" {
		err = fmt.Errorf("input must specify an amount with the --amt flag")
		return
	}

	var pubKeyBytes []byte
	if pubkey == "" && addr == "" {
		err = fmt.Errorf("at least one of --pubkey or --addr must be given")
		return
	} else if pubkey != "" {
		if addr != "" {
			logging.InfoMsg(nodeClient.Logger(), "Both a public key and an address have been specified. The public key takes precedent.",
				"public_key", pubkey,
				"address", addr,
			)
		}
		pubKeyBytes, err = hex.DecodeString(pubkey)
		if err != nil {
			err = fmt.Errorf("pubkey is bad hex: %v", err)
			return
		}
	} else {
		// grab the pubkey from monax-keys
		addressBytes, err2 := hex.DecodeString(addr)
		if err2 != nil {
			err = fmt.Errorf("Bad hex string for address (%s): %v", addr, err)
			return
		}
		address, err2 := acm.AddressFromBytes(addressBytes)
		if err2 != nil {
			err = fmt.Errorf("Could not convert bytes (%X) to address: %v", addressBytes, err2)
		}
		pubKeyBytes, err2 = keyClient.PublicKey(address)
		if err2 != nil {
			err = fmt.Errorf("Failed to fetch pubkey for address (%s): %v", addr, err2)
			return
		}
	}

	if len(pubKeyBytes) == 0 {
		err = fmt.Errorf("Error resolving public key")
		return
	}

	amt, err = strconv.ParseUint(amtS, 10, 64)
	if err != nil {
		err = fmt.Errorf("amt is misformatted: %v", err)
	}

	var pubArray [32]byte
	copy(pubArray[:], pubKeyBytes)
	pub = crypto.PubKeyEd25519(pubArray).Wrap()
	address, err := acm.AddressFromBytes(pub.Address())
	if err != nil {
		return
	}

	if nonceS == "" {
		if nodeClient == nil {
			err = fmt.Errorf("input must specify a nonce with the --nonce flag or use --node-addr (or BURROW_CLIENT_NODE_ADDR) to fetch the nonce from a node")
			return
		}
		// fetch nonce from node
		account, err2 := nodeClient.GetAccount(address)
		if err2 != nil {
			return pub, amt, nonce, err2
		}
		nonce = account.Sequence() + 1
		logging.TraceMsg(nodeClient.Logger(), "Fetch nonce from node",
			"nonce", nonce,
			"account address", address,
		)
	} else {
		nonce, err = strconv.ParseUint(nonceS, 10, 64)
		if err != nil {
			err = fmt.Errorf("nonce is misformatted: %v", err)
			return
		}
	}

	return
}
