package core

import (
	"testing"

	evbus "github.com/asaskevich/EventBus"
	"github.com/benbjohnson/clock"
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/core/types"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGreenShare(t *testing.T) {
	tc := []struct {
		title                                                 string
		grid, pv, battery, home, lp                           float64
		greenShareTotal, greenShareHome, greenShareLoadpoints float64
	}{
		{
			"half grid, half pv, green home",
			1000, 1000, 0, 1000, 1000,
			0.5, 1, 0,
		},
		{
			"half grid, half pv, no home",
			1000, 1000, 0, 0, 2000,
			0.5, 1, 0.5,
		},
		{
			"half grid, half pv, no lp",
			2500, 2500, 0, 5000, 0,
			0.5, 0.5, 0,
		},
		{
			"full pv",
			0, 5000, 0, 1000, 4000,
			1, 1, 1,
		},
		{
			"full grid",
			5000, 0, 0, 1000, 4000,
			0, 0, 0,
		},
		{
			"half grid, half battery, green home",
			1000, 0, 1000, 1000, 1000,
			0.5, 1, 0,
		},
		{
			"half grid, half battery, no home",
			1000, 0, 1000, 0, 2000,
			0.5, 1, 0.5,
		},
		{
			"half grid, half battery, no lp",
			1000, 0, 1000, 2000, 0,
			0.5, 0.5, 0,
		},
		{
			"full pv, pv export",
			-5000, 10000, 0, 1000, 4000,
			1, 1, 1,
		},
		{
			"full pv, pv export, no lp",
			-5000, 10000, 0, 5000, 0,
			1, 1, 1,
		},
		{
			"full pv, pv export, battery charge",
			-2500, 10000, -2500, 1000, 4000,
			1, 1, 1,
		},
		{
			"full grid, battery charge",
			3000, 0, -1000, 1000, 1000,
			0, 0, 0,
		},
		{
			"full grid, battery charge, no lp",
			2000, 0, -1000, 1000, 0,
			0, 0, 0,
		},
		{
			"half grid, half pv, battery charge, no lp",
			1000, 1000, -1000, 1000, 0,
			0.5, 1, 0,
		},
		{
			"half grid, half pv, battery charge, home, lp",
			1000, 1000, -1000, 500, 500,
			0.5, 1, 0,
		},
		{
			"pv ac limited, battery charge & grid import",
			1000, 3000, -1000, 1000, 2000,
			0.75, 1, 0.5,
		},
	}

	for _, tc := range tc {
		t.Log(tc.title)

		s := &Site{
			gridPower: tc.grid,
			pvPower:   tc.pv,
			battery: types.BatteryState{
				Power: tc.battery,
			},
		}

		totalPower := tc.grid + tc.pv + max(0, tc.battery)
		greenShareTotal := s.greenShare(0, totalPower)
		if greenShareTotal != tc.greenShareTotal {
			t.Errorf("greenShareTotal wanted %.3f, got %.3f", tc.greenShareTotal, greenShareTotal)
		}
		greenShareHome := s.greenShare(0, tc.home)
		if greenShareHome != tc.greenShareHome {
			t.Errorf("greenShareHome wanted %.3f, got %.3f", tc.greenShareHome, greenShareHome)
		}
		greenShareLoadpoints := s.greenShare(tc.home+max(0, -tc.battery), totalPower)
		if greenShareLoadpoints != tc.greenShareLoadpoints {
			t.Errorf("greenShareLoadpoints wanted %.3f, got %.3f", tc.greenShareLoadpoints, greenShareLoadpoints)
		}
	}
}

// chargerExWrapper wraps MockCharger and adds MaxCurrentMillis for sub-amp control
type chargerExWrapper struct {
	*api.MockCharger
	current float64
}

func (c *chargerExWrapper) MaxCurrentMillis(current float64) error {
	c.current = current
	return nil
}

