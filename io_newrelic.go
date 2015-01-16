package main

import (
	"flag"
	"github.com/yvasiyarov/newrelic_platform_go"
	"log"
	"net/http"
	// _ "net/http/pprof"
	"os"
)

func main() {

	// for profiling
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	var err error
	var verbose bool
	var newrelic_key string
	flag.StringVar(&newrelic_key, "key", "", "Newrelic license key")
	flag.BoolVar(&verbose, "v", false, "Verbose mode")

	flag.Parse()

	if newrelic_key == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	top := &IOTopCollector{}
	top.Run()

	component := NewDynamicPluginComponent(hostname, "com.github.maciejmrowiec.io_newrelic")

	plugin := newrelic_platform_go.NewNewrelicPlugin("0.0.1", "dfa8fd72df76280f63342a77af576f97023e7f74", 60)
	plugin.AddComponent(component)

	component.AddDynamicMetrica(NewTotalIOPerCommand(top, "io/process/total_io_percentage"))
	component.AddDynamicMetrica(NewReadRatePerCommand(top, "io/process/read_rate"))
	component.AddDynamicMetrica(NewWriteRatePerCommand(top, "io/process/write_rate"))

	plugin.Verbose = verbose
	plugin.Run()
}
