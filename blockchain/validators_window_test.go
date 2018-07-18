package blockchain

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var pubA = pubKey(1)
var pubB = pubKey(2)
var pubC = pubKey(3)

func TestValidatorsWindow_AlterPower(t *testing.T) {
	vs := NewValidators()
	vw := NewValidatorsWindow(5)

	var powA, powB, powC int64 = 100, 23, 1113
	powerChange, totalFlow := alterPowers(t, vw, vs, powA, powB, powC)
	expectedChange := big.NewInt(powA + powB + powC)
	assert.Equal(t, expectedChange, powerChange)
	assert.Equal(t, expectedChange, totalFlow)

	powerChange, totalFlow = alterPowers(t, vw, vs, powA, powB, powC)
	assert.Equal(t, expectedChange, powerChange)
	assert.Equal(t, big0, totalFlow)

	powerChange, totalFlow = alterPowers(t, vw, vs, powA, powB, powC)
	assert.Equal(t, expectedChange, powerChange)
	assert.Equal(t, big0, totalFlow)

	powerChange, totalFlow = alterPowers(t, vw, vs, powA, powB, powC)
	assert.Equal(t, expectedChange, powerChange)
	assert.Equal(t, big0, totalFlow)

	// Now we have filled the window there will be no change
	powerChange, totalFlow = alterPowers(t, vw, vs, powA, powB, powC)
	assert.True(t, big0.Cmp(powerChange) == 0)
	assert.Equal(t, big0, totalFlow)

	powA, powB, powC = 50, 43, 1103
	powerChange, totalFlow = alterPowers(t, vw, vs, powA, powB, powC)
	assert.Equal(t, big.NewInt(-40), powerChange)
	assert.Equal(t, big.NewInt(80), totalFlow)
}

func alterPowers(t testing.TB, vw *ValidatorsWindow, vs ValidatorSet, powA, powB, powC int64) (powerChange, totalFlow *big.Int) {
	t.Log(vw)
	var err error
	_, err = vw.AlterPower(pubA, big.NewInt(powA))
	require.NoError(t, err)
	_, err = vw.AlterPower(pubB, big.NewInt(powB))
	require.NoError(t, err)
	_, err = vw.AlterPower(pubC, big.NewInt(powC))
	require.NoError(t, err)
	powerChange, totalFlow, err = vw.FlushInto(vs)
	require.NoError(t, err)
	return powerChange, totalFlow
}
