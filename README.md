moxspec-ocpsel
===

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
