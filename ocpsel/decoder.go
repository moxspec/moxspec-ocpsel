package ocpsel

import (
	"fmt"
	"strconv"
)

// Decode decodes given sensor number, generator number, and evend data
func Decode(sens, gens, eds string) (record, error) {
	sen, err := strconv.ParseUint(sens, 16, 32)
	if err != nil {
		return record{}, err
	}

	gen, err := strconv.ParseUint(gens, 16, 32)
	if err != nil {
		return record{}, err
	}

	ed, err := strconv.ParseUint(eds, 16, 32)
	if err != nil {
		return record{}, err
	}

	title, d, err := getDecoder(sen, gen)
	if err != nil {
		return record{}, err
	}

	ed1, ed2, ed3 := splitEventData(ed)
	r1, r2, r3 := d(ed1, ed2, ed3)

	return record{
		title: title,
		ed1:   ed1,
		r1:    r1,
		ed2:   ed2,
		r2:    r2,
		ed3:   ed3,
		r3:    r3,
	}, nil
}

// 8.10.3 Events
//
// Events are triggered if a GPIO transition is detected. The Event only discrete sensors that
// are required and the Event Data format is provided in Table 8-4.
// Generator ID in the event log shows which piece of firmware generates the log:
//
//  0x602C = Intel® SPS ME Firmware
//  0x0001 = BIOS/UEFI system Firmware
//  0x0020 = BMC Firmware
const (
	intel = 0x602C
	bios  = 0x0001
	bmc   = 0x0020
)

type decoder func(ed1, ed2, ed3 byte) (r1, r2, r3 string)

func getDecoder(sen, gen uint64) (string, decoder, error) {
	match := func(s, g uint64) bool {
		if sen == s && gen == g {
			return true
		}
		return false
	}

	var (
		title string
		d     decoder
	)
	switch {
	case match(0x17, intel):
		title = "SPS FW Health"
		d = decodeSPSFWHealth
	case match(0x18, intel):
		title = "NM Exception"
		d = decodeNMException
	case match(0x19, intel):
		title = "NM Health"
		d = decodeNMHealth
	case match(0x1A, intel):
		title = "NM Capabilities"
		d = decodeNMCapabilities
	case match(0x1B, intel):
		title = "NM Threshold"
		d = decodeNMThreshold
	case match(0x1C, intel):
		title = "CPU0 Therm Statu"
		d = decodeCPUThermStatu
	case match(0x1D, intel):
		title = "CPU1 Therm Statu"
		d = decodeCPUThermStatu
	case match(0x2B, bios):
		title = "POST Error"
		d = decodePOSTError
	case match(0x3B, intel):
		title = "Pwr Thresh Evt"
		d = decodePwrThreshEvt
	case match(0x40, bios):
		title = "Machine Chk Err"
		d = decodeMachineChkErr
	case match(0x41, bios):
		title = "PCIe Error"
		d = decodePCIeError
	case match(0x43, bios):
		title = "Other IIO Err"
		d = decodeOtherIIOErr
	case match(0x51, bmc):
		title = "ProcHot Ext"
		d = decodeProcHotExt
	case match(0x52, bmc):
		title = "MemHot Ext"
		d = decodeMemHotExt
	case match(0x56, bmc):
		title = "Power Error"
		d = decodePowerError
	case match(0x63, bios):
		title = "Memory ECC Error"
		d = decodeMemoryECCError
	case match(0x90, bios):
		title = "Software NMI"
		d = decodeSoftwareNMI
	case match(0xAA, bmc):
		title = "Button"
		d = decodeButton
	case match(0xAB, bmc):
		title = "Power State"
		d = decodePowerState
	case match(0xAC, bmc):
		title = "Power Policy"
		d = decodePowerPolicy
	case match(0xAE, bmc):
		title = "ME Status"
		d = decodeMEStatus
	case match(0xB1, bmc):
		title = "Network Status"
		d = decodeNetworkStatus
	case match(0xBF, bmc):
		title = "PCH Thermal Trip"
		d = decodePCHThermalTrip
	case match(0xC5, intel):
		title = "ME Gl Reset Warn"
		d = decodeMEGlResetWarn
	case match(0xE9, bmc):
		title = "System Event"
		d = decodeSystemEvent
	case match(0xEA, bmc):
		title = "Critical IRQ"
		d = decodeCriticalIRQ
	case match(0xEB, bmc):
		title = "CATERR/MS MI"
		d = decodeCATERRorMSMI
	case match(0xEF, bmc):
		title = "Dual BIOS Up Sts"
		d = decodeDualBIOSUpSts
	}

	if d == nil {
		return "", nil, fmt.Errorf("unsupported pattern")
	}

	return title, d, nil
}

