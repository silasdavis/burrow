package blockchain

import (
	"fmt"
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
	powAInitial := int64(10000)
	vs.AlterPower(pubA, big.NewInt(powAInitial))
	vw := NewValidatorsRing(vs, 5)

	// Just allowable validator tide
	var powA, powB, powC int64 = 7000, 23, 309
	powerChange, totalFlow, err := alterPowers(t, vw, powA, powB, powC)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(powA+powB+powC-powAInitial), powerChange)
	assert.Equal(t, big.NewInt(powAInitial/3 - 1), totalFlow)

	// This one is not
	vw = NewValidatorsRing(vs, 5)
	powA, powB, powC = 7000, 23, 308
	powerChange, totalFlow, err = alterPowers(t, vw, powA, powB, powC)
	require.Error(t, err)

	//powerChange, totalFlow, err = alterPowers(t, vw, powA, powB, powC)
	//require.NoError(t, err)
	//assertZero(t, powerChange)
	//assert.Equal(t, big0, totalFlow)
	//
	//powerChange, totalFlow, err = alterPowers(t, vw, powA, powB, powC)
	//require.NoError(t, err)
	//assertZero(t, powerChange)
	//assert.Equal(t, big0, totalFlow)
	//
	//powerChange, totalFlow, err = alterPowers(t, vw, powA, powB, powC)
	//require.NoError(t, err)
	//assertZero(t, powerChange)
	//assert.Equal(t, big0, totalFlow)
	//
	//// Now we have filled the window there will be no change
	//powerChange, totalFlow, err = alterPowers(t, vw, powA, powB, powC)
	//require.NoError(t, err)
	//assertZero(t, powerChange)
	//assert.Equal(t, big0, totalFlow)
	//
	//powA, powB, powC = 50, 43, 1103
	//powerChange, totalFlow, err = alterPowers(t, vw, powA, powB, powC)
	//require.NoError(t, err)
	//assert.Equal(t, big.NewInt(-40), powerChange)
	//assert.Equal(t, big.NewInt(80), totalFlow)
	//
	//powA, powB, powC = 0, 43, 1103
	//powerChange, totalFlow, err = alterPowers(t, vw, powA, powB, powC)
	//require.NoError(t, err)
	//assert.Equal(t, big.NewInt(-50), powerChange)
	//assert.Equal(t, big.NewInt(50), totalFlow)
}

func alterPowers(t testing.TB, vw *ValidatorsRing, powA, powB, powC int64) (powerChange, totalFlow *big.Int, err error) {
	t.Log(vw)
	_, err = vw.AlterPower(pubA, big.NewInt(powA))
	if err != nil {
		return nil, nil, err
	}
	_, err = vw.AlterPower(pubB, big.NewInt(powB))
	if err != nil {
		return nil, nil, err
	}
	_, err = vw.AlterPower(pubC, big.NewInt(powC))
	if err != nil {
		return nil, nil, err
	}
	maxFlow := vw.MaxFlow()
	powerChange, totalFlow = vw.Rotate()
	if maxFlow.Cmp(totalFlow) == 1 {
		return powerChange, totalFlow, fmt.Errorf("totalFlow (%v) exceeded maxFlow (%v)", totalFlow, maxFlow)
	}

	return powerChange, totalFlow, nil
}

// Since we have -0 and 0 with big.Int due to its representation with a neg flag
func assertZero(t testing.TB, i *big.Int) {
	assert.True(t, big0.Cmp(i) == 0, "expected 0 but got %v", i)
}
