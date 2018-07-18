package blockchain

import (
	"fmt"
	"math/big"

	"github.com/hyperledger/burrow/crypto"
)

type ValidatorsWindow struct {
	Buckets []*Validators
	Total   *Validators
	head    int
}

// Provides a sliding window over the last size buckets of validator power changes
func NewValidatorsWindow(size int) *ValidatorsWindow {
	if size < 1 {
		size = 1
	}
	vw := &ValidatorsWindow{
		Buckets: make([]*Validators, size),
		Total:   NewValidators(),
	}
	vw.Buckets[vw.head] = NewValidatorsNoTrim()
	return vw
}

// Updates the current head bucket (accumulator)
func (vw *ValidatorsWindow) AlterPower(publicKey crypto.PublicKey, power *big.Int) (*big.Int, error) {
	return vw.Buckets[vw.head].AlterPower(publicKey, power)
}

// Returns the change in power over the window during this commit and the total flow which is the sum of absolute values
// of all changes in particular validator's power
func (vw *ValidatorsWindow) FlushInto(validatorsToUpdate ValidatorSet) (totalPowerChange *big.Int, totalFlow *big.Int, err error) {
	totalPowerChange = new(big.Int).Add(big0, vw.Total.TotalPower())
	// The absolute flow from validators
	totalFlow = new(big.Int)
	var flow *big.Int
	if vw.Buckets[vw.head].Iterate(func(publicKey crypto.PublicKey, power *big.Int) (stop bool) {
		// Update the sink validators
		flow, err = validatorsToUpdate.AlterPower(publicKey, power)
		if err != nil {
			return true
		}
		totalFlow.Add(totalFlow, flow.Abs(flow))
		// Add to total power
		err = vw.Total.AddPower(publicKey, power)
		if err != nil {
			return true
		}
		return
	}) {
		// If iteration stopped there was an error
		return nil, nil, err
	}
	// move the ring buffer on
	vw.head = (vw.head + 1) % len(vw.Buckets)
	// Subtract the tail bucket (if any) from the total
	if vw.Buckets[vw.head].Iterate(func(publicKey crypto.PublicKey, power *big.Int) (stop bool) {
		err = vw.Total.SubtractPower(publicKey, power)
		if err != nil {
			return true
		}
		return
	}) {
		return nil, nil, err
	}
	// Clear new head bucket (and possibly previous tail)
	vw.Buckets[vw.head] = NewValidatorsNoTrim()
	return totalPowerChange.Sub(vw.Total.TotalPower(), totalPowerChange), totalFlow, nil
}

func (vw *ValidatorsWindow) Size() int {
	return len(vw.Buckets)
}

// Returns buckets in order
func (vw *ValidatorsWindow) OrderedBuckets() []*Validators {
	length := len(vw.Buckets)
	buckets := make([]*Validators, length)
	for i := 0; i < length; i++ {
		buckets[i] = vw.Buckets[(length+vw.head-i)%length]
	}
	return buckets
}

func (vw *ValidatorsWindow) String() string {
	return fmt.Sprintf("ValidatorsWindow{Total: %v; Buckets: Head->%v<-Tail}", vw.Total, vw.OrderedBuckets())
}
