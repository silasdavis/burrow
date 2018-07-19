package blockchain

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"sort"

	"github.com/hyperledger/burrow/crypto"
	"github.com/tendermint/go-amino"
)

var big0 = big.NewInt(0)

const maxAminoReadSizeBytes int64 = 1 << 8

// A Validator multiset - can be used to capture the global state of validators or as an accumulator each block
type Validators struct {
	powers     map[crypto.Address]*big.Int
	publicKeys map[crypto.Address]crypto.Addressable
	totalPower *big.Int
}

type ValidatorSet interface {
	AlterPower(id crypto.Addressable, power *big.Int) (flow *big.Int, err error)
}

// Create a new Validators which can act as an accumulator for validator power changes
func NewValidators() *Validators {
	return &Validators{

		totalPower: new(big.Int),
		powers:     make(map[crypto.Address]*big.Int),
		publicKeys: make(map[crypto.Address]crypto.Addressable),
	}
}

// Add the power of a validator and returns the flow into that validator
func (vs *Validators) AlterPower(id crypto.Addressable, power *big.Int) *big.Int {
	address := id.Address()
	// Calculcate flow into this validator (postive means in, negative means out)
	flow := new(big.Int).Sub(power, vs.Power(address))
	if power.Cmp(big0) == 0 {
		// Remove from set so that we return an accurate length
		delete(vs.publicKeys, address)
		delete(vs.powers, address)
		return flow
	}
	vs.totalPower.Add(vs.totalPower, flow)
	vs.publicKeys[address] = crypto.MemoizeAddressable(id)
	vs.powers[address] = power
	return flow
}

func (vs *Validators) AddPower(id crypto.Addressable, power *big.Int) {
	// Current power + power
	vs.AlterPower(id, new(big.Int).Add(vs.Power(id.Address()), power))
}

func (vs *Validators) SubtractPower(id crypto.Addressable, power *big.Int) {
	// Current power - power
	vs.AlterPower(id, new(big.Int).Sub(vs.Power(id.Address()), power))
}

func (vs *Validators) Power(address crypto.Address) *big.Int {
	if vs.powers[address] == nil {
		return new(big.Int)
	}
	return vs.powers[address]
}

// Iterates over validators sorted by address
func (vs *Validators) Iterate(iter func(id crypto.Addressable, power *big.Int) (stop bool)) (stopped bool) {
	if vs == nil {
		return
	}
	addresses := make(crypto.Addresses, 0, len(vs.powers))
	for address := range vs.powers {
		addresses = append(addresses, address)
	}
	sort.Sort(addresses)
	for _, address := range addresses {
		if iter(vs.publicKeys[address], vs.Power(address)) {
			return true
		}
	}
	return
}

func (vs *Validators) Count() int {
	return len(vs.publicKeys)
}

func (vs *Validators) TotalPower() *big.Int {
	return new(big.Int).Set(vs.totalPower)
}

func (vs *Validators) Copy() *Validators {
	vsCopy := NewValidators()
	vs.Iterate(func(id crypto.Addressable, power *big.Int) (stop bool) {
		vsCopy.AlterPower(id, power)
		return
	})
	return vsCopy
}

type PersistedValidator struct {
	PublicKey  crypto.PublicKey
	PowerBytes []byte
}

var cdc = amino.NewCodec()

// Uses the fixed width public key encoding to
func (vs *Validators) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	// varint buffer
	var err error
	vs.Iterate(func(id crypto.Addressable, power *big.Int) (stop bool) {
		_, err = cdc.MarshalBinaryWriter(buffer, PersistedValidator{
			PublicKey:  id.PublicKey(),
			PowerBytes: power.Bytes(),
		})
		if err != nil {
			return true
		}
		return
	})
	return buffer.Bytes(), err
}

func (vs *Validators) String() string {
	return fmt.Sprintf("Validators{TotalPower: %v; Count: %v}", vs.TotalPower(), vs.Count())
}

// Decodes validators encoded with Encode - expects the exact encoded size with no trailing bytes
func DecodeValidators(encoded []byte, validators *Validators) error {
	buffer := bytes.NewBuffer(encoded)
	var err error
	for err == nil {
		pv := new(PersistedValidator)
		_, err = cdc.UnmarshalBinaryReader(buffer, pv, maxAminoReadSizeBytes)
		if err == nil {
			power := new(big.Int).SetBytes(pv.PowerBytes)
			validators.AlterPower(pv.PublicKey, power)
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}
