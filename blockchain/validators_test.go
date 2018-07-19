package blockchain

import (
	"math/big"
	"testing"

	"fmt"

	"math/rand"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidators_AlterPower(t *testing.T) {
	vs := NewValidators()
	pow1 := big.NewInt(2312312321)
	vs.AlterPower(pubKey(1), pow1)
	assert.Equal(t, pow1, vs.TotalPower())
}

func TestValidators_Encode(t *testing.T) {
	vs := NewValidators()
	rnd := rand.New(rand.NewSource(43534543))
	for i := 0; i < 100; i++ {
		power := big.NewInt(rnd.Int63n(10))
		vs.AlterPower(pubKey(rnd.Int63()), power)
	}
	encoded, err := vs.Encode()
	require.NoError(t, err)
	vsOut := NewValidators()
	require.NoError(t, DecodeValidators(encoded, vsOut))
	// Check decoded matches encoded
	var publicKeyPower []interface{}
	vs.Iterate(func(id crypto.Addressable, power *big.Int) (stop bool) {
		publicKeyPower = append(publicKeyPower, id, power)
		return
	})
	vsOut.Iterate(func(id crypto.Addressable, power *big.Int) (stop bool) {
		assert.Equal(t, publicKeyPower[0], id)
		assert.Equal(t, publicKeyPower[1], power)
		publicKeyPower = publicKeyPower[2:]
		return
	})
	assert.Len(t, publicKeyPower, 0, "should exhaust all validators in decoded multiset")
}

func pubKey(secret interface{}) crypto.PublicKey {
	return acm.NewConcreteAccountFromSecret(fmt.Sprintf("%v", secret)).PublicKey
}
