package eebus

// EEBUS use case scenario numbers per the respective Use Case Technical Specifications.
//
// Spec scenario numbers diverge between use cases (e.g. MPC scenario 1 = active power,
// MGCP scenario 1 = power factor; MPC scenario 2 = energy, MGCP scenario 2 = active power).
// Passing the wrong number to IsScenarioAvailableAtEntity gates reads on the wrong feature.
//
// Each block mirrors the scenarios registered in the corresponding eebus-go usecase, which
// in turn matches the EEBus UC TS document.

// MGCP — Monitoring of Grid Connection Point (UC TS v1.0.0)
const (
	MGCPScenarioPowerFactor     uint = 1 // S1 power factor (cos phi)
	MGCPScenarioPower           uint = 2 // S2 active power per phase + total
	MGCPScenarioEnergyFeedIn    uint = 3 // S3 total feed-in energy
	MGCPScenarioEnergyConsumed  uint = 4 // S4 total consumed energy
	MGCPScenarioCurrentPerPhase uint = 5 // S5 phase-specific currents
	MGCPScenarioVoltagePerPhase uint = 6 // S6 phase-specific voltages
	MGCPScenarioFrequency       uint = 7 // S7 frequency
)

// MPC — Monitoring of Power Consumption (UC TS v1.0.0)
const (
	MPCScenarioPower           uint = 1 // S1 active power per phase + total
	MPCScenarioEnergyConsumed  uint = 2 // S2 total consumed energy
	MPCScenarioCurrentPerPhase uint = 3 // S3 phase-specific currents
	MPCScenarioVoltagePerPhase uint = 4 // S4 phase-specific voltages
	MPCScenarioFrequency       uint = 5 // S5 frequency
)

// LPC — Limitation of Power Consumption (UC TS v1.0.0). Same scenario layout for CS and EG roles.
const (
	LPCScenarioLimit                uint = 1 // S1 LoadControl: consumption limit
	LPCScenarioFailsafe             uint = 2 // S2 DeviceConfiguration: failsafe values
	LPCScenarioHeartbeat            uint = 3 // S3 DeviceDiagnosis: heartbeat
	LPCScenarioElectricalConnection uint = 4 // S4 ElectricalConnection (optional)
)

// LPP — Limitation of Power Production (UC TS v1.0.0). Same scenario layout for CS and EG roles.
const (
	LPPScenarioLimit                uint = 1 // S1 LoadControl: production limit
	LPPScenarioFailsafe             uint = 2 // S2 DeviceConfiguration: failsafe values
	LPPScenarioHeartbeat            uint = 3 // S3 DeviceDiagnosis: heartbeat
	LPPScenarioElectricalConnection uint = 4 // S4 ElectricalConnection (optional)
)

// OPEV — Overload Protection by EV Charging Current Curtailment (UC TS v1.0.1)
const (
	OPEVScenarioObligationLimit uint = 1 // S1 LoadControl + ElectricalConnection
	OPEVScenarioChargingState   uint = 2 // S2 charging state
	OPEVScenarioChargingPlan    uint = 3 // S3 charging plan
)

// OSCEV — Optimization of Self-Consumption during EV Charging (UC TS v1.0.1)
const (
	OSCEVScenarioRecommendationLimit uint = 1 // S1 LoadControl + ElectricalConnection
	OSCEVScenarioChargingState       uint = 2 // S2 charging state
	OSCEVScenarioChargingPlan        uint = 3 // S3 charging plan
)

// EVCEM — Measurement of Electricity during EV Charging (UC TS v1.0.1)
const (
	EVCEMScenarioPowerPerPhase uint = 1 // S1 phase-specific active power + ElectricalConnection (currents)
	EVCEMScenarioPowerTotal    uint = 2 // S2 total active power only
	EVCEMScenarioEnergy        uint = 3 // S3 charging energy summary
)

// EVSOC — EV State of Charge (UC TS v1.0.0 RC1)
const (
	EVSOCScenarioStateOfCharge uint = 1 // S1 state of charge
)
