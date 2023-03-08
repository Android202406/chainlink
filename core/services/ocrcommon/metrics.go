package ocrcommon

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
)

var _ commontypes.Metric = &DefaultMetric{nil}

type DefaultMetric struct {
	prometheus.Gauge
}

var _ commontypes.Metrics = &MetricVecFactory{nil}

type MetricVecFactory struct {
	generatorFn func(name string, help string, labelNames ...string) (commontypes.MetricVec, error)
}

func (f *MetricVecFactory) NewMetricVec(name string, help string, labelNames ...string) (commontypes.MetricVec, error) {
	return f.generatorFn(name, help, labelNames...)
}

func NewMetricVecFactory(generator func(name string, help string, labelNames ...string) (commontypes.MetricVec, error)) *MetricVecFactory {
	return &MetricVecFactory{
		generatorFn: generator,
	}
}

var _ commontypes.MetricVec = &DefaultMetricVec{nil}

type DefaultMetricVec struct {
	*prometheus.GaugeVec
}

func NewDefaultMetricVec(name string, help string, labelNames ...string) (commontypes.MetricVec, error) {
	gv := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: name,
		Subsystem: "",
		Name:      name,
		Help:      help,
	}, labelNames)

	return &DefaultMetricVec{
		GaugeVec: gv,
	}, nil
}

func (mv *DefaultMetricVec) GetMetricWith(labels map[string]string) (commontypes.Metric, error) {
	return mv.GaugeVec.GetMetricWith(labels)
}
