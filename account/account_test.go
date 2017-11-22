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

package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-wire"
)

func TestAddress(t *testing.T) {
	bs := []byte{
		1, 2, 3, 4, 5,
		1, 2, 3, 4, 5,
		1, 2, 3, 4, 5,
		1, 2, 3, 4, 5,
	}
	addr, err := AddressFromBytes(bs)
	assert.NoError(t, err)
	word256 := addr.Word256()
	leadingZeroes := []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
	assert.Equal(t, leadingZeroes, word256[:12])
	addrFromWord256 := AddressFromWord256(word256)
	assert.Equal(t, bs, addrFromWord256[:])
	assert.Equal(t, addr, addrFromWord256)
}

func TestAccountSerialise(t *testing.T) {

	type AccountContainingStruct struct {
		Account Account
		ChainID string
	}
	// This test is really testing this go wire declaration in account.go
	//var _ = wire.RegisterInterface(struct{ Account }{}, wire.ConcreteType{concreteAccountWrapper{}, 0x01})

	acc := AsConcreteAccount(FromAddressable(GeneratePrivateAccountFromSecret("Super Secret Secret")))

	// Go wire cannot serialise a top-level interface type it needs to be a field or sub-field of a struct
	// at some depth. i.e. you MUST wrap an interface if you want it to be decoded (they do not document this)
	var accStruct = AccountContainingStruct{
		Account: acc.Account(),
		ChainID: "TestChain",
	}

	// We will write into this
	accStructOut := AccountContainingStruct{}

	// We must pass in a value type to read from (accStruct), but provide a pointer type to write into (accStructout
	wire.ReadBinaryBytes(wire.BinaryBytes(accStruct), &accStructOut)

	assert.Equal(t, accStruct, accStructOut)
}

func TestDecodeConcrete(t *testing.T) {
	concreteAcc := AsConcreteAccount(FromAddressable(GeneratePrivateAccountFromSecret("Super Semi Secret")))
	acc := concreteAcc.Account()
	concreteAccOut, err := DecodeConcrete(acc.Encode())
	assert.NoError(t, err)
	assert.Equal(t, concreteAcc, *concreteAccOut)

	concreteAccOut, err = DecodeConcrete([]byte("flungepliffery munknut tolopops"))
	assert.Error(t, err)
}

func TestDecode(t *testing.T) {
	concreteAcc := AsConcreteAccount(FromAddressable(GeneratePrivateAccountFromSecret("Super Semi Secret")))
	acc := concreteAcc.Account()
	accOut := Decode(acc.Encode())
	assert.Equal(t, concreteAcc, AsConcreteAccount(accOut))

	accOut = Decode([]byte("flungepliffery munknut tolopops"))
	assert.Nil(t, accOut)
}
