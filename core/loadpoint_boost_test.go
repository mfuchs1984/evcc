package core

import (
	"testing"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/core/site"
	"github.com/evcc-io/evcc/util"
	"github.com/stretchr/testify/assert"
)

type mockSite struct {
	site.API
	maxDischargePower float64
	residualPower     float64
}

func (m *mockSite) GetBatteryMaxDischargePower() float64 {
	return m.maxDischargePower
}

func (m *mockSite) GetResidualPower() float64 {
	return m.residualPower
}

func TestBoostPower(t *testing.T) {
	Voltage = 230
	lp := &Loadpoint{
		log:          util.NewLogger("lp"),
		status:       api.StatusC,
		batteryBoost: boostStart,
		maxCurrent:   16,
		phases:       3,
	}
	s := &mockSite{}
	lp.site = s

	// No max discharge power limit
	s.maxDischargePower = 0
	// EffectiveMaxPower will be 230 * 16 * 3 = 11040
	delta := lp.boostPower(0)
	assert.Equal(t, 11040.0, delta)
	assert.Equal(t, boostContinue, lp.batteryBoost)

	// With max discharge power limit
	s.maxDischargePower = 5000
	lp.batteryBoost = boostStart
	delta = lp.boostPower(0)
	assert.Equal(t, 5000.0, delta)
	assert.Equal(t, boostContinue, lp.batteryBoost)

	// boostContinue with limit
	lp.batteryBoost = boostContinue
	s.residualPower = 0
	// delta = math.Max(100, 0) = 100
	// plus EffectiveStepPower = 690
	// delta = 790
	// delta = min(790, max(0, 5000 - 0)) = 790
	// res = 0 + 790 + 0 = 790
	delta = lp.boostPower(0)
	assert.Equal(t, 790.0, delta)

	// boostContinue at limit
	// delta = min(790, max(0, 5000 - 5000)) = 0
	// res = 5000 + 0 + 0 = 5000
	delta = lp.boostPower(5000)
	assert.Equal(t, 5000.0, delta)

	// boostContinue over limit
	// delta = min(790, max(0, 5000 - 6000)) = 0
	// res = 6000 + 0 + 0 = 6000
	delta = lp.boostPower(6000)
	assert.Equal(t, 6000.0, delta)

	// boostStart while battery is charging (negative power)
	// battery charging at 2000W, limit is 5000W
	// max discharge capacity = 5000 - (-2000) = 7000W
	// res = max(0, -2000) + 7000 + 0 = 7000W
	lp.batteryBoost = boostStart
	delta = lp.boostPower(-2000)
	assert.Equal(t, 7000.0, delta)

	// boostContinue while battery is charging (negative power)
	// limit is 50W (less than the standard 790W delta)
	// without raw negative power, delta would be restricted to 50W
	// with raw negative power (-2000W), headroom is 2050W, so delta is allowed to be 790W
	s.maxDischargePower = 50
	s.residualPower = 0 // base delta = 100 + 690 = 790
	lp.batteryBoost = boostContinue
	delta = lp.boostPower(-2000)
	// res = max(0, -2000) + 790 + 0 = 790W
	assert.Equal(t, 790.0, delta)
}

func TestPvMaxCurrentForcedBatteryCharge(t *testing.T) {
	Voltage = 230
	lp := &Loadpoint{
		log:          util.NewLogger("lp"),
		status:       api.StatusC,
		batteryBoost: boostContinue,
		maxCurrent:   16,
		minCurrent:   6,
		phases:       3,
		enabled:      true,
		chargePower:  4000.0,
	}
	s := &mockSite{
		maxDischargePower: 5000,
		residualPower:     0,
	}
	lp.site = s

	// Scenario: Battery is forced charging at 3000W, car is drawing 4000W.
	// Since the battery refuses to yield its forced charge, the grid supplies 7000W total.
	// sitePower = gridPower (7000) + batteryPower (-3000) = 4000
	sitePower := 4000.0
	batteryBoostPower := -3000.0

	// In pvMaxCurrent:
	// boostPower(-3000) will authorize 3000W.
	// sitePower becomes 4000 - 3000 = 1000W.
	// deltaCurrent = powerToCurrent(-1000) = -1.44A.
	// targetCurrent drops by 1.44A. Since we didn't even set offeredCurrent,
	// the base is 0, making targetCurrent 0. Either way, it's < minCurrent (6A).
	// pvMaxCurrent will return minCurrent to throttle down the over-budget charging.
	current := lp.pvMaxCurrent(api.ModePV, sitePower, batteryBoostPower, false, false)

	assert.Equal(t, 6.0, current)
}

func TestPvMaxCurrentSpeedUpStep(t *testing.T) {
	Voltage = 230
	lp := &Loadpoint{
		log:          util.NewLogger("lp"),
		status:       api.StatusC,
		batteryBoost: boostContinue,
		maxCurrent:   32, // Allow high current for the test
		minCurrent:   6,
		phases:       3,
		enabled:      true,
	}
	s := &mockSite{
		maxDischargePower: 5000,
		residualPower:     0,
	}
	lp.site = s

	// Cycle 1: Car is at 0W. PV is 3000W. Battery is charging at -3000W.
	// sitePower = grid + battery = 0 + (-3000) = -3000W
	lp.chargePower = 0.0
	lp.offeredCurrent = 0.0
	sitePower1 := -3000.0
	batteryBoostPower1 := -3000.0

	targetCurrent1 := lp.pvMaxCurrent(api.ModePV, sitePower1, batteryBoostPower1, false, false)
	
	// Because batteryBoostPower is -3000, delta = 3000.
	// Target current will be PV + 3000 = 6000W (8.69A).
	assert.InDelta(t, 8.69, targetCurrent1, 0.1)

	// Cycle 2: Car followed the target and is now drawing 6000W!
	// Since PV is 3000W and Car is 6000W, the battery MUST supply the remaining 3000W.
	// So the battery flips into discharging at 3000W.
	// sitePower = grid (0) + battery (3000) = 3000W
	lp.chargePower = 6000.0
	lp.offeredCurrent = targetCurrent1
	sitePower2 := 3000.0
	batteryBoostPower2 := 3000.0 // positive now, battery is discharging

	targetCurrent2 := lp.pvMaxCurrent(api.ModePV, sitePower2, batteryBoostPower2, false, false)

	// In Cycle 2, because battery is discharging, delta reverts to the safe 790W.
	// targetCurrent will be 6000W + 790W = 6790W (9.84A).
	// This proves that targetCurrent increases and the hand-off is seamless!
	assert.Greater(t, targetCurrent2, targetCurrent1, "Target current gracefully increases when car follows the target")
	assert.InDelta(t, 9.84, targetCurrent2, 0.1)
}
