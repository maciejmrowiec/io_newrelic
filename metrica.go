package main

type TotalIOPerCommand struct {
	processor *ProcessIOProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewTotalIOPerCommand(processor *ProcessIOProcessor, path string) *TotalIOPerCommand {
	return &TotalIOPerCommand{
		processor: processor,
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
	t.samples = t.processor.GetAndPurgeIOPercent()
	return t.processor.GetUniqKeys(t.samples)
}

type ReadRatePerCommand struct {
	processor *ProcessIOProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewReadRatePerCommand(processor *ProcessIOProcessor, path string) *ReadRatePerCommand {
	return &ReadRatePerCommand{
		processor: processor,
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
	t.samples = t.processor.GetAndPurgeDiskReadRate()
	return t.processor.GetUniqKeys(t.samples)
}

type WriteRatePerCommand struct {
	processor *ProcessIOProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewWriteRatePerCommand(processor *ProcessIOProcessor, path string) *WriteRatePerCommand {
	return &WriteRatePerCommand{
		processor: processor,
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
	t.samples = t.processor.GetAndPurgeDiskWriteRate()
	return t.processor.GetUniqKeys(t.samples)
}

type SwapinPerCommand struct {
	processor *ProcessIOProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewSwapinPerCommand(processor *ProcessIOProcessor, path string) *SwapinPerCommand {
	return &SwapinPerCommand{
		processor: processor,
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
	t.samples = t.processor.GetAndPurgeSwapinPercent()
	return t.processor.GetUniqKeys(t.samples)
}

type RrqmpsPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewRrqmpsPerDevice(processor *DeviceStatsProcessor, path string) *RrqmpsPerDevice {
	return &RrqmpsPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *RrqmpsPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *RrqmpsPerDevice) GetUnits() string {
	return "requestps"
}

func (r *RrqmpsPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *RrqmpsPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeRrqmps()
	return r.processor.GetUniqKeys(r.samples)
}

type WrqmpsPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewWrqmpsPerDevice(processor *DeviceStatsProcessor, path string) *WrqmpsPerDevice {
	return &WrqmpsPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *WrqmpsPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *WrqmpsPerDevice) GetUnits() string {
	return "requestps"
}

func (r *WrqmpsPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *WrqmpsPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeWrqmps()
	return r.processor.GetUniqKeys(r.samples)
}

type RpsPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewRpsPerDevice(processor *DeviceStatsProcessor, path string) *RpsPerDevice {
	return &RpsPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *RpsPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *RpsPerDevice) GetUnits() string {
	return "rps"
}

func (r *RpsPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *RpsPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeRps()
	return r.processor.GetUniqKeys(r.samples)
}

type WpsPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewWpsPerDevice(processor *DeviceStatsProcessor, path string) *WpsPerDevice {
	return &WpsPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *WpsPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *WpsPerDevice) GetUnits() string {
	return "wps"
}

func (r *WpsPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *WpsPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeWps()
	return r.processor.GetUniqKeys(r.samples)
}

type RkbpsPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewRkbpsPerDevice(processor *DeviceStatsProcessor, path string) *RkbpsPerDevice {
	return &RkbpsPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *RkbpsPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *RkbpsPerDevice) GetUnits() string {
	return "kbps"
}

func (r *RkbpsPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *RkbpsPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeRkbps()
	return r.processor.GetUniqKeys(r.samples)
}

type WkbpsPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewWkbpsPerDevice(processor *DeviceStatsProcessor, path string) *WkbpsPerDevice {
	return &WkbpsPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *WkbpsPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *WkbpsPerDevice) GetUnits() string {
	return "kbps"
}

func (r *WkbpsPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *WkbpsPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeWkbps()
	return r.processor.GetUniqKeys(r.samples)
}

type AvgrqszPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewAvgrqszPerDevice(processor *DeviceStatsProcessor, path string) *AvgrqszPerDevice {
	return &AvgrqszPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *AvgrqszPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *AvgrqszPerDevice) GetUnits() string {
	return "sectors"
}

func (r *AvgrqszPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *AvgrqszPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeAvgrq_sz()
	return r.processor.GetUniqKeys(r.samples)
}

type AvgquszPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewAvgquszPerDevice(processor *DeviceStatsProcessor, path string) *AvgquszPerDevice {
	return &AvgquszPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *AvgquszPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *AvgquszPerDevice) GetUnits() string {
	return "count"
}

func (r *AvgquszPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *AvgquszPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeAvgqu_sz()
	return r.processor.GetUniqKeys(r.samples)
}

type AwaitPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewAwaitPerDevice(processor *DeviceStatsProcessor, path string) *AwaitPerDevice {
	return &AwaitPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *AwaitPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *AwaitPerDevice) GetUnits() string {
	return "ms"
}

func (r *AwaitPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *AwaitPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeAwait()
	return r.processor.GetUniqKeys(r.samples)
}

type RawaitPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewRawaitPerDevice(processor *DeviceStatsProcessor, path string) *RawaitPerDevice {
	return &RawaitPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *RawaitPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *RawaitPerDevice) GetUnits() string {
	return "ms"
}

func (r *RawaitPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *RawaitPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeRawait()
	return r.processor.GetUniqKeys(r.samples)
}

type WawaitPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewWawaitPerDevice(processor *DeviceStatsProcessor, path string) *WawaitPerDevice {
	return &WawaitPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *WawaitPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *WawaitPerDevice) GetUnits() string {
	return "ms"
}

func (r *WawaitPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *WawaitPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeWawait()
	return r.processor.GetUniqKeys(r.samples)
}

type SvctmPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewSvctmPerDevice(processor *DeviceStatsProcessor, path string) *SvctmPerDevice {
	return &SvctmPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *SvctmPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *SvctmPerDevice) GetUnits() string {
	return "ms"
}

func (r *SvctmPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *SvctmPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeSvctm()
	return r.processor.GetUniqKeys(r.samples)
}

type UtilPerDevice struct {
	processor *DeviceStatsProcessor
	path      string
	samples   map[string]ISampleStats
}

func NewUtilPerDevice(processor *DeviceStatsProcessor, path string) *UtilPerDevice {
	return &UtilPerDevice{
		processor: processor,
		path:      path,
	}
}

func (r *UtilPerDevice) GetName(id string) string {
	return r.path + "/" + id
}

func (r *UtilPerDevice) GetUnits() string {
	return "%"
}

func (r *UtilPerDevice) GetValue(id string) (float64, error) {
	return r.samples[id].GetAverage(), nil
}

func (r *UtilPerDevice) GetIdList() []string {
	r.samples = r.processor.GetAndPurgeUtil()
	return r.processor.GetUniqKeys(r.samples)
}
