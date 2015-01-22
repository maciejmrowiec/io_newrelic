
[![Build Status](https://drone.io/github.com/maciejmrowiec/io_newrelic/status.png)](https://drone.io/github.com/maciejmrowiec/io_newrelic/latest)

# io_newrelic
NewRelic plugin for measuring I/O usage.

Metricas:
* /io/process/ - I/O usage per process (iotop based measurements) with 1s sample rate
* /io/device/ - I/O stats per device (iostat based measurements)

Requirements:
* iotop installed in the system and on the PATH.
* iostat installed in the system and on the PATH.
* Tested on CentOS and Ubuntu (OS X not supported)

Build instructions:

```

go get github.com/yvasiyarov/newrelic_platform_go
go get github.com/maciejmrowiec/io_newrelic
go build

```
