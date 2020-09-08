# ports2packets

Generates a CSV file with the packets that are to be sent in order to check if a UDP port is open.

## Overview

`ports2packets` is a tool that is meant to solve the following problem: Knowing the packets that need to be sent to a UDP port in order to check whether it is open. This is done by UDP scanning Go UDP listeners on registered ports with [nmap](https://nmap.org/) and taking note of the port and packet, which are then writen in base64 format to a CSV file. The output CSV file based on the latest service names and port numbers by the [IANA](https://www.iana.org/) is also built weekly and available to download below.

## Installation

### Prebuilt Binaries

Linux, macOS and Windows binaries are available on [GitHub Releases](https://github.com/pojntfx/ports2packets/releases).

### Go Package

A Go package [is available](https://pkg.go.dev/mod/github.com/pojntfx/ports2packets).

### Prebuilt CSV File

As mentioned above, the CSV file is also pre-built every week and can be downloaded from [GitHub Releases](https://github.com/pojntfx/ports2packets/releases).

## Usage

```bash
% ports2packets -help
Usage of ports2packets:
  -in string
        Path to the CSV input file containing the registered services. Download from https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml (default "/etc/ports2packets/service-names-port-numbers.csv")
  -out string
        Path of the generated CSV output file (default "ports2packets.csv")
```

## License

ports2packets (c) 2020 Felicitas Pojtinger

SPDX-License-Identifier: AGPL-3.0