func decodeSPSFWHealth(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	var et string
	switch mustPick(ed1, 0, 3) {
	case 0x00:
		et = "Firmware Status"
	}

	r1 = fmt.Sprintf("ED2 OEM code: %d, ED3 OEM code: %d, Health Event Type: %s", mustPick(ed1, 6, 7), mustPick(ed1, 4, 5), et)
	r2 = "Follow the Intel® SPS FW specification"
	r3 = "Follow the Intel® SPS FW specification"

	return
}

func decodeNMException(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	r1 = fmt.Sprintf("ED2 OEM code: %d, ED3 OEM code: %d", mustPick(ed1, 4, 5), mustPick(ed1, 6, 7))
	if mustPick(ed1, 3, 3) == 1 {
		r1 = fmt.Sprintf("%s, Policy Correction Time Exceeded", r1)
	}

	switch ed2 {
	case 0x00:
		r2 = "Entire platform"
	case 0x01:
		r2 = "CPU subsystem"
	case 0x02:
		r2 = "Memory subsystem"
	case 0x03:
		r2 = "HW Protection"
	case 0x04:
		r2 = "High Power I/O subsystem."
	}

	r3 = fmt.Sprintf("Policy ID: %d", ed3)
	return
}

func decodeNMHealth(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x02:
		r1 = "Sensor Intel® Node Manager"
	}

	r2 = fmt.Sprintf("0x%02X, Follow the Intel® SPS FW specification", ed2)
	r3 = fmt.Sprintf("0x%02X, Follow the Intel® SPS FW specification", ed3)
	return
}

func decodeNMCapabilities(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	pol := "Available"
	if mustPick(ed1, 0, 0) == 0 {
		pol = "Not Available"
	}

	mon := "Available"
	if mustPick(ed1, 1, 1) == 0 {
		mon = "Not Available"
	}

	pow := "Available"
	if mustPick(ed1, 2, 2) == 0 {
		pow = "Not Available"
	}

	r1 = fmt.Sprintf("Policy interface capability: %s, Monitoring capability: %s, Power limiting capability: %s", pol, mon, pow)
	return
}

func decodeNMThreshold(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	var reason string
	switch mustPick(ed1, 3, 3) {
	case 0x00:
		reason = "Threshold exceeded"
	case 0x01:
		reason = "Policy Correction Time Exceeded"
	}
	r1 = fmt.Sprintf("Threshold Number: %d, %s", mustPick(ed1, 0, 1), reason)

	switch ed2 {
	case 0x00:
		r2 = "Entire platform"
	case 0x01:
		r2 = "CPU subsystem"
	case 0x02:
		r2 = "Memory subsystem"
	case 0x03:
		r2 = "HW Protection"
	case 0x04:
		r2 = "High Power I/O subsystem"
	}

	r3 = fmt.Sprintf("Policy ID: %d", ed3)

	return
}

func decodeCPUThermStatu(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch ed1 {
	case 0x00:
		r1 = "CPU Critical Temperature. Indicates whether CPU temperature is above critical temperature point."
	case 0x01:
		r1 = "PROCHOT# Assertions. Indicates whether PROCHOT# signal is asserted."
	case 0x02:
		r1 = "TCC Activation. Indicates whether CPU thermal throttling functionality is activated due to CPU temperature being above Thermal Circuit Control Activation point."
	}

	return
}

