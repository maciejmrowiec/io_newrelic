package main

import (
	"bytes"
	"github.com/yvasiyarov/newrelic_platform_go"
	"log"
	"math"
	"sync"
)

type IDynamicMetrica interface {
	GetUnits() string
	GetName(metrica_id string) string
	GetValue(metrica_id string) (float64, error)
	GetIdList() []string
}

type MetricaContainer struct {
	MetricaModel  newrelic_platform_go.IMetrica
	MetricaIDList []string
	RWLock        sync.RWMutex
}

type DynamicPluginComponent struct {
	Name          string                                       `json:"name"`
	GUID          string                                       `json:"guid"`
	Duration      int                                          `json:"duration"`
	Metrics       map[string]newrelic_platform_go.MetricaValue `json:"metrics"`
	MetricaModels []IDynamicMetrica                            `json:"-"`
}

func NewDynamicPluginComponent(name string, guid string) *DynamicPluginComponent {
	c := &DynamicPluginComponent{
		Name: name,
		GUID: guid,
	}
	return c
}

func (component *DynamicPluginComponent) AddMetrica(model newrelic_platform_go.IMetrica) {
	log.Fatal("Unsupported Method")
}

func (component *DynamicPluginComponent) AddDynamicMetrica(model IDynamicMetrica) {
	component.MetricaModels = append(component.MetricaModels, model)
}

func (component *DynamicPluginComponent) ClearSentData() {
	component.Metrics = nil
}

func (component *DynamicPluginComponent) SetDuration(duration int) {
	component.Duration = duration
}

func (component *DynamicPluginComponent) GetMetricaKey(metrica IDynamicMetrica, id string) string {
	var keyBuffer bytes.Buffer

	keyBuffer.WriteString("Component/")
	keyBuffer.WriteString(metrica.GetName(id))
	keyBuffer.WriteString("[")
	keyBuffer.WriteString(metrica.GetUnits())
	keyBuffer.WriteString("]")

	return keyBuffer.String()
}

func (component *DynamicPluginComponent) Harvest(plugin newrelic_platform_go.INewrelicPlugin) newrelic_platform_go.ComponentData {

	component.Metrics = make(map[string]newrelic_platform_go.MetricaValue, len(component.MetricaModels))

	for i := 0; i < len(component.MetricaModels); i++ {
		model := component.MetricaModels[i]

		ids := model.GetIdList()

		for _, id := range ids {

			metricaKey := component.GetMetricaKey(model, id)

			if newValue, err := model.GetValue(id); err == nil {
				if math.IsInf(newValue, 0) || math.IsNaN(newValue) {
					newValue = 0
				}

				if existMetric, ok := component.Metrics[metricaKey]; ok {
					if floatExistVal, ok := existMetric.(float64); ok {
						component.Metrics[metricaKey] = newrelic_platform_go.NewAggregatedMetricaValue(floatExistVal, newValue)
					} else if aggregatedValue, ok := existMetric.(*newrelic_platform_go.AggregatedMetricaValue); ok {
						aggregatedValue.Aggregate(newValue)
					} else {
						panic("Invalid type in metrica value")
					}
				} else {
					component.Metrics[metricaKey] = newValue
				}
			} else {
				log.Printf("Can not get metrica: %v, got error:%v", metricaKey, err)
			}
		}
	}
	return component
}
