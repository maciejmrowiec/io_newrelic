package main

import (
	"strconv"
	"strings"
	"testing"
)

type DummyItem struct {
	name  string
	value float64
}

func (d *DummyItem) GetName() string {
	return d.name
}

func (d *DummyItem) Aggregate(item IItem) {}

type DummyProcessor struct {
	data map[string]ISampleStats
}

func (d *DummyProcessor) GetCmdName() string {
	return ""
}

func (d *DummyProcessor) GetCmdArgs() []string {
	return []string{}
}

func (d *DummyProcessor) AggregateSample(sample ISample) {
	for key, item := range sample.GetMap() {
		ds := item.(*DummyItem)

		if entry, has := d.data[key]; has {
			entry.Append(ds.value)
		} else {
			d.data[key] = NewStatSample(ds.value)
		}
	}
}

func (d *DummyProcessor) NewSample() ISample {
	return NewSample()
}

func (d *DummyProcessor) IsHeader(row string) bool {
	if row == "head" {
		return true
	}

	return false
}

func (d *DummyProcessor) ParseRow(row string) (IItem, error) {
	f := strings.Fields(row)
	val, _ := strconv.ParseFloat(f[1], 64)
	return &DummyItem{
		name:  f[0],
		value: val,
	}, nil
}

func Test_dynamic_collector(t *testing.T) {
	c := NewDynamicCollector(&DummyProcessor{
		data: make(map[string]ISampleStats),
	})

	ch := make(chan string, 100)
	ch <- "head"
	ch <- "t1 1"
	ch <- "t2 2"
	ch <- "t3 3"
	ch <- "head"
	ch <- "t1 3"
	ch <- "t2 2"
	ch <- "t3 1"
	ch <- "head"
	close(ch)

	c.processOutput(ch)

	if !IsFloat64Equal(c.processor.(*DummyProcessor).data["t1"].GetAverage(), 2) {
		t.Log(c.processor.(*DummyProcessor).data["t1"])
		t.FailNow()
	}

	if !IsFloat64Equal(c.processor.(*DummyProcessor).data["t2"].GetAverage(), 2) {
		t.FailNow()
	}

	if !IsFloat64Equal(c.processor.(*DummyProcessor).data["t3"].GetAverage(), 2) {
		t.FailNow()
	}
}
