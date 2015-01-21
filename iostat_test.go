package main

import (
	"testing"
)

func Test_devicestats_parse_ubuntu(t *testing.T) {
	str := "sda               0.10     0.20    0.30    0.40     0.50     0.60     0.70     0.80    0.90    1.00    1.10   1.20   1.30"

	d, err := DeviceStatsParseRow(str)
	device := d.(*DeviceStats)

	if err != nil {
		t.FailNow()
	}

	if device.name != "sda" {
		t.Fail()
	}

	if !IsFloat64Equal(device.rrqmps, 0.1) {
		t.Fail()
	}

	if !IsFloat64Equal(device.wrqmps, 0.2) {
		t.Fail()
	}

	if !IsFloat64Equal(device.rps, 0.3) {
		t.Fail()
	}

	if !IsFloat64Equal(device.wps, 0.4) {
		t.Fail()
	}

	if !IsFloat64Equal(device.rkbps, 0.5) {
		t.Fail()
	}

	if !IsFloat64Equal(device.wkbps, 0.6) {
		t.Fail()
	}

	if !IsFloat64Equal(device.avgrq_sz, 0.7) {
		t.Fail()
	}

	if !IsFloat64Equal(device.avgqu_sz, 0.8) {
		t.Fail()
	}

	if !IsFloat64Equal(device.await, 0.9) {
		t.Fail()
	}

	if !IsFloat64Equal(device.r_await, 1.0) {
		t.Fail()
	}
	if !IsFloat64Equal(device.w_await, 1.1) {
		t.Fail()
	}
	if !IsFloat64Equal(device.svctm, 1.2) {
		t.Fail()
	}
	if !IsFloat64Equal(device.util, 1.3) {
		t.Fail()
	}
}

func Test_devicestats_parse_centos(t *testing.T) {
	str := "sda               0.10     0.20    0.30    0.40     0.50     0.60     0.70     0.80    0.90    1.20   1.30"

	d, err := DeviceStatsParseRow(str)
	device := d.(*DeviceStats)

	if err != nil {
		t.FailNow()
	}

	if device.name != "sda" {
		t.Fail()
	}

	if !IsFloat64Equal(device.rrqmps, 0.1) {
		t.Fail()
	}

	if !IsFloat64Equal(device.wrqmps, 0.2) {
		t.Fail()
	}

	if !IsFloat64Equal(device.rps, 0.3) {
		t.Fail()
	}

	if !IsFloat64Equal(device.wps, 0.4) {
		t.Fail()
	}

	if !IsFloat64Equal(device.rkbps, 0.5) {
		t.Fail()
	}

	if !IsFloat64Equal(device.wkbps, 0.6) {
		t.Fail()
	}

	if !IsFloat64Equal(device.avgrq_sz, 0.7) {
		t.Fail()
	}

	if !IsFloat64Equal(device.avgqu_sz, 0.8) {
		t.Fail()
	}

	if !IsFloat64Equal(device.await, 0.9) {
		t.Fail()
	}

	if !IsFloat64Equal(device.r_await, 0.0) {
		t.Fail()
	}
	if !IsFloat64Equal(device.w_await, 0.0) {
		t.Fail()
	}
	if !IsFloat64Equal(device.svctm, 1.2) {
		t.Fail()
	}
	if !IsFloat64Equal(device.util, 1.3) {
		t.Fail()
	}
}

func Test_devicestats_header(t *testing.T) {
	str := "Device:         rrqm/s   wrqm/s     r/s     w/s    rkB/s    wkB/s avgrq-sz avgqu-sz   await r_await w_await  svctm  %util"

	if !DeviceStatsIsHeader(str) {
		t.FailNow()
	}
}

func Test_aggregate_sample(t *testing.T) {
	str_1 := "sda               0.10     0.20    0.30    0.40     0.50     0.60     0.70     0.80    0.90    1.00    1.10   1.20   1.30"
	str_2 := "sda               0.20     0.20    0.30    0.40     0.50     0.60     0.70     0.80    0.90    1.00    1.10   1.20   1.30"

	processor := NewDeviceStatsProcessor()

	sample_1 := NewSample()
	item_1, _ := processor.ParseRow(str_1)
	sample_1.Append(item_1)

	processor.AggregateSample(sample_1)

	if !IsFloat64Equal(processor.rrqmps["sda"].GetAverage(), 0.1) {
		t.FailNow()
	}

	sample_2 := NewSample()
	item_2, _ := processor.ParseRow(str_2)
	sample_2.Append(item_2)

	processor.AggregateSample(sample_2)

	if !IsFloat64Equal(processor.rrqmps["sda"].GetAverage(), 0.15) {
		t.FailNow()
	}
}

func Test_get_purge(t *testing.T) {
	str := "sda               0.10     0.20    0.30    0.40     0.50     0.60     0.70     0.80    0.90    1.00    1.10   1.20   1.30"

	processor := NewDeviceStatsProcessor()

	sample := NewSample()
	item, _ := processor.ParseRow(str)
	sample.Append(item)

	processor.AggregateSample(sample)

	m := processor.GetAndPurgeRrqmps()

	if processor.GetUniqKeys(m)[0] != "sda" {
		t.FailNow()
	}

	if len(processor.rrqmps) != 0 {
		t.FailNow()
	}
}
