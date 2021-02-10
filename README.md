moxspec-ocpsel
===

[![CircleCI](https://circleci.com/gh/actapio/moxspec-ocpsel.svg?style=shield&circle-token=461cacb242f68ee6efb7ec5b15b7b2984ccdd49a461cacb242f68ee6efb7ec5b15b7b2984ccdd49a)](https://circleci.com/gh/actapio/moxspec-occpsel)
[![Maintainability](https://api.codeclimate.com/v1/badges/95c9c7e8b41f759331d0/maintainability)](https://codeclimate.com/repos/6023e424911fcc41ac002f11/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/95c9c7e8b41f759331d0/test_coverage)](https://codeclimate.com/repos/6023e424911fcc41ac002f11/test_coverage)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

An OCP SEL decoder

## Installation

```
$ make bin
$ bin/ocpsel
```

## How to use

```
$ ocpsel -s 40 -g 0001 -e ab0700 
Machine Chk Err
  ed1 : 0xAB, 1010 1011 ... Uncorrectable
  ed2 : 0x07, 0000 0111 ... Machine Check bank Number: 7
  ed3 : 0x00, 0000 0000 ... CPU Number: 0, Core Number: 0
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
 Machine Chk Err
   ed1 : 0xAB, 1010 1011 ... Uncorrectable
   ed2 : 0x07, 0000 0111 ... Machine Check bank Number: 7
   ed3 : 0x00, 0000 0000 ... CPU Number: 0, Core Number: 0
```
