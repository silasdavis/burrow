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

package mock

import (
	"crypto/rand"
	"fmt"

	acm "github.com/hyperledger/burrow/account"
	. "github.com/hyperledger/burrow/keys"
	"github.com/tendermint/ed25519"
	"github.com/tendermint/go-crypto"
	"golang.org/x/crypto/ripemd160"
)

//---------------------------------------------------------------------
// Mock ed25510 key for mock keys client

// Simple ed25519 key structure for mock purposes with ripemd160 address
type MockKey struct {
	Address    acm.Address
	PrivateKey [ed25519.PrivateKeySize]byte
	PublicKey  []byte
}

func newMockKey() (*MockKey, error) {
	key := &MockKey{
		PublicKey: make([]byte, ed25519.PublicKeySize),
	}
	// this is a mock key, so the entropy of the source is purely
	// for testing
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	copy(key.PrivateKey[:], privateKey[:])
	copy(key.PublicKey[:], publicKey[:])

	// prepend 0x01 for ed25519 public key
	typedPublicKeyBytes := append([]byte{0x01}, key.PublicKey...)
	hasher := ripemd160.New()
	hasher.Write(typedPublicKeyBytes)
	key.Address, err = acm.AddressFromBytes(hasher.Sum(nil))
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (mockKey *MockKey) Sign(message []byte) (crypto.Signature, error) {
	sigEd25519 := crypto.SignatureEd25519{}
	signatureBytes := ed25519.Sign(&mockKey.PrivateKey, message)
	copy(sigEd25519[:], signatureBytes[:])
	return sigEd25519.Wrap(), nil
}

//---------------------------------------------------------------------
// Mock client for replacing signing done by monax-keys

// Implementation assertion
var _ KeyClient = (*MockKeyClient)(nil)

type MockKeyClient struct {
	knownKeys map[acm.Address]*MockKey
}

func NewMockKeyClient() *MockKeyClient {
	return &MockKeyClient{
		knownKeys: make(map[acm.Address]*MockKey),
	}
}

func (mock *MockKeyClient) NewKey() acm.Address {
	// Only tests ED25519 curve and ripemd160.
	key, err := newMockKey()
	if err != nil {
		panic(fmt.Sprintf("Mocked key client failed on key generation: %s", err))
	}
	mock.knownKeys[key.Address] = key
	return key.Address
}

func (mock *MockKeyClient) Sign(signAddress acm.Address, message []byte) (crypto.Signature, error) {
	key := mock.knownKeys[signAddress]
	if key == nil {
		return crypto.Signature{}, fmt.Errorf("Unknown address (%s)", signAddress)
	}
	return key.Sign(message)
}

func (mock *MockKeyClient) PublicKey(address acm.Address) (acm.PublicKey, error) {
	key := mock.knownKeys[address]
	if key == nil {
		return acm.PublicKey{}, fmt.Errorf("Unknown address (%s)", address)
	}
	pubKeyEd25519 := acm.PublicKeyEd25519{}
	copy(pubKeyEd25519[:], key.PublicKey)
	return pubKeyEd25519.Wrap(), nil
}

func (mock *MockKeyClient) Generate(keyName string, keyType KeyType) (acm.Address, error) {
	return mock.NewKey(), nil
}