// boostTestCurrents exercises the real lp.Update → pvMaxCurrent → boostPower code path
// for above and below prioritySoc. It returns the target currents set on the charger for
// each case so the caller can compare them. When useChargerEx is true, a ChargerEx wrapper
// is used for sub-amp precision; otherwise plain MockCharger (integer amps) is used.
func boostTestCurrents(t *testing.T, batteryPow float64, initialCurrent float64, useChargerEx bool) (aboveCurrent, belowCurrent float64) {
	t.Helper()

	tcs := []struct {
		title       string
		batterySoc  float64
		prioritySoc float64
	}{
		{"above prioritySoc", 80, 50},
		{"below prioritySoc", 30, 50},
	}

	var currents []float64

	for _, tc := range tcs {
		t.Run(tc.title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockCharger := api.NewMockCharger(ctrl)

			Voltage = 230

			site := &Site{
				log:           util.NewLogger("site"),
				batteryMeters: []config.Device[api.Meter]{nil}, // non-empty to trigger priority logic
				battery: types.BatteryState{
					Power: batteryPow,
					Soc:   tc.batterySoc,
				},
				prioritySoc:   tc.prioritySoc,
				ResidualPower: 100,
			}

			// choose charger: ChargerEx for sub-amp precision, plain Charger for integer amps
			var charger api.Charger
			var wrapper *chargerExWrapper
			if useChargerEx {
				wrapper = &chargerExWrapper{MockCharger: mockCharger}
				charger = wrapper
			} else {
				charger = mockCharger
			}

			lp := &Loadpoint{
				log:               util.NewLogger("lp"),
				bus:               evbus.New(),
				clock:             clock.NewMock(),
				charger:           charger,
				chargeMeter:       &Null{},
				chargeRater:       &Null{},
				chargeTimer:       &Null{},
				wakeUpTimer:       NewTimer(),
				minCurrent:        minA,
				maxCurrent:        maxA,
				phases:            1,
				measuredPhases:    1,
				mode:              api.ModePV,
				batteryBoost:      boostContinue,
				batteryBoostLimit: 100, // disabled
				status:            api.StatusC,
				enabled:           true,
				offeredCurrent:    0,
			}

			uiChan, pushChan, lpChan := createChannels(t)

			// Prepare expects Enabled + MaxCurrent during init
			mockCharger.EXPECT().Enabled().Return(true, nil)
			if useChargerEx {
				// ChargerEx: Prepare calls setLimit which uses MaxCurrentMillis
			} else {
				mockCharger.EXPECT().MaxCurrent(int64(minA)).Return(nil)
			}
			lp.Prepare(site, uiChan, pushChan, lpChan)

			// simulate that the charger is already running at initialCurrent
			lp.offeredCurrent = initialCurrent

			// replicate sitePower computation including priorityAdjustment (site.go)
			batteryPower := site.battery.Power
			residualPower := site.GetResidualPower()
			var priorityAdjustment float64
			if len(site.batteryMeters) > 0 && site.battery.Soc < site.prioritySoc && residualPower <= 0 {
				priorityAdjustment -= 100 - residualPower
				residualPower = 100
			}
			if len(site.batteryMeters) > 0 && site.battery.Soc < site.prioritySoc && batteryPower < 0 {
				priorityAdjustment += batteryPower
				batteryPower = 0
			}
			sitePower := 0.0 + batteryPower + residualPower // grid=0

			// replicate batteryBoostPower as site.update() passes it
			batteryBoostPower := max(0.0, site.battery.Power)

			// undo battery priority adjustment for boost-active loadpoints (site.go update)
			lpSitePower := sitePower + priorityAdjustment

			// mock the charger calls that lp.Update will make
			mockCharger.EXPECT().Status().Return(api.StatusC, nil)
			mockCharger.EXPECT().Enabled().Return(true, nil)

			// capture the current set on the charger
			var setCurrent float64
			if useChargerEx {
				// MaxCurrentMillis is captured via the wrapper
			} else {
				mockCharger.EXPECT().MaxCurrent(gomock.Any()).DoAndReturn(func(c int64) error {
					setCurrent = float64(c)
					return nil
				}).AnyTimes()
			}

			lp.Update(lpSitePower, batteryBoostPower, nil, nil, false, false, 0, nil, nil)

			if useChargerEx {
				setCurrent = wrapper.current
			}
			currents = append(currents, setCurrent)
		})
	}

	return currents[0], currents[1]
}

func TestBatteryBoostPowerConsistency(t *testing.T) {
	// Battery charging at 2000W from PV, charger starting from 0A.
	// Battery boost should ramp up the charger identically regardless of SoC vs prioritySoc.
	above, below := boostTestCurrents(t, -2000, 0, false)
	assert.Equal(t, above, below,
		"battery boost must produce the same target current above and below prioritySoc, "+
			"got above=%.1fA, below=%.1fA", above, below)
}

func TestBatteryBoostPVExceedsMaxChargePower(t *testing.T) {
	// PV generation exceeds current charge power: charger at 14A, battery absorbing
	// the remaining 1000W surplus. Battery boost should increase the charger current
	// identically in both cases. Below prioritySoc the bug causes a lower target current
	// because sitePower hides the battery's surplus while batteryBoostPower can't compensate.
	above, below := boostTestCurrents(t, -1000, 14, false)
	assert.Equal(t, above, below,
		"battery boost must produce the same target current above and below prioritySoc, "+
			"got above=%.1fA, below=%.1fA", above, below)
}

func TestBatteryBoostDischarging(t *testing.T) {
	// Battery actively discharging at 2000W during boost. With sub-amp precision (ChargerEx),
	// the 100W residualPower override leak below prioritySoc becomes visible as a ~0.4A
	// difference that integer-amp chargers would mask.
	above, below := boostTestCurrents(t, 2000, 6, true)
	assert.InDelta(t, above, below, 0.01,
		"battery boost during discharge must produce the same target current above and below prioritySoc, "+
			"got above=%.1fA, below=%.1fA", above, below)
}

func TestRequiredBatteryMode(t *testing.T) {
	tc := []struct {
		gridChargeActive bool
		mode, res        api.BatteryMode
	}{
		{false, api.BatteryUnknown, api.BatteryUnknown}, // ignore
		{false, api.BatteryNormal, api.BatteryUnknown},  // ignore
		{false, api.BatteryHold, api.BatteryNormal},
		{false, api.BatteryCharge, api.BatteryNormal},

		{true, api.BatteryUnknown, api.BatteryCharge},
		{true, api.BatteryNormal, api.BatteryCharge},
		{true, api.BatteryHold, api.BatteryCharge},
		{true, api.BatteryCharge, api.BatteryUnknown}, // ignore
	}

	{
		// no battery
		res := new(Site).requiredBatteryMode(true, api.Rate{})
		assert.Equal(t, api.BatteryUnknown, res, "expected %s, got %s", api.BatteryUnknown, res)
	}

	for _, tc := range tc {
		t.Logf("%+v", tc)

		s := &Site{
			batteryMeters: []config.Device[api.Meter]{nil},
			batteryMode:   tc.mode,
		}

		res := s.requiredBatteryMode(tc.gridChargeActive, api.Rate{})
		assert.Equal(t, tc.res, res, "expected %s, got %s", tc.res, res)
	}
}
