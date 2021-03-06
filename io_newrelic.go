package main

import (
	"flag"
	"github.com/yvasiyarov/newrelic_platform_go"
	"log"
	// "net/http"
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

	iotop_processor := NewProcessIOProcessor()
	NewDynamicCollector(iotop_processor).Run()

	iostat_processor := NewDeviceStatsProcessor()
	NewDynamicCollector(iostat_processor).Run()

	component := NewDynamicPluginComponent(hostname, "com.github.maciejmrowiec.io_newrelic")

	plugin := newrelic_platform_go.NewNewrelicPlugin("0.0.2", "dfa8fd72df76280f63342a77af576f97023e7f74", 60)
	plugin.AddComponent(component)

	// Collect IO usage per process
	component.AddDynamicMetrica(NewTotalIOPerCommand(iotop_processor, "io/process/total_io_percentage"))
	component.AddDynamicMetrica(NewReadRatePerCommand(iotop_processor, "io/process/read_rate"))
	component.AddDynamicMetrica(NewWriteRatePerCommand(iotop_processor, "io/process/write_rate"))
	component.AddDynamicMetrica(NewSwapinPerCommand(iotop_processor, "io/process/swapin_percentage"))

	// Collect IO stats per device
	component.AddDynamicMetrica(NewRrqmpsPerDevice(iostat_processor, "io/device/rrqmps"))
	component.AddDynamicMetrica(NewWrqmpsPerDevice(iostat_processor, "io/device/wrqmps"))
	component.AddDynamicMetrica(NewRpsPerDevice(iostat_processor, "io/device/rps"))
	component.AddDynamicMetrica(NewWpsPerDevice(iostat_processor, "io/device/wps"))
	component.AddDynamicMetrica(NewRkbpsPerDevice(iostat_processor, "io/device/rkbps"))
	component.AddDynamicMetrica(NewWkbpsPerDevice(iostat_processor, "io/device/wkbps"))
	component.AddDynamicMetrica(NewAvgrqszPerDevice(iostat_processor, "io/device/avgrq_sz"))
	component.AddDynamicMetrica(NewAvgquszPerDevice(iostat_processor, "io/device/avgqu_sz"))
	component.AddDynamicMetrica(NewAwaitPerDevice(iostat_processor, "io/device/await"))
	component.AddDynamicMetrica(NewRawaitPerDevice(iostat_processor, "io/device/r_await"))
	component.AddDynamicMetrica(NewWawaitPerDevice(iostat_processor, "io/device/w_await"))
	component.AddDynamicMetrica(NewSvctmPerDevice(iostat_processor, "io/device/svctm"))
	component.AddDynamicMetrica(NewUtilPerDevice(iostat_processor, "io/device/util"))

	plugin.Verbose = verbose
	plugin.Run()
}
