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

func Test_parse_normal(t *testing.T) {
	str := "13:05:19 24086 be/4 cfapache    0.01 K/s    0.02 K/s  0.03 %  0.04 % httpd -k start"

	pr, err := NewProcessIOFromString(str)
	if err != nil {
		t.FailNow()
	}

	if math.Abs(pr.disk_read_rate-0.01) >= Epsilon ||
		math.Abs(pr.disk_write_rate-0.02) >= Epsilon ||
		math.Abs(pr.swapin_percent-0.03) >= Epsilon ||
		math.Abs(pr.io_percent-0.04) >= Epsilon ||
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
