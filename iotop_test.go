package main

import (
	"testing"
)

func Test_processio_parse(t *testing.T) {
	str := "13:05:19 24086 be/4 cfapache    0.01 K/s    0.02 K/s  0.03 %  0.04 % httpd -k start"

	p, err := ProcessIOParseRow(str)
	process := p.(*ProcessIO)

	if err != nil {
		t.FailNow()
	}

	if !IsFloat64Equal(process.disk_read_rate, 0.01) {
		t.FailNow()
	}

	if !IsFloat64Equal(process.disk_write_rate, 0.02) {
		t.FailNow()
	}

	if !IsFloat64Equal(process.swapin_percent, 0.03) {
		t.FailNow()
	}

	if !IsFloat64Equal(process.io_percent, 0.04) {
		t.FailNow()
	}
}

func Test_processio_header(t *testing.T) {
	if !ProcessIOIsHeader("Total DISK READ:") {
		t.Fail()
	}
}

func Test_processio_aggregate_sample(t *testing.T) {
	str_1 := "13:05:19 24086 be/4 cfapache    0.1 K/s   0.2 K/s  0.3 %  0.4 % test"
	str_2 := "13:05:19 24086 be/4 cfapache    0.2 K/s   0.3 K/s  0.4 %  0.5 % test"

	processor := NewProcessIOProcessor()

	sample_1 := NewSample()
	item_1, _ := processor.ParseRow(str_1)
	sample_1.Append(item_1)

	processor.AggregateSample(sample_1)

	if !IsFloat64Equal(processor.disk_read_rate["test"].GetAverage(), 0.1) {
		t.FailNow()
	}

	if !IsFloat64Equal(processor.disk_write_rate["test"].GetAverage(), 0.2) {
		t.FailNow()
	}

	if !IsFloat64Equal(processor.swapin_percent["test"].GetAverage(), 0.3) {
		t.FailNow()
	}

	if !IsFloat64Equal(processor.io_percent["test"].GetAverage(), 0.4) {
		t.FailNow()
	}

	sample_2 := NewSample()
	item_2, _ := processor.ParseRow(str_2)
	sample_2.Append(item_2)

	processor.AggregateSample(sample_2)

	if !IsFloat64Equal(processor.disk_read_rate["test"].GetAverage(), 0.15) {
		t.FailNow()
	}

	if !IsFloat64Equal(processor.disk_write_rate["test"].GetAverage(), 0.25) {
		t.FailNow()
	}

	if !IsFloat64Equal(processor.swapin_percent["test"].GetAverage(), 0.35) {
		t.FailNow()
	}

	if !IsFloat64Equal(processor.io_percent["test"].GetAverage(), 0.45) {
		t.FailNow()
	}
}
