package main

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

type DeviceStats struct {
	name     string
	rrqmps   float64
	wrqmps   float64
	rps      float64
	wps      float64
	rkbps    float64
	wkbps    float64
	avgrq_sz float64
	avgqu_sz float64
	await    float64
	r_await  float64
	w_await  float64
	svctm    float64
	util     float64
}

func (d *DeviceStats) GetName() string {
	return d.name
}

func (d *DeviceStats) Aggregate(item IItem) {}

func DeviceStatsIsHeader(row string) bool {
	return strings.HasPrefix(row, "Device")
}

func DeviceStatsParseRow(row string) (IItem, error) {
	tokens := strings.Fields(row)

	if len(tokens) != 14 {
		return nil, errors.New("Failed to parse iostat")
	}

	var err error
	device_stat := &DeviceStats{}

	device_stat.name = tokens[0]

	if device_stat.rrqmps, err = strconv.ParseFloat(tokens[1], 64); err != nil {
		return nil, err
	}

	if device_stat.wrqmps, err = strconv.ParseFloat(tokens[2], 64); err != nil {
		return nil, err
	}

	if device_stat.rps, err = strconv.ParseFloat(tokens[3], 64); err != nil {
		return nil, err
	}

	if device_stat.wps, err = strconv.ParseFloat(tokens[4], 64); err != nil {
		return nil, err
	}

	if device_stat.rkbps, err = strconv.ParseFloat(tokens[5], 64); err != nil {
		return nil, err
	}

	if device_stat.wkbps, err = strconv.ParseFloat(tokens[6], 64); err != nil {
		return nil, err
	}

	if device_stat.avgrq_sz, err = strconv.ParseFloat(tokens[7], 64); err != nil {
		return nil, err
	}

	if device_stat.avgqu_sz, err = strconv.ParseFloat(tokens[8], 64); err != nil {
		return nil, err
	}

	if device_stat.await, err = strconv.ParseFloat(tokens[9], 64); err != nil {
		return nil, err
	}

	if device_stat.r_await, err = strconv.ParseFloat(tokens[10], 64); err != nil {
		return nil, err
	}

	if device_stat.w_await, err = strconv.ParseFloat(tokens[11], 64); err != nil {
		return nil, err
	}

	if device_stat.svctm, err = strconv.ParseFloat(tokens[12], 64); err != nil {
		return nil, err
	}

	if device_stat.util, err = strconv.ParseFloat(tokens[13], 64); err != nil {
		return nil, err
	}

	return device_stat, nil
}

type DeviceStatsProcessor struct {
	lock sync.Mutex
	// Looks like this could be a seperate type
	rrqmps   map[string]ISampleStats
	wrqmps   map[string]ISampleStats
	rps      map[string]ISampleStats
	wps      map[string]ISampleStats
	rkbps    map[string]ISampleStats
	wkbps    map[string]ISampleStats
	avgrq_sz map[string]ISampleStats
	avgqu_sz map[string]ISampleStats
	await    map[string]ISampleStats
	r_await  map[string]ISampleStats
	w_await  map[string]ISampleStats
	svctm    map[string]ISampleStats
	util     map[string]ISampleStats
}

func NewDeviceStatsProcessor() *DeviceStatsProcessor {
	return &DeviceStatsProcessor{
		rrqmps:   make(map[string]ISampleStats),
		wrqmps:   make(map[string]ISampleStats),
		rps:      make(map[string]ISampleStats),
		wps:      make(map[string]ISampleStats),
		rkbps:    make(map[string]ISampleStats),
		wkbps:    make(map[string]ISampleStats),
		avgrq_sz: make(map[string]ISampleStats),
		avgqu_sz: make(map[string]ISampleStats),
		await:    make(map[string]ISampleStats),
		r_await:  make(map[string]ISampleStats),
		w_await:  make(map[string]ISampleStats),
		svctm:    make(map[string]ISampleStats),
		util:     make(map[string]ISampleStats),
	}
}

func (p *DeviceStatsProcessor) GetCmdName() string {
	return "iostat"
}

func (p *DeviceStatsProcessor) GetCmdArgs() []string {
	return []string{"-xkd", "10"}
}

