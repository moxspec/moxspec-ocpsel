moxspec-ocpsel
===

[![CircleCI](https://circleci.com/gh/moxspec/moxspec-ocpsel.svg?style=shield&circle-token=eba3eea470549e9eb8de10b6275735e12c622ab3)](https://circleci.com/gh/moxspec/moxspec-ocpsel)
[![Maintainability](https://api.codeclimate.com/v1/badges/4a615055a788795e5384/maintainability)](https://codeclimate.com/github/moxspec/moxspec-ocpsel/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/4a615055a788795e5384/test_coverage)](https://codeclimate.com/github/moxspec/moxspec-ocpsel/test_coverage)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

An OCP SEL decoder

# Why needed / Story behind MoxSpec
1. Operation at scale with [Open Compute Project (OCP)](https://www.opencompute.org/) / White boxes

When you adopt Open Compute, instead of getting the benefit out of "part level" replacement for a better standardization and operational efficiency, the onsite service team might face to a challenge in the troubleshooting to identify what component needs to be replaced. [MoxSpec-OCPSEL](https://github.com/moxspec/moxspec-ocpsel) is developed and allows SEL (System Event Log) decoded to human readable information, which allows onsite service team easily identify what to be replaced such as a DIMM slot.
Democratizing the chance to gain another level of operational scalability to us by adopting open technologies like OCP was crucial and a holistic approach to even a 19” OEM servers was needed.

## Installation

```
$ make bin
$ bin/ocpsel
```

## How to use

```
$ ocpsel -s 40 -g 0001 -e ab0700 
Decoded Info:
  Summary      : Machine Chk Err
  Event Data 1 : 0xAB, 1010 1011 ... Uncorrectable
  Event Data 2 : 0x07, 0000 0111 ... Machine Check bank Number: 7
  Event Data 3 : 0x00, 0000 0000 ... CPU Number: 0, Core Number: 0
  Generator    : 0x0001          ... BIOS/UEFI system Firmware
```

```
$ sudo ipmitool sel elist -v > selelistv.log
$ ocpsel -f selelistv.log
SEL Record ID          : 00a1
 Record Type           : 02
 Timestamp             : 08/08/2019 10:37:17
 Generator ID          : 0001
 EvM Revision          : 04
 Sensor Type           : Processor
 Sensor Number         : 40
 Event Type            : Sensor-specific Discrete
 Event Direction       : Assertion Event
 Event Data            : ab0700
 Description           : Uncorrectable machine check exception
 Decoded Info:
   Summary      : Machine Chk Err
   Event Data 1 : 0xAB, 1010 1011 ... Uncorrectable
   Event Data 2 : 0x07, 0000 0111 ... Machine Check bank Number: 7
   Event Data 3 : 0x00, 0000 0000 ... CPU Number: 0, Core Number: 0
   Generator    : 0x0001          ... BIOS/UEFI system Firmware
```
