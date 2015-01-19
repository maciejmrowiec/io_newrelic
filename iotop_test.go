package main

import (
	"fmt"
	"math"
	"runtime"
	"testing"
)

// For floating point testing
var Epsilon float64

func init() {
	runtime.GOMAXPROCS(1)
	Epsilon = math.Nextafter(1.0, 2.0) - 1.0
}

func IsFloat64Equal(a float64, b float64) bool {
	return math.Abs(a-b) < Epsilon
}

func Test_parse_normal(t *testing.T) {
	str := "13:05:19 24086 be/4 cfapache    0.01 K/s    0.02 K/s  0.03 %  0.04 % httpd -k start"

	pr, err := NewProcessIOFromString(str)
	if err != nil {
		t.FailNow()
	}

	if !IsFloat64Equal(pr.disk_read_rate, 0.01) ||
		!IsFloat64Equal(pr.disk_write_rate, 0.02) ||
		!IsFloat64Equal(pr.swapin_percent, 0.03) ||
		!IsFloat64Equal(pr.io_percent, 0.04) ||
		pr.Name != "httpd -k start" {
		t.FailNow()
	}
}

func Test_parse_corrupted(t *testing.T) {
	str := "13:05:19 24086 be/4 cfapache    0.01    0.02 K/s  0.03 %  0.04 % httpd -k start"

	_, err := NewProcessIOFromString(str)
	if err == nil {
		t.FailNow()
	}
}

func Test_parse_empty(t *testing.T) {
	str := ""

	_, err := NewProcessIOFromString(str)
	if err == nil {
		t.FailNow()
	}
}

func Test_parse_header(t *testing.T) {
	if IsSampleSummary("lala") {
		t.Fail()
	}

	if !IsSampleSummary("Total DISK READ:") {
		t.Fail()
	}
}

func Test_get_and_purge(t *testing.T) {

	top := &IOTopCollector{}
	command_pipe := make(chan string, 1000)
	command_pipe <- "13:05:19 24086 be/4 cfapache    1.1 K/s   2.2 K/s  3.3 %  4.4 % test"
	command_pipe <- "Total DISK READ:"
	command_pipe <- "13:05:19 24086 be/4 cfapache    1.1 K/s   2.2 K/s  3.3 %  4.4 % test"
	command_pipe <- "Total DISK READ:"
	close(command_pipe)

	top.processOutput(command_pipe)

	if average := top.GetAndPurgeIOPercent(); !IsFloat64Equal(average["test"].GetAverage(), 4.4) {
		t.Fail()
	}

	if average := top.GetAndPurgeSwapinPercent(); !IsFloat64Equal(average["test"].GetAverage(), 3.3) {
		t.Fail()
	}

	if average := top.GetAndPurgeDiskReadRate(); !IsFloat64Equal(average["test"].GetAverage(), 1.1) {
		t.Fail()
	}

	if average := top.GetAndPurgeDiskWriteRate(); !IsFloat64Equal(average["test"].GetAverage(), 2.2) {
		t.Fail()
	}

}

func Benchmark_Processing(b *testing.B) {

	for n := 0; n < b.N; n++ {
		top := &IOTopCollector{}
		command_pipe := make(chan string, 1000)

		go func(ch chan<- string) {
			for i := 0; i < 50; i++ {
				head := "Total DISK READ:"
				ch <- head
				for i := 1; i <= 150; i++ {
					row := fmt.Sprintf("13:05:19 24086 be/4 cfapache    %d K/s    %d K/s  0.01 x  0.01 x process %d", i, i, i)
					ch <- row
				}
			}

			close(ch)
		}(command_pipe)

		top.processOutput(command_pipe)
	}
}
