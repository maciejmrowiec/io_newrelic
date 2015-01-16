package main

type TotalIOPerCommand struct {
	collector *IOTopCollector
	path      string
	samples   map[string]float64
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

	return t.samples[id], nil
}

func (t *TotalIOPerCommand) GetIdList() []string {
	t.samples = t.collector.GetAndPurgeIOPercent()
	return t.collector.GetUniqKeys(t.samples)
}

type ReadRatePerCommand struct {
	collector *IOTopCollector
	path      string
	samples   map[string]float64
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

	return t.samples[id], nil
}

func (t *ReadRatePerCommand) GetIdList() []string {
	t.samples = t.collector.GetAndPurgeDiskReadRate()
	return t.collector.GetUniqKeys(t.samples)
}

type WriteRatePerCommand struct {
	collector *IOTopCollector
	path      string
	samples   map[string]float64
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

	return t.samples[id], nil
}

func (t *WriteRatePerCommand) GetIdList() []string {
	t.samples = t.collector.GetAndPurgeDiskWriteRate()
	return t.collector.GetUniqKeys(t.samples)
}
