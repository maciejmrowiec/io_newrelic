# io_newrelic
NewRelic plugin for measuring I/O usage per process.

Requirements:
* iotop installed in the system.
* Currenlty supports only CentsOS (iotop cmd options differ between OS)

Build instructions:

```

go get github.com/yvasiyarov/newrelic_platform_go
go get github.com/maciejmrowiec/io_newrelic
go build

```