func (d *DeviceStatsProcessor) AggregateSample(sample ISample) {
	d.lock.Lock()

	for key, item := range sample.GetMap() {
		ds := item.(*DeviceStats)

		if entry, has := d.rrqmps[key]; has {
			entry.Append(ds.rrqmps)
		} else {
			d.rrqmps[key] = NewStatSample(ds.rrqmps)
		}

		if entry, has := d.wrqmps[key]; has {
			entry.Append(ds.wrqmps)
		} else {
			d.wrqmps[key] = NewStatSample(ds.wrqmps)
		}

		if entry, has := d.rps[key]; has {
			entry.Append(ds.rps)
		} else {
			d.rps[key] = NewStatSample(ds.rps)
		}

		if entry, has := d.wps[key]; has {
			entry.Append(ds.wps)
		} else {
			d.wps[key] = NewStatSample(ds.wps)
		}

		if entry, has := d.rkbps[key]; has {
			entry.Append(ds.rkbps)
		} else {
			d.rkbps[key] = NewStatSample(ds.rkbps)
		}

		if entry, has := d.wkbps[key]; has {
			entry.Append(ds.wkbps)
		} else {
			d.wkbps[key] = NewStatSample(ds.wkbps)
		}

		if entry, has := d.avgrq_sz[key]; has {
			entry.Append(ds.avgrq_sz)
		} else {
			d.avgrq_sz[key] = NewStatSample(ds.avgrq_sz)
		}

		if entry, has := d.avgqu_sz[key]; has {
			entry.Append(ds.avgqu_sz)
		} else {
			d.avgqu_sz[key] = NewStatSample(ds.avgqu_sz)
		}

		if entry, has := d.await[key]; has {
			entry.Append(ds.await)
		} else {
			d.await[key] = NewStatSample(ds.await)
		}

		if entry, has := d.r_await[key]; has {
			entry.Append(ds.r_await)
		} else {
			d.r_await[key] = NewStatSample(ds.r_await)
		}

		if entry, has := d.w_await[key]; has {
			entry.Append(ds.w_await)
		} else {
			d.w_await[key] = NewStatSample(ds.w_await)
		}

		if entry, has := d.svctm[key]; has {
			entry.Append(ds.svctm)
		} else {
			d.svctm[key] = NewStatSample(ds.svctm)
		}

		if entry, has := d.util[key]; has {
			entry.Append(ds.util)
		} else {
			d.util[key] = NewStatSample(ds.util)
		}

	}

	d.lock.Unlock()
}

func (d *DeviceStatsProcessor) NewSample() ISample {
	return NewSample()
}

func (d *DeviceStatsProcessor) IsHeader(row string) bool {
	return DeviceStatsIsHeader(row)
}

func (d *DeviceStatsProcessor) ParseRow(row string) (IItem, error) {
	return DeviceStatsParseRow(row)
}

func (d *DeviceStatsProcessor) GetAndPurgeRrqmps() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.rrqmps
	d.rrqmps = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeWrqmps() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.wrqmps
	d.wrqmps = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeRps() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.rps
	d.rps = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeWps() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.wps
	d.wps = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeRkbps() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.rkbps
	d.rkbps = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeWkbps() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.wkbps
	d.wkbps = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeAvgrq_sz() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.avgrq_sz
	d.avgrq_sz = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeAvgqu_sz() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.avgqu_sz
	d.avgqu_sz = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeAwait() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.await
	d.await = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeRawait() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.r_await
	d.r_await = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeWawait() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.w_await
	d.w_await = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeSvctm() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.svctm
	d.svctm = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

func (d *DeviceStatsProcessor) GetAndPurgeUtil() map[string]ISampleStats {
	d.lock.Lock()
	samples := d.util
	d.util = make(map[string]ISampleStats)
	d.lock.Unlock()

	return samples
}

// this could be a method on map[string]*StatSample
func (d *DeviceStatsProcessor) GetUniqKeys(data map[string]ISampleStats) []string {
	var sample_list []string

	for key := range data {
		sample_list = append(sample_list, key)
	}

	return sample_list
}
