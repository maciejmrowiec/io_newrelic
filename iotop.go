package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
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
	Name            string
}

// expected format:
//13:05:19 24086 be/4 cfapache    0.00 K/s    0.00 K/s  0.00 %  0.00 % httpd -k start
func NewProcessIOFromString(row_data string) (*ProcessIO, error) {
	tokens := RowRegex.Split(row_data, 13)
	if len(tokens) != 13 {
		return nil, errors.New("Failed to parse row data for sample")
	}

	var err error
	process_io := &ProcessIO{}

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
	process_io.Name = tokens[12]

	return process_io, nil
}

func (p *ProcessIO) Aggregate(process *ProcessIO) {
	p.disk_read_rate += process.disk_read_rate
	p.disk_write_rate += process.disk_write_rate
	p.swapin_percent += process.swapin_percent
	p.io_percent += process.io_percent
}

// Total DISK READ: 0.20 K/s | Total DISK WRITE: 0.01 K/s
func IsSampleSummary(row_data string) bool {
	return strings.Contains(row_data, "Total DISK READ")
}

type Sample struct {
	processes_io map[string]*ProcessIO
}

func NewSample() *Sample {
	return &Sample{
		processes_io: make(map[string]*ProcessIO),
	}
}

func (s *Sample) Append(process *ProcessIO) {
	if val, has := s.processes_io[process.Name]; has {
		val.Aggregate(process)
	} else {
		s.processes_io[process.Name] = process
	}
}

func (s *Sample) Empty() bool {
	if len(s.processes_io) > 0 {
		return false
	}

	return true
}

type IOTopCollector struct {
	disk_read_rate    map[string]*StatSample
	disk_write_rate   map[string]*StatSample
	swapin_percent    map[string]*StatSample
	io_percent        map[string]*StatSample
	measurements_lock sync.RWMutex
}

func (i *IOTopCollector) Run() {
	i.measurements_lock.Lock()

	command_pipe := make(chan string, 1000)

	go i.executeIOTop(command_pipe)
	go i.processOutput(command_pipe)

	i.measurements_lock.Unlock()

}

func (i *IOTopCollector) GetAndPurgeIOPercent() map[string]*StatSample {
	i.measurements_lock.Lock()
	samples := i.io_percent
	i.io_percent = nil
	i.measurements_lock.Unlock()

	return samples
}

func (i *IOTopCollector) GetAndPurgeSwapinPercent() map[string]*StatSample {
	i.measurements_lock.Lock()
	samples := i.swapin_percent
	i.swapin_percent = nil
	i.measurements_lock.Unlock()

	return samples
}

func (i *IOTopCollector) GetAndPurgeDiskReadRate() map[string]*StatSample {
	i.measurements_lock.Lock()
	samples := i.disk_read_rate
	i.disk_read_rate = nil
	i.measurements_lock.Unlock()

	return samples
}

func (i *IOTopCollector) GetAndPurgeDiskWriteRate() map[string]*StatSample {
	i.measurements_lock.Lock()
	samples := i.disk_write_rate
	i.disk_write_rate = nil
	i.measurements_lock.Unlock()

	return samples
}

func (i *IOTopCollector) GetUniqKeys(data map[string]*StatSample) []string {
	var sample_list []string

	for key := range data {
		sample_list = append(sample_list, key)
	}

	return sample_list
}

func (i *IOTopCollector) processOutput(data <-chan string) {
	sample := NewSample()

	for row := range data {

		if IsSampleSummary(row) {
			if sample != nil && !sample.Empty() {
				i.FlushSample(sample)
			}

			sample = NewSample()
		} else {
			p, err := NewProcessIOFromString(row)
			if err != nil {
				continue
			}

			sample.Append(p)
		}
	}
}

func (i *IOTopCollector) FlushSample(sample *Sample) {
	i.measurements_lock.Lock()
	if i.disk_read_rate == nil {
		i.disk_read_rate = make(map[string]*StatSample, len(sample.processes_io))
	}

	if i.disk_write_rate == nil {
		i.disk_write_rate = make(map[string]*StatSample, len(sample.processes_io))
	}

	if i.io_percent == nil {
		i.io_percent = make(map[string]*StatSample, len(sample.processes_io))
	}

	if i.swapin_percent == nil {
		i.swapin_percent = make(map[string]*StatSample, len(sample.processes_io))
	}

	for key, val := range sample.processes_io {
		if entry, has := i.disk_read_rate[key]; has {
			entry.Append(val.disk_read_rate)
		} else {
			i.disk_read_rate[key] = NewStatSample(val.disk_read_rate)
		}

		if entry, has := i.disk_write_rate[key]; has {
			entry.Append(val.disk_write_rate)
		} else {
			i.disk_write_rate[key] = NewStatSample(val.disk_write_rate)
		}

		if entry, has := i.io_percent[key]; has {
			entry.Append(val.io_percent)
		} else {
			i.io_percent[key] = NewStatSample(val.io_percent)
		}

		if entry, has := i.swapin_percent[key]; has {
			entry.Append(val.swapin_percent)
		} else {
			i.swapin_percent[key] = NewStatSample(val.swapin_percent)
		}
	}
	i.measurements_lock.Unlock()
}

func (i *IOTopCollector) executeIOTop(ch chan<- string) {
	cmd := exec.Command("iotop", "-bPkqqt")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	in := bufio.NewScanner(stdout)

	for in.Scan() {
		ch <- in.Text()
	}

	if err := in.Err(); err != nil {
		log.Fatal(err)
	}

	close(ch)
}

type StatSample struct {
	total float64
	count float64
}

func NewStatSample(value float64) *StatSample {
	return &StatSample{
		total: value,
		count: 1,
	}
}

func (s *StatSample) Append(value float64) {
	s.count += 1
	s.total += value
}

func (s *StatSample) GetAverage() float64 {
	return s.total / s.count
}
