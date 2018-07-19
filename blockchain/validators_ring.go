package blockchain

import (
	"fmt"
	"math/big"

	"github.com/hyperledger/burrow/crypto"
)

type ValidatorsRing struct {
	buckets []*Validators
	total   *Validators
	flow    *big.Int
	head    int
}

var big1 = big.NewInt(1)
var big3 = big.NewInt(3)

// Provides a sliding window over the last size buckets of validator power changes
func NewValidatorsRing(initialValidators *Validators, size int) *ValidatorsRing {
	if size < 2 {
		size = 2
	}
	vw := &ValidatorsRing{
		buckets: make([]*Validators, size),
		total:   NewValidators(),
		flow:    big.NewInt(0),
	}
	// Existing set
	vw.buckets[vw.index(-1)] = initialValidators.Copy()
	// Current accumulator
	vw.buckets[vw.head] = initialValidators.Copy()
	return vw
}

// Updates the current head bucket (accumulator) with some safety checks
func (vw *ValidatorsRing) AlterPower(id crypto.Addressable, power *big.Int) (*big.Int, error) {
	if power.Sign() == -1 {
		return nil, fmt.Errorf("cannot set negative validator power: %v", power)
	}
	// if flow > maxflow then we cannot alter the power
	flow := vw.Flow(id, power)
	maxFlow := vw.MaxFlow()
	if new(big.Int).Add(flow, vw.flow).Cmp(maxFlow) == 1 {
		allowable := new(big.Int).Sub(maxFlow, vw.flow)
		return nil, fmt.Errorf("cannot change validator power of %v to %v because that would result in a flow " +
			"greater than or equal to 1/3 of total power for the next commit: flow induced by change: %v, " +
			"current total flow: %v/%v (cumulative/max), remaining allowable flow: %v",
			id.Address(), power, flow, vw.flow, maxFlow, allowable)
	}
	vw.flow.Add(vw.flow, flow)
	return vw.Head().AlterPower(id, power), nil
}

// Returns the flow that would be induced by a validator change by comparing the current accumulator with the previous
// bucket
func (vw *ValidatorsRing) Flow(id crypto.Addressable, power *big.Int) *big.Int {
	flow := new(big.Int)
	return flow.Abs(flow.Sub(power, vw.Prev().Power(id.Address())))
}

// To ensure that in the maximum valildator shift at least one unit
// of validator power in the intersection of last block validators and this block validators must have at least one
// non-byzantine validator who can tell you if you've been lied to about the validator set
// So need at most ceiling((Total Power)/3) - 1, in integer division we have ceiling(X*p/q) = (p(X+1)-1)/q
// For p = 1 just X/q
// So we want (Total Power)/3 - 1
func (vw *ValidatorsRing) MaxFlow() *big.Int {
	max := vw.Prev().TotalPower()
	return max.Sub(max.Div(max, big3), big1)
}

// Advance the current head bucket to the next bucket and returns the change in total power between the previous bucket
// and the current head, and the total flow which is the sum of absolute values of all changes each validator's power
// after rotation the next head is a copy of the current head
func (vw *ValidatorsRing) Rotate() (totalPowerChange *big.Int, totalFlow *big.Int) {
	// The absolute flow from validators
	// initialise the nextValidators bucket to be a copy of the previous bucket so we can can calculate the flow between
	// the previous bucket and the current bucket when we flush into nextValidators
	prevAndNext := vw.Prev().Copy()
	totalFlow = new(big.Int)
	// Subtract the previous bucket total power so we can add on the current buckets power after this
	totalPowerChange = new(big.Int).Sub(vw.Head().totalPower, prevAndNext.totalPower)
	vw.buckets[vw.head].Iterate(func(id crypto.Addressable, power *big.Int) (stop bool) {
		// Update the sink validators
		flow := prevAndNext.AlterPower(id, power)
		totalFlow.Add(totalFlow, flow.Abs(flow))
		// Add to total power
		vw.total.AddPower(id, power)
		return
	})
	// move the ring buffer on
	vw.head = vw.index(1)
	// Subtract the tail bucket (if any) from the total
	vw.buckets[vw.head].Iterate(func(id crypto.Addressable, power *big.Int) (stop bool) {
		vw.total.SubtractPower(id, power)
		return
	})
	// Overwrite new head bucket (previous tail) with previous bucket copy updated with current head
	vw.buckets[vw.head] = prevAndNext
	vw.flow.SetInt64(0)
	return totalPowerChange, totalFlow
}

func (vw *ValidatorsRing) Prev() *Validators {
	return vw.buckets[vw.index(-1)]
}

func (vw *ValidatorsRing) Head() *Validators {
	return vw.buckets[vw.head]
}

func (vw *ValidatorsRing) Next() *Validators {
	return vw.buckets[vw.index(1)]
}

func (vw *ValidatorsRing) index(i int) int {
	return (len(vw.buckets) + vw.head + i) % len(vw.buckets)
}

func (vw *ValidatorsRing) Size() int {
	return len(vw.buckets)
}

// Returns buckets in order head, previous, ...
func (vw *ValidatorsRing) OrderedBuckets() []*Validators {
	buckets := make([]*Validators, len(vw.buckets))
	for i := 0; i < len(buckets); i++ {
		buckets[i] = vw.buckets[vw.index(-i)]
	}
	return buckets
}

func (vw *ValidatorsRing) String() string {
	return fmt.Sprintf("ValidatorsWindow{Total: %v; Buckets: Head->%v<-Tail}", vw.total, vw.OrderedBuckets())
}
