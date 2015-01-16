package main

import (
	"bufio"
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

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
	tokens := regexp.MustCompile("[\\s]+").Split(row_data, 13)
	if len(tokens) != 13 {
		return nil, errors.New("Failed to parse row data for sample")
	}

	var err error
	process_io := &ProcessIO{}

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
	return strings.HasPrefix(row_data, "Total DISK READ:")
}

type Sample struct {
	processes_io map[string]*ProcessIO
}

func NewSample() *Sample {
	return &Sample{
		processes_io: make(map[string]*ProcessIO),
	}
}

func (s *Sample) Append(name string, process *ProcessIO) {
	if val, has := s.processes_io[name]; has {
		val.Aggregate(process)
	} else {
		s.processes_io[name] = process
	}
}

func (s *Sample) Empty() bool {
	if len(s.processes_io) > 0 {
		return false
	}

	return true
}

type IOTopCollector struct {
	command_pipe      chan string
	disk_read_rate    []map[string]float64
	disk_write_rate   []map[string]float64
	swapin_percent    []map[string]float64
	io_percent        []map[string]float64
	measurements_lock sync.RWMutex
}

func (i *IOTopCollector) Run() {
	i.measurements_lock.Lock()

	i.command_pipe = make(chan string, 1000)

	go i.executeIOTop(i.command_pipe)
	go i.processOutput(i.command_pipe)

	i.measurements_lock.Unlock()

}

func (i *IOTopCollector) AverageSamples(sample_list []map[string]float64) map[string]float64 {

	aggregated := make(map[string]struct {
		total float64
		count float64
	})

	for _, sample := range sample_list {
		for key, value := range sample {
			entry := aggregated[key]
			entry.total += value
			entry.count += 1
			aggregated[key] = entry
		}
	}

	summary := make(map[string]float64, len(aggregated))
	for key, val := range aggregated {
		if val.count > 0 {
			summary[key] = val.total / val.count
		} else {
			summary[key] = 0
		}
	}

	return summary
}

func (i *IOTopCollector) GetAndPurgeIOPercent() map[string]float64 {
	i.measurements_lock.Lock()
	samples := i.io_percent
	i.io_percent = nil
	i.measurements_lock.Unlock()

	return i.AverageSamples(samples)
}

func (i *IOTopCollector) GetAndPurgeSwapinPercent() map[string]float64 {
	i.measurements_lock.Lock()
	samples := i.swapin_percent
	i.swapin_percent = nil
	i.measurements_lock.Unlock()

	return i.AverageSamples(samples)
}

func (i *IOTopCollector) GetAndPurgeDiskReadRate() map[string]float64 {
	i.measurements_lock.Lock()
	samples := i.disk_read_rate
	i.disk_read_rate = nil
	i.measurements_lock.Unlock()

	return i.AverageSamples(samples)
}

func (i *IOTopCollector) GetAndPurgeDiskWriteRate() map[string]float64 {
	i.measurements_lock.Lock()
	samples := i.disk_write_rate
	i.disk_write_rate = nil
	i.measurements_lock.Unlock()

	return i.AverageSamples(samples)
}

func (i *IOTopCollector) GetUniqKeys(data map[string]float64) []string {
	var sample_list []string

	for key := range data {
		sample_list = append(sample_list, key)
	}

	return sample_list
}

func (i *IOTopCollector) processOutput(data <-chan string) {
	var sample *Sample

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

			sample.Append(p.Name, p)
		}
	}
}

func (i *IOTopCollector) FlushSample(sample *Sample) {
	i.measurements_lock.Lock()

	for key, val := range sample.processes_io {
		iop := make(map[string]float64)
		iop[key] = val.io_percent
		i.io_percent = append(i.io_percent, iop)

		swp := make(map[string]float64)
		swp[key] = val.swapin_percent
		i.swapin_percent = append(i.swapin_percent, swp)

		drr := make(map[string]float64)
		drr[key] = val.disk_read_rate
		i.disk_read_rate = append(i.disk_read_rate, drr)

		dwr := make(map[string]float64)
		dwr[key] = val.disk_write_rate
		i.disk_write_rate = append(i.disk_write_rate, dwr)

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
