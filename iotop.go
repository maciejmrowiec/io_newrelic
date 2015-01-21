package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func init() {
	RowRegex = regexp.MustCompile("[\\s]+")
}

var RowRegex *regexp.Regexp

type ProcessIO struct {
	disk_read_rate  float64
	disk_write_rate float64
	swapin_percent  float64
	io_percent      float64
	name            string
}

func (p *ProcessIO) GetName() string {
	return p.name
}

func (p *ProcessIO) Aggregate(item IItem) {
	process := item.(*ProcessIO)
	p.disk_read_rate += process.disk_read_rate
	p.disk_write_rate += process.disk_write_rate
	p.swapin_percent += process.swapin_percent
	p.io_percent += process.io_percent
}

func ProcessIOParseRow(row_data string) (IItem, error) {
	tokens := RowRegex.Split(row_data, 13)
	if len(tokens) != 13 {
		return nil, errors.New("Failed to parse row data for sample")
	}

	var err error
	process_io := new(ProcessIO)

	if process_io.disk_read_rate, err = strconv.ParseFloat(tokens[4], 64); err != nil {
		return nil, err
	}
	if process_io.disk_write_rate, err = strconv.ParseFloat(tokens[6], 64); err != nil {
		return nil, err
	}
	if process_io.swapin_percent, err = strconv.ParseFloat(tokens[8], 64); err != nil {
		return nil, err
	}
	if process_io.io_percent, err = strconv.ParseFloat(tokens[10], 64); err != nil {
		return nil, err
	}
	process_io.name = tokens[12]

	return process_io, nil
}

func ProcessIOIsHeader(row_data string) bool {
	return strings.Contains(row_data, "Total DISK READ")
}

type ProcessIOProcessor struct {
	disk_read_rate    map[string]ISampleStats
	disk_write_rate   map[string]ISampleStats
	swapin_percent    map[string]ISampleStats
	io_percent        map[string]ISampleStats
	measurements_lock sync.Mutex
}

func NewProcessIOProcessor() *ProcessIOProcessor {
	return &ProcessIOProcessor{
		disk_read_rate:  make(map[string]ISampleStats, 200),
		disk_write_rate: make(map[string]ISampleStats, 200),
		swapin_percent:  make(map[string]ISampleStats, 200),
		io_percent:      make(map[string]ISampleStats, 200),
	}
}

func (i *ProcessIOProcessor) GetAndPurgeIOPercent() map[string]ISampleStats {
	i.measurements_lock.Lock()
	samples := i.io_percent
	i.io_percent = make(map[string]ISampleStats, 200)
	i.measurements_lock.Unlock()

	return samples
}

func (i *ProcessIOProcessor) GetAndPurgeSwapinPercent() map[string]ISampleStats {
	i.measurements_lock.Lock()
	samples := i.swapin_percent
	i.swapin_percent = make(map[string]ISampleStats, 200)
	i.measurements_lock.Unlock()

	return samples
}

func (i *ProcessIOProcessor) GetAndPurgeDiskReadRate() map[string]ISampleStats {
	i.measurements_lock.Lock()
	samples := i.disk_read_rate
	i.disk_read_rate = make(map[string]ISampleStats, 200)
	i.measurements_lock.Unlock()

	return samples
}

func (i *ProcessIOProcessor) GetAndPurgeDiskWriteRate() map[string]ISampleStats {
	i.measurements_lock.Lock()
	samples := i.disk_write_rate
	i.disk_write_rate = make(map[string]ISampleStats, 200)
	i.measurements_lock.Unlock()

	return samples
}

func (i *ProcessIOProcessor) GetUniqKeys(data map[string]ISampleStats) []string {
	var sample_list []string

	for key := range data {
		sample_list = append(sample_list, key)
	}

	return sample_list
}

func (p *ProcessIOProcessor) GetCmdName() string {
	return "iotop"
}

func (p *ProcessIOProcessor) GetCmdArgs() []string {
	return []string{"-bPkqqt"}
}

func (i *ProcessIOProcessor) AggregateSample(sample ISample) {
	i.measurements_lock.Lock()

	for key, val := range sample.GetMap() {
		pio := val.(*ProcessIO)

		if entry, has := i.disk_read_rate[key]; has {
			entry.Append(pio.disk_read_rate)
		} else {
			i.disk_read_rate[key] = NewStatSample(pio.disk_read_rate)
		}

		if entry, has := i.disk_write_rate[key]; has {
			entry.Append(pio.disk_write_rate)
		} else {
			i.disk_write_rate[key] = NewStatSample(pio.disk_write_rate)
		}

		if entry, has := i.io_percent[key]; has {
			entry.Append(pio.io_percent)
		} else {
			i.io_percent[key] = NewStatSample(pio.io_percent)
		}

		if entry, has := i.swapin_percent[key]; has {
			entry.Append(pio.swapin_percent)
		} else {
			i.swapin_percent[key] = NewStatSample(pio.swapin_percent)
		}
	}
	i.measurements_lock.Unlock()
}

func (p *ProcessIOProcessor) NewSample() ISample {
	return NewSample()
}

func (p *ProcessIOProcessor) IsHeader(row string) bool {
	return ProcessIOIsHeader(row)
}

func (p *ProcessIOProcessor) ParseRow(row string) (IItem, error) {
	return ProcessIOParseRow(row)
}
