package main

import (
	"testing"
)

func Test_read_metrica(t *testing.T) {
	top := &IOTopCollector{}
	top.disk_read_rate = make(map[string]*StatSample)
	top.disk_read_rate["test"] = NewStatSample(1.1)

	r := NewReadRatePerCommand(top, "path")
	if r.GetName("test") != "path/test" {
		t.Fail()
	}

	if r.GetUnits() != "kbps" {
		t.Fail()
	}

	if val, _ := r.GetValue("test"); !IsFloat64Equal(val, 1.1) {
		t.Fail()
	}
}

func Test_write_metrica(t *testing.T) {
	top := &IOTopCollector{}
	top.disk_write_rate = make(map[string]*StatSample)
	top.disk_write_rate["test"] = NewStatSample(1.1)

	r := NewWriteRatePerCommand(top, "path")
	if r.GetName("test") != "path/test" {
		t.Fail()
	}

	if r.GetUnits() != "kbps" {
		t.Fail()
	}

	if val, _ := r.GetValue("test"); !IsFloat64Equal(val, 1.1) {
		t.Fail()
	}
}

func Test_io_metrica(t *testing.T) {
	top := &IOTopCollector{}
	top.io_percent = make(map[string]*StatSample)
	top.io_percent["test"] = NewStatSample(1.1)

	r := NewTotalIOPerCommand(top, "path")
	if r.GetName("test") != "path/test" {
		t.Fail()
	}

	if r.GetUnits() != "kbps" {
		t.Fail()
	}

	if val, _ := r.GetValue("test"); !IsFloat64Equal(val, 1.1) {
		t.Fail()
	}
}

func Test_swapin_metrica(t *testing.T) {
	top := &IOTopCollector{}
	top.swapin_percent = make(map[string]*StatSample)
	top.swapin_percent["test"] = NewStatSample(1.1)

	r := NewSwapinPerCommand(top, "path")
	if r.GetName("test") != "path/test" {
		t.Fail()
	}

	if r.GetUnits() != "kbps" {
		t.Fail()
	}

	if val, _ := r.GetValue("test"); !IsFloat64Equal(val, 1.1) {
		t.Fail()
	}
}
