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
	trim       bool
	powers     map[crypto.Address]*big.Int
	publicKeys map[crypto.Address]crypto.PublicKey
	totalPower *big.Int
}

type ValidatorSet interface {
	AlterPower(publicKey crypto.PublicKey, power *big.Int) (flow *big.Int, err error)
}

func NewValidators() *Validators {
	return newValidators(true)
}

// Returns validators that continue to store zeroed (removed) validators so total flow can be calculated
func NewValidatorsNoTrim() *Validators {
	return newValidators(false)
}

// Create a new Validators which can act as an accumulator for validator power changes
func newValidators(trim bool) *Validators {
	return &Validators{
		trim:       trim,
		totalPower: new(big.Int),
		powers:     make(map[crypto.Address]*big.Int),
		publicKeys: make(map[crypto.Address]crypto.PublicKey),
	}
}

// Add the power of a validator and returns the flow into that validator
func (vs *Validators) AlterPower(publicKey crypto.PublicKey, power *big.Int) (*big.Int, error) {
	address := publicKey.Address()
	// Calculcate flow into this validator (postive means in, negative means out)
	flow := new(big.Int).Sub(power, vs.Power(address))
	if vs.trim && power.Cmp(big0) == 0 {
		// Remove from set so that we return an accurate length
		delete(vs.publicKeys, address)
		delete(vs.powers, address)
		return flow, nil
	}
	vs.totalPower.Add(vs.totalPower, flow)
	vs.publicKeys[address] = publicKey
	vs.powers[address] = power
	return flow, nil
}

func (vs *Validators) AddPower(publicKey crypto.PublicKey, power *big.Int) error {
	// Current power + power
	_, err := vs.AlterPower(publicKey, new(big.Int).Add(vs.Power(publicKey.Address()), power))
	return err
}

func (vs *Validators) SubtractPower(publicKey crypto.PublicKey, power *big.Int) error {
	// Current power - power
	_, err := vs.AlterPower(publicKey, new(big.Int).Sub(vs.Power(publicKey.Address()), power))
	return err
}

func (vs *Validators) Power(address crypto.Address) *big.Int {
	if vs.powers[address] == nil {
		zero := new(big.Int)
		vs.powers[address] = zero
		return zero
	}
	return vs.powers[address]
}

// Iterates over validators sorted by address
func (vs *Validators) Iterate(iter func(publicKey crypto.PublicKey, power *big.Int) (stop bool)) (stopped bool) {
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

func (vs *Validators) Size() int {
	return len(vs.publicKeys)
}

func (vs *Validators) TotalPower() *big.Int {
	return vs.totalPower
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
	vs.Iterate(func(publicKey crypto.PublicKey, power *big.Int) (stop bool) {
		_, err = cdc.MarshalBinaryWriter(buffer, PersistedValidator{
			PublicKey:  publicKey,
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
	return fmt.Sprintf("Validators{TotalPower: %v; Size: %v}", vs.TotalPower(), vs.Size())
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
			_, err = validators.AlterPower(pv.PublicKey, power)
			if err != nil {
				return err
			}
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}
