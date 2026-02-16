# perdoon

Port sniffer. Tool used to track a wide range of ports, emulate responses and store data for future analysis.


## What is it?

A simple tool that can listen on a range of UDP or TCP ports, record incoming data and respond using different strategies.

This tool starts UDP and TCP servers on every port you defined. Then it records access and activity log for each port and stores it inside database file. It also can generate reponse data using some pre-defined strategies.


## Setup

Build:

```bash
git clone https://github.com/bitrate16/perdoon.git
cd perdoon/src

go build -o perdoon
```

Install using root:

```bash
sudo cp -f perdoon /usr/bin/perdoon
```


## Configuration

Sample configuration:

```yml
tcp:
  ports:
    - 1-65535
  chunk-size: 1024
udp:
  ports:
    - 1-1000
    - 2000-15000
  chunk-size: 1024
database:
  path: perdoon.db
  table: perdoon
  record:
    request-payload: true
    response-payload: true
response:
  sizes:
    - 1-16
    - 64-128
    - 512-65536
  bytes: 3a33
strategy: random
```

## Port range

Port ranges for each prtocol are defined in `udp.ports` and `tcp.ports`. You can also define a wide range of ports and exclude vcertain ports with `udp.exclude` and `tcp.exclude` sections.


## Strategies

Strategy defines how to respond to clients.

Supported strategies are:
- echo - echo data back without change
- random - produce random data
  - `response.sizes` - random size ranges
- zero - respond with zero bytes
  - `response.sizes` - random size ranges
- potato - respond with "potato" string repeated
  - `response.sizes` - random size ranges
- bytes - respond with fixed string repeated
  - `response.sizes` - random size ranges
  - `response.bytes` - hex string. If non even, padded with zero at left
- qt - respond with ":3" string repeated
  - `response.sizes` - random size ranges

Strategiy parameters can be controlled with `response` section. For example you can generate specific sizes of responses randomly.


## Database

`database` section defines the configuration of the database like path and table name. Currently this tool uses sqlite3 to store records.

`database.record` section defines what extra data to record. You can disable recording payload data to save space.



## Privilleged ports

In order to bind privilleged ports in range 1-1023 process should be run as root or have certain capabilities.

You can allow executable to bind these ports by setting `CAP_NET_BIND_SERVICE`:

```bash
sudo setcap 'CAP_NET_BIND_SERVICE+ep' /usr/bin/perdoon
```


## ulimit

Each process has certain limits (like open files, handles, memory and others). Opening too many ports may not be allows by current limit.

You can check the limits by executing:

```bash
ulimit -a
```


# LICENSE

```
perdoon - port sniffer
Copyright (C) 2026  bitrate16

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```
