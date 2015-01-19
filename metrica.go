package main

type TotalIOPerCommand struct {
	collector *IOTopCollector
	path      string
	samples   map[string]*StatSample
}

func NewTotalIOPerCommand(collector *IOTopCollector, path string) *TotalIOPerCommand {
	return &TotalIOPerCommand{
		collector: collector,
		path:      path,
	}
}

func (t *TotalIOPerCommand) GetUnits() string {
	return "%"
}

func (t *TotalIOPerCommand) GetName(id string) string {
	return t.path + "/" + id
}

func (t *TotalIOPerCommand) GetValue(id string) (float64, error) {
	return t.samples[id].GetAverage(), nil
}

func (t *TotalIOPerCommand) GetIdList() []string {
	t.samples = t.collector.GetAndPurgeIOPercent()
	return t.collector.GetUniqKeys(t.samples)
}

type ReadRatePerCommand struct {
	collector *IOTopCollector
	path      string
	samples   map[string]*StatSample
}

func NewReadRatePerCommand(collector *IOTopCollector, path string) *ReadRatePerCommand {
	return &ReadRatePerCommand{
		collector: collector,
		path:      path,
	}
}

func (t *ReadRatePerCommand) GetUnits() string {
	return "kbps"
}

func (t *ReadRatePerCommand) GetName(id string) string {
	return t.path + "/" + id
}

func (t *ReadRatePerCommand) GetValue(id string) (float64, error) {
	return t.samples[id].GetAverage(), nil
}

func (t *ReadRatePerCommand) GetIdList() []string {
	t.samples = t.collector.GetAndPurgeDiskReadRate()
	return t.collector.GetUniqKeys(t.samples)
}

type WriteRatePerCommand struct {
	collector *IOTopCollector
	path      string
	samples   map[string]*StatSample
}

func NewWriteRatePerCommand(collector *IOTopCollector, path string) *WriteRatePerCommand {
	return &WriteRatePerCommand{
		collector: collector,
		path:      path,
	}
}

func (t *WriteRatePerCommand) GetUnits() string {
	return "kbps"
}

func (t *WriteRatePerCommand) GetName(id string) string {
	return t.path + "/" + id
}

func (t *WriteRatePerCommand) GetValue(id string) (float64, error) {
	return t.samples[id].GetAverage(), nil
}

func (t *WriteRatePerCommand) GetIdList() []string {
	t.samples = t.collector.GetAndPurgeDiskWriteRate()
	return t.collector.GetUniqKeys(t.samples)
}

type SwapinPerCommand struct {
	collector *IOTopCollector
	path      string
	samples   map[string]*StatSample
}

func NewSwapinPerCommand(collector *IOTopCollector, path string) *SwapinPerCommand {
	return &SwapinPerCommand{
		collector: collector,
		path:      path,
	}
}

func (t *SwapinPerCommand) GetUnits() string {
	return "%"
}

func (t *SwapinPerCommand) GetName(id string) string {
	return t.path + "/" + id
}

func (t *SwapinPerCommand) GetValue(id string) (float64, error) {
	return t.samples[id].GetAverage(), nil
}

func (t *SwapinPerCommand) GetIdList() []string {
	t.samples = t.collector.GetAndPurgeSwapinPercent()
	return t.collector.GetUniqKeys(t.samples)
}