func decodePOSTError(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x00:
		r1 = "System Firmware Error"
	}

	switch mustPick(ed1, 6, 7) {
	case 0x02:
		r2 = fmt.Sprintf("LSB of OEM POST Error Code: 0x%02X", ed2)
	case 0x03:
		r2 = "Per IPMI Spec"
	}

	r3 = fmt.Sprintf("MSB of OEM POST Error Code: 0x%02X", ed3)
	return
}

func decodePwrThreshEvt(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	if ed1 == 0x01 {
		r1 = "Limit Exceeded"
	}
	return
}

func decodeMachineChkErr(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x0B:
		r1 = "Uncorrectable"
	case 0x0C:
		r1 = "Correctable"
	}

	r2 = fmt.Sprintf("Machine Check bank Number: %d", ed2)
	r3 = fmt.Sprintf("CPU Number: %d, Core Number: %d", mustPick(ed3, 5, 7), mustPick(ed3, 0, 4))

	return
}

func decodePCIeError(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x04:
		r1 = "PCI PERR"
	case 0x05:
		r1 = "PCI SERR"
	case 0x07:
		r1 = "correctable"
	case 0x08:
		r1 = "uncorrectable"
	case 0x0A:
		r1 = "Bus Fatal"
	}

	r2 = fmt.Sprintf("Device Number: 0x%02X, Func Number: 0x%02X", mustPick(ed2, 3, 7), mustPick(ed2, 0, 2))
	r3 = fmt.Sprintf("Bus Number: 0x%02X", ed3)

	return
}

func decodeOtherIIOErr(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x00:
		r1 = "Offset 0x00 (Other IIO)"
	}

	r2 = fmt.Sprintf("Error ID: 0x%02X [Refer to Intel® Xeon® Processor E5 v3 Product Family Datasheet, Vol. 1 Sec 11.1.7 IIO module error codes]", ed2)

	var src string
	switch mustPick(ed3, 0, 2) {
	case 0x00:
		src = "IRP0"
	case 0x01:
		src = "IRP1"
	case 0x02:
		src = "IIO- Core"
	case 0x03:
		src = "VT-d"
	case 0x04:
		src = "TBD"
	case 0x05:
		src = "Misc"
	}

	r3 = fmt.Sprintf("CPU#: %d, Source: %s", mustPick(ed3, 5, 7), src)

	return
}

func decodeProcHotExt(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x0A:
		r1 = "Processor thermal throttling offset"
	}

	switch mustPick(ed2, 0, 1) {
	case 0x00:
		r2 = "Native"
	case 0x01:
		r2 = "External (VR)"
	case 0x02:
		r2 = "External(Throttle)"
	}

	r3 = fmt.Sprintf("CPU/VR Number: %d", mustPick(ed3, 5, 7))

	return
}

func decodeMemHotExt(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x09:
		r1 = "Memory thermal throttling offset"
	}

	switch mustPick(ed2, 0, 1) {
	case 0x00:
		r2 = "Native"
	case 0x01:
		r2 = "External (VR)"
	case 0x02:
		r2 = "External(Throttle)"
	}

	r3 = fmt.Sprintf("CPU/VR Number: %d, Channel Number: %d, DIMM Number: %d", mustPick(ed3, 5, 7), mustPick(ed3, 3, 4), mustPick(ed3, 0, 2))

	return
}

func decodePowerError(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x01:
		r1 = "SYS_PWROK Failure"
	case 0x02:
		r1 = "PCH_PWROK Failure"
	}

	r2 = ed2ShouldBe(ed2, 0xFF)
	r3 = ed3ShouldBe(ed3, 0xFF)

	return
}

