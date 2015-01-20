package main

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

type Sample struct {
	data map[string]IItem
}

func NewSample() *Sample {
	return &Sample{
		data: make(map[string]IItem, 10),
	}
}

func (d *Sample) Append(item IItem) {
	name := item.GetName()

	if val, has := d.data[name]; has {
		val.Aggregate(item)
	} else {
		d.data[name] = item
	}
}

func (d *Sample) Empty() bool {
	if len(d.data) > 0 {
		return false
	}
	return true
}

func (d *Sample) GetMap() map[string]IItem {
	return d.data
}