func decodeMemoryECCError(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x00:
		r1 = "Correctable"
	case 0x01:
		r1 = "Uncorrectable"
	case 0x05:
		r1 = "Correctable ECC error Logging Limit Reached"
	}

	switch mustPick(ed2, 2, 3) {
	case 0x00:
		r2 = "All info available"
	case 0x01:
		r2 = "DIMM info not valid"
	case 0x02:
		r2 = "CHN info not valid"
	case 0x03:
		r2 = "CPU info not valid"
	}
	r2 = fmt.Sprintf("%s, Logical Rank: %d", r2, mustPick(ed2, 0, 1))

	cpu := mustPick(ed3, 5, 7)
	channel := mustPick(ed3, 2, 4)
	dimm := mustPick(ed3, 0, 1)

	posNumber := func() int {
		switch channel {
		case 0:
			if dimm == 0 {
				return 0
			}
			return 1
		case 1:
			return 2
		case 2:
			return 3
		case 3:
			if dimm == 0 {
				return 4
			}
			return 5
		case 4:
			return 6
		case 5:
			return 7
		}
		return 0
	}

	posLetter := func() string {
		if cpu == 0 {
			return "A"
		}
		return "B"
	}

	r3 = fmt.Sprintf("Label: %s%d, CPU: %d, Channel: %d, DIMM: %d", posLetter(), posNumber(), cpu, channel, dimm)

	return
}

func decodeSoftwareNMI(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x03:
		r1 = "Software NMI offset (03h)"
	}

	r2 = ed2ShouldBe(ed2, 0xFF)
	r3 = ed3ShouldBe(ed3, 0xFF)

	return
}

func decodeButton(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x00:
		r1 = "Power button pressed"
	case 0x02:
		r1 = "Reset button pressed"
	}

	r2 = ed2ShouldBe(ed2, 0xFF)
	r3 = ed3ShouldBe(ed3, 0xFF)

	return
}

func decodePowerState(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x00:
		r1 = "Transition to running"
	case 0x02:
		r1 = "Transition to power off"
	}

	r2 = ed2ShouldBe(ed2, 0xFF)
	r3 = ed3ShouldBe(ed3, 0xFF)

	return
}

func decodePowerPolicy(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x05:
		r1 = "Soft-power control failure"
	}

	r2 = ed2ShouldBe(ed2, 0xFF)
	r3 = ed3ShouldBe(ed3, 0xFF)

	return
}

func decodeMEStatus(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x00:
		r1 = "Controller access degraded or unavailable"
	case 0x03:
		r1 = "Management controller unavailable"
	}

	r2 = ed2ShouldBe(ed2, 0xFF)
	r3 = ed3ShouldBe(ed3, 0xFF)

	return
}

func decodeNetworkStatus(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch mustPick(ed1, 0, 3) {
	case 0x00:
		r1 = "After BMC has IP assigned, and Both IPv4 and IPv6 network cannot ping gateway. BMC in disconnection state"
	case 0x01:
		r1 = "Either IPv4 or IPv6 can ping gateway after disconnection state"
	}

	r2 = ed2ShouldBe(ed2, 0xFF)
	r3 = ed3ShouldBe(ed3, 0xFF)

	return
}

func decodePCHThermalTrip(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	if ed1 == 0x01 {
		r1 = "State Asserted"
	}

	r2 = ed2ShouldBe(ed2, 0xFF)
	r3 = ed3ShouldBe(ed3, 0xFF)

	return
}

func decodeMEGlResetWarn(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	if ed1 != 0xA0 {
		r1 = fmt.Sprintf("ed1 should be 0xA0 but 0x%02X", ed1)
		return
	}

	r2 = "Time Units for which Intel® ME will delay Global Platform Reset. "
	if ed2 == 0xFF {
		r2 = r2 + "Infinite delay. For debug purposes, you could configure the delay time as infinity. In this case, the BMC is not required to respond to the event – global reset is suppressed unconditionally."
		return
	}
	r2 = r2 + "Time in unites specified in Event Data 3."

	switch ed3 {
	case 0x01:
		r3 = "minutes"
	}

	return
}

func decodeSystemEvent(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch ed1 {
	case 0xE5:
		r1 = "Timestamp Clock Synch"
	case 0xC4:
		r1 = "PEF Action"
	}

	if ed1 == 0xE5 {
		switch ed2 {
		case 0x00:
			r2 = "event is first of pair"
		case 0x80:
			r2 = "event is second of pair"
		}
	}

	if ed1 == 0xC4 {
		switch ed2 {
		case 0x01:
			r2 = "PEF Action"
		}

	}

	if ed1 == 0xE5 {
		switch ed3 {
		case 0x00:
			r3 = "Cause of time changed is NTP"
		case 0x01:
			r3 = "Cause of time changed is Host RTC"
		case 0x02:
			r3 = "Cause of time changed is Set SEL time Command"
		case 0x03:
			r3 = "Cause of time changed is Set SEL time UTC offset Command"
		case 0xFF:
			r3 = "Cause of time changed is Unknown"
		}
	}

	return
}

func decodeCriticalIRQ(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch ed1 {
	case 0x00:
		r1 = "Front Panel NMI / Diagnostic Interrupt"
	}

	r2 = ed2ShouldBe(ed2, 0xFF)
	r3 = ed3ShouldBe(ed3, 0xFF)

	return
}

func decodeCATERRorMSMI(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch ed1 {
	case 0x00:
		r1 = "IERR (CATERR_N Hold Low)"
	case 0x01:
		r1 = "MSMI_N Hold Low"
	case 0x0B:
		r1 = "MCERR (CATERR_N 16 BCLK pulse)"
	case 0x0C:
		r1 = "MSMI_N 16 BCLK Pulse"
	}

	r2 = ed2ShouldBe(ed2, 0xFF)
	r3 = ed3ShouldBe(ed3, 0xFF)

	return
}

func decodeDualBIOSUpSts(ed1, ed2, ed3 byte) (r1, r2, r3 string) {
	switch ed1 {
	case 0x01:
		r1 = "Auto Recovery"
	case 0x02:
		r1 = "Manual Recovery"
	case 0x03:
		r1 = "OOB Directly"
	case 0x04:
		r1 = "Auto Detect"
	case 0x05:
		r1 = "BIOS Crash by SLP_S3_N cycling Recovery"
	}

	if ed1 == 0x01 {
		switch ed2 {
		case 0x01:
			r2 = "FRB2 WDT timeout"
		case 0x02:
			r2 = "BIOS Good de-assert (GPION2)"
		case 0x07:
			r2 = "Watchdog not enable"
		}
	}

	if ed1 == 0x02 {
		switch ed2 {
		case 0x03:
			r2 = "Recovery from Gold to Primary"
		case 0x04:
			r2 = "Recovery from Primary to Gold"
		}
	}

	if ed1 == 0x03 {
		switch ed2 {
		case 0x05:
			r2 = "Primary directly"
		case 0x06:
			r2 = "Gold directly"
		}
	}

	if ed1 == 0x04 {
		switch ed2 {
		case 0x08:
			r2 = "BMC Self test failed"
		case 0x09:
			r2 = "Unexpected power off"
		case 0x0A:
			r2 = "BMC ready pin de-assert (GPIOQ4)"
		}
	}

	if ed1 == 0x05 {
		switch ed2 {
		case 0x03:
			r2 = "Recovery from Gold to Primary"
		}
	}

	if ed2 == 0x01 || ed2 == 0x02 {
		switch ed3 {
		case 0x01:
			r3 = "Start the progress for recovery"
		case 0x02:
			r3 = "End the progress for recovery"
		case 0x03:
			r3 = "Checksum compare failed"
		case 0x04:
			r3 = "Primary BIOS is not present"
		case 0x05:
			r3 = "Gold BIOS is not present"
		}
	}

	return
}